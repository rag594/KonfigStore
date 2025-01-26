package locks

import "context"

type LockManager interface {
	Lock(context.Context, string) error
	Unlock(context.Context) error
}
