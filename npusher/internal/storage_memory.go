package internal

import "sync"

func init() {
	registerTokenStorage("memory", newMemoryTokenStorage)
}

// MemoryStorage TBD
type memoryTokenStorage struct {
	mutex  *sync.Mutex
	arnMap map[string]string
}

// newMemoryTokenStorage TBD
func newMemoryTokenStorage(_ *tokenStorageCtx) (tokenStorager, error) {
	return &memoryTokenStorage{
		mutex:  &sync.Mutex{},
		arnMap: make(map[string]string),
	}, nil
}

// Get TBD
func (ms *memoryTokenStorage) get(deviceToken string) (arn string, err error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	arn, _ = ms.arnMap[deviceToken]
	return arn, nil
}

// Set TBD
func (ms *memoryTokenStorage) set(deviceToken string, arn string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.arnMap[deviceToken] = arn
	return nil
}
