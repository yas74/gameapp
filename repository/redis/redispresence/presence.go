package redispresence

import (
	"context"
	"gocasts/gameapp/pkg/richerror"
	"time"
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
