package lock

// A function to generate a lock
type IdFunc func(interface{}) (lockId string, err error)

// A basic lock Id function which return a LockIdFunc which uses static lock ID supplied
func NewLockIdFunc(lockIdToUse string) IdFunc {
	return func(i interface{}) (lockId string, err error) {
		return lockIdToUse, nil
	}
}

// Lock information
type Lock struct {
	Group      string
	LockIdFunc IdFunc
}

type RunFunc func(interface{}) (out interface{}, err error)

type DistributedLock interface {
	RunInLock(lock Lock, runFunc RunFunc) (out interface{}, err error)
}
