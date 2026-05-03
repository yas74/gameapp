package redispresence

import (
	"context"
	"gocasts/gameapp/pkg/richerror"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (d DB) Upsert(ctx context.Context, key string, timeStamp int64,
	expTime time.Duration) error {

	const op = "redispresence.Upsert"

	_, err := d.adaptor.Client().Set(ctx, key, timeStamp, expTime).Result()
	if err != nil {
		richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d DB) GetPresence(ctx context.Context, key string) (int64, error) {
	const op = richerror.Op("redispresence.GetPresence")

	timeStampString, err := d.adaptor.Client().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, richerror.New(op).WithErr(err).WithKind(richerror.KindNotFound)
		}

		return 0, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	timeStamp, _ := strconv.Atoi(timeStampString)

	return int64(timeStamp), nil
}
