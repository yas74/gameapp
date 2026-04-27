package redismatching

import (
	"context"
	"fmt"
	"gocasts/gameapp/entity"
	"gocasts/gameapp/pkg/richerror"
	"time"

	"github.com/redis/go-redis/v9"
)

// TODO - add to config in usecase layer
const WaitingListPrefix = "waitinglist"

func (d DB) AddToWaitingList(userID uint, category entity.Category) error {
	const op = richerror.Op("redismatching.AddToWaitingList")

	ctx := context.Background()
	zsetKey := fmt.Sprintf("%s:%s", WaitingListPrefix, category)
	_, err := d.adapter.Client().ZAdd(ctx, zsetKey, redis.Z{
		Score:  float64(time.Now().UnixMicro()),
		Member: fmt.Sprintf("%d", userID),
	}).Result()

	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil

}
