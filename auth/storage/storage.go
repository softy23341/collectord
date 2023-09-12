package storage

// Storage TBD
type Storage interface {
	Set(token, val string) error
	Get(token string) (val string, err error)
	Del(token string) error
}

