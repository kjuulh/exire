package main

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/kjuulh/exire"
)

type externalStore struct {
	mu   *sync.Mutex
	file *os.File
}

func newExternalStore() exire.StoreContract {
	file, err := os.Create("_example/with_external_lock/locks/store.txt")
	if err != nil {
		panic(err)
	}

	return &externalStore{
		mu:   &sync.Mutex{},
		file: file,
	}
}

func (es *externalStore) Set(ctx context.Context, key string, value string) error {
	return nil
}

func (es *externalStore) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (es *externalStore) Delete(ctx context.Context, key string) error {
	return nil
}

type externalLock struct {
	mu   *sync.Mutex
	file *os.File
}

func newExternalLock() exire.LockContract {
	file, err := os.Create("_example/with_external_lock/locks/lock.json")
	if err != nil {
		panic(err)
	}

	return &externalLock{
		mu:   &sync.Mutex{},
		file: file,
	}
}

func (el *externalLock) Lock(ctx context.Context) error {
	return nil
}

func (el *externalLock) Unlock(ctx context.Context) error {
	return nil
}

func (el *externalLock) Cleanup(ctx context.Context) error {
	return nil
}

// The goal is to provide a lock on an external dependency (such as yarn.lock)
// No expiry in this case
func main() {
	e := exire.New(
		exire.WithStore(newExternalStore()),
		exire.WithGlobalLock(newExternalLock()),
	)
	defer e.Cleanup()

	ctx := context.Background()

	err := e.Set(ctx, "hello", "world")
	if err != nil {
		panic(err)
	}

	_, err = e.Get(ctx, "hello")
	if err != nil {
		panic(err)
	}

	err = e.Delete(ctx, "hello")
	if err != nil {
		panic(err)
	}

	_, err = e.Get(ctx, "hello")
	if err == nil {
		panic(errors.New("should panic, because hello must not exist"))
	}
}
