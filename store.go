package main

import (
	"context"
	"log"
	"sync"

	"github.com/bjjb/urleen/base62"
	"github.com/go-redis/redis/v8"
)

type store interface {
	put(string, string) string
	get(string) string
}

type redisStore struct {
	client       *redis.Client
	counter, url string
}

func (r *redisStore) open() {
	if r.client != nil {
		return
	}
	opts, err := redis.ParseURL(r.url)
	if err != nil {
		log.Fatal(err)
	}
	r.client = redis.NewClient(opts)
}

func (r *redisStore) ping() error {
	r.open()
	ctx := context.Background()
	return r.client.Ping(ctx).Err()
}

func (r *redisStore) put(id, s string) string {
	r.open()
	ctx := context.Background()
	if id == "" {
		id = base62.Encode(uint64(r.client.Incr(ctx, r.counter).Val()))
	}
	r.client.Set(ctx, id, s, 0)
	return id
}

func (r *redisStore) get(id string) string {
	r.open()
	ctx := context.Background()
	s, _ := r.client.Get(ctx, id).Result()
	return s
}

type mapStore struct {
	counter uint64
	data    map[string]string
	mutex   *sync.Mutex
}

func (m *mapStore) put(id, s string) string {
	m.open()
	m.mutex.Lock()
	if id == "" {
		id = base62.Encode(m.counter)
	}
	m.counter++
	m.mutex.Unlock()
	m.data[id] = s
	return id
}

func (m *mapStore) get(id string) string {
	m.open()
	return m.data[id]
}

func (m *mapStore) open() {
	if m.mutex == nil {
		m.mutex = new(sync.Mutex)
	}
	if m.data == nil {
		m.data = make(map[string]string)
	}
}
