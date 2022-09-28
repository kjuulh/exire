package exire

func WithStore(s StoreContract) With {
	return func(e *Exire) {
		e.store = s
	}
}

func WithGlobalLock(l LockContract) With {
	return func(e *Exire) {
		e.globallock = l
	}
}

func WithPerkeyLock(l LockContract) With {
	return func(e *Exire) {
		e.perkeylock = l
	}
}

func WithTTL(ttl TTLContract) With {
	return func(e *Exire) {
		e.ttl = ttl
	}
}
