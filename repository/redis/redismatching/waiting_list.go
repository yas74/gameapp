package redismatching

import (
	"context"
	"fmt"
	"gocasts/gameapp/entity"
	"gocasts/gameapp/pkg/richerror"
	"gocasts/gameapp/pkg/timestamp"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// TODO - add to config in usecase layer
const WaitingListPrefix = "waitinglist"

func (d DB) AddToWaitingList(userID uint, category entity.Category) error {
	const op = richerror.Op("redismatching.AddToWaitingList")

	ctx := context.Background()
	_, err := d.adapter.Client().ZAdd(ctx, getCategoryKey(category), redis.Z{
		Score:  float64(timestamp.Now()),
		Member: fmt.Sprintf("%d", userID),
	}).Result()

	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil

}

func (d DB) GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error) {
	const op = richerror.Op("redismatching.GetWaitingListByCategory")

	min := fmt.Sprintf("%d", timestamp.Add(-2*time.Hour))
	max := fmt.Sprintf("%d", timestamp.Now())
	list, err := d.adapter.Client().ZRangeByScoreWithScores(ctx, getCategoryKey(category),
		&redis.ZRangeBy{
			Min:    min,
			Max:    max,
			Offset: 0,
			Count:  0,
		}).Result()

	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	result := []entity.WaitingMember{}

	for _, item := range list {
		userID, _ := strconv.Atoi(item.Member.(string))

		result = append(result, entity.WaitingMember{
			UserID:    uint(userID),
			TimeStamp: int64(item.Score),
			Category:  category,
		})
	}

	return result, nil

}

func getCategoryKey(category entity.Category) string {

	return fmt.Sprintf("%s:%s", WaitingListPrefix, category)
}
