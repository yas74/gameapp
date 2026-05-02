package redispresence

import "gocasts/gameapp/adapter/redis"

type DB struct {
	adaptor redis.Adapter
}

func New(adaptor redis.Adapter) DB {
	return DB{adaptor: adaptor}
}
