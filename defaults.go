package exire

import "context"

type noopStore struct {
}

func newNoopStore() StoreContract {
	return &noopStore{}
}

func (es *noopStore) Set(ctx context.Context, key string, value string) error {
	return nil
}

func (es *noopStore) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (es *noopStore) Delete(ctx context.Context, key string) error {
	return nil
}

type noopLock struct {
}

func newNoopLock() LockContract {
	return &noopLock{}
}

func (el *noopLock) Lock(ctx context.Context) error {
	return nil
}

func (el *noopLock) Unlock(ctx context.Context) error {
	return nil
}

func (el *noopLock) Cleanup(ctx context.Context) error {
	return nil
}

type noopTTL struct {
}

func newNoopTTL() TTLContract {
	return &noopTTL{}
}

func (*noopTTL) Add(ctx context.Context, key string) error {
	return nil
}

func (*noopTTL) Remove(ctx context.Context, key string) error {
	return nil
}

func (*noopTTL) Valid(ctx context.Context, key string) (bool, error) {
	return true, nil
}
