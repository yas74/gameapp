package scheduler

import (
	"fmt"
	"gocasts/gameapp/dto"
	"gocasts/gameapp/service/matchingservice"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Config struct {
	MatchWaitadUsersIntervalInSeconds int `koanf:"match_waitad_users_interval_in_seconds"`
}

type Scheduler struct {
	config      Config
	sch         gocron.Scheduler
	matchingSvc matchingservice.Service
}

func New(config Config, matchingSvc matchingservice.Service) Scheduler {
	sch, err := gocron.NewScheduler()
	if err != nil {
		fmt.Println("scheduler error", err)
	}
	return Scheduler{
		sch:         sch,
		matchingSvc: matchingSvc,
		config:      config,
	}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("scheduler started")

	_, err := s.sch.NewJob(gocron.DurationJob(time.Duration(s.config.MatchWaitadUsersIntervalInSeconds)),
		gocron.NewTask(s.matchWaitadUsers))
	if err != nil {
		fmt.Println("job scheduler error:", err)
	}

	s.sch.Start()

	<-done
	fmt.Println("scheduler stopped!")
	s.sch.StopJobs()

}

func (s Scheduler) matchWaitadUsers() {

	fmt.Println(time.Now())

	req := dto.MatchWaitedUsersRequest{}
	s.matchingSvc.MatchWaitedUsers(req)
}
