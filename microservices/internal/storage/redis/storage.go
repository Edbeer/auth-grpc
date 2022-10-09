package redstorage

import "github.com/go-redis/redis/v9"

type Storage struct {
	Session *sessionStorage
}

func NewStorage(redis *redis.Client) *Storage {
	return &Storage{
		Session: newSessionStorage(redis),
	}
}