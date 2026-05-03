package matchingservice

import (
	"context"
	"fmt"
	"gocasts/gameapp/dto"
	"gocasts/gameapp/entity"
	"gocasts/gameapp/pkg/richerror"
	"gocasts/gameapp/pkg/timestamp"
	"sync"
	"time"

	funk "github.com/thoas/go-funk"
)

type Config struct {
	WaitingTimeout time.Duration `koanf:"timeout"`
}

type Repository interface {
	AddToWaitingList(userID uint, category entity.Category) error
	GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
}

type PresenceClient interface {
	GetPresence(ctx context.Context, req dto.GetPresenceRequest) (dto.GetPresenceResponse, error)
}

type Service struct {
	config         Config
	repo           Repository
	presenceClient PresenceClient
}

func New(config Config, repo Repository, presenceClient PresenceClient) Service {
	return Service{config: config, repo: repo, presenceClient: presenceClient}
}

func (s Service) AddToWaitingList(req dto.AddToWaitingListRequest) (dto.AddToWaitingListResponse, error) {
	op := richerror.Op("matchingservice.AddToWaitingList")

	// add the user to the waiting list for the given category if they already dont exist
	err := s.repo.AddToWaitingList(req.UserID, req.Category)
	if err != nil {
		return dto.AddToWaitingListResponse{}, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return dto.AddToWaitingListResponse{Timeout: s.config.WaitingTimeout}, nil
}

func (s Service) MatchWaitedUsers(ctx context.Context, _ dto.MatchWaitedUsersRequest) (dto.MatchWaitedUsersResponse, error) {
	const op = richerror.Op("matchingservice.MatchingWaitedUsers")

	//get a list of existing categories
	// get the waiting lists for all categories
	// see who is online in each list
	// match users 2 by 2 randomly
	// create event and publish it to be broker and remove from the waiting list
	var wg sync.WaitGroup
	for _, category := range entity.CategoryList() {
		wg.Add(1)
		go s.match(ctx, category, &wg)
	}

	wg.Wait()

	return dto.MatchWaitedUsersResponse{}, nil
}

func (s Service) match(ctx context.Context, category entity.Category, wg *sync.WaitGroup) {
	const op = richerror.Op("matchingservice.match")

	defer wg.Done()

	waitingList, err := s.repo.GetWaitingListByCategory(ctx, category)
	if err != nil {
		// TODO - log error
		// TODO - update metrics
		return
	}

	userIDs := make([]uint, len(waitingList))
	for _, i := range waitingList {
		userIDs = append(userIDs, i.UserID)
	}

	presenceList, err := s.presenceClient.GetPresence(ctx, dto.GetPresenceRequest{
		UserIDs: userIDs,
	})
	if err != nil {
		// TODO - log error
		// TODO - update metrics
		return
	}

	//TODO - merge presenceList with waitingList based on userID
	// also consider the presence timestamp of each user
	// and remove offline users from waiting list
	presenceUserIDs := make([]uint, len(presenceList.Items))
	for _, i := range presenceList.Items {
		presenceUserIDs = append(presenceUserIDs, i.UserID)
	}

	finlaList := make([]entity.WaitingMember, 0)
	for _, i := range waitingList {
		if funk.ContainsUInt(presenceUserIDs, i.UserID) && i.TimeStamp > timestamp.Add(-5*time.Minute) {
			finlaList = append(finlaList, i)
		} else {
			//remove from waiting list
		}
	}

	for i := 0; i < len(finlaList)-1; i = i + 2 {
		mu := entity.MatchedUsers{
			Category: category,
			UserID:   []uint{finlaList[i].UserID, finlaList[i+1].UserID},
		}

		fmt.Println("mu", mu)

		// publish a new event for mu
		// remove mu users from waiting list

	}

}
