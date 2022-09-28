package exire

import "context"

type (
	StoreContract interface {
		Set(ctx context.Context, key string, value string) error
		Get(ctx context.Context, key string) (string, error)
		Delete(ctx context.Context, key string) error
	}
	LockContract interface {
		Lock(ctx context.Context) error
		Unlock(ctx context.Context) error
		Cleanup(ctx context.Context) error
	}
	TTLContract interface {
		Valid(ctx context.Context, key string) (bool, error)
		Add(ctx context.Context, key string) error
		Remove(ctx context.Context, key string) error
	}

	Exire struct {
		store      StoreContract
		perkeylock LockContract
		globallock LockContract
		ttl        TTLContract
	}

	With func(e *Exire)
)

func New(opts ...With) *Exire {
	e := &Exire{
		perkeylock: &noopLock{},
		globallock: &noopLock{},
		ttl:        &noopTTL{},
	}

	for _, o := range opts {
		o(e)
	}

	if e.store == nil {
		panic("cannot initialize a Exire without a store")
	}

	return e
}

func (e *Exire) Cleanup() {
	err := e.perkeylock.Cleanup(context.TODO())
	if err != nil {
		panic(err)
	}

	err = e.globallock.Unlock(context.TODO())
	if err != nil {
		panic(err)
	}
	err = e.globallock.Cleanup(context.TODO())
	if err != nil {
		panic(err)
	}
}

func (e *Exire) Set(ctx context.Context, key, value string) error {
	err := e.perkeylock.Lock(ctx)
	if err != nil {
		return err
	}
	defer e.perkeylock.Unlock(ctx)
	err = e.globallock.Lock(ctx)
	if err != nil {
		return err
	}

	err = e.ttl.Add(ctx, key)
	if err != nil {
		return err
	}

	err = e.store.Set(ctx, key, value)
	if err != nil {
		return err
	}

	return nil
}

func (e *Exire) Get(ctx context.Context, key string) (string, error) {
	err := e.perkeylock.Lock(ctx)
	if err != nil {
		return "", err
	}
	defer e.perkeylock.Unlock(ctx)
	err = e.globallock.Lock(ctx)
	if err != nil {
		return "", err
	}

	valid, err := e.ttl.Valid(ctx, key)
	if err != nil {
		return "", nil
	}

	if !valid {
		err = e.ttl.Remove(ctx, key)
		if err != nil {
			return "", err
		}

		return "", nil
	}

	val, err := e.store.Get(ctx, key)
	if err != nil {
		return "", err
	}

	return val, err
}

func (e *Exire) Delete(ctx context.Context, key string) error {
	err := e.perkeylock.Lock(ctx)
	if err != nil {
		return err
	}
	defer e.perkeylock.Unlock(ctx)
	err = e.globallock.Lock(ctx)
	if err != nil {
		return err
	}

	err = e.ttl.Remove(ctx, key)
	if err != nil {
		return err
	}

	err = e.store.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
