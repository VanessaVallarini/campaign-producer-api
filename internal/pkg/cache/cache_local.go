package cache

import (
	"sync"

	easyzap "github.com/lockp111/go-easyzap"
)

type LocalMapService struct {
	it map[string]interface{}
	mu sync.RWMutex
}

func NewLocalMapService() *LocalMapService {
	m := make(map[string]interface{})
	return &LocalMapService{
		it: m,
	}
}

func (l *LocalMapService) Set(key string, value interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	easyzap.Infof("SET key: %s value: %s", key, value)
	l.it[key] = value
}

func (l *LocalMapService) Get(key string) (interface{}, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	value, exists := l.it[key]
	if exists {
		easyzap.Infof("GET key: %s value: %v", key, value)
	} else {
		easyzap.Infof("GET key: %s not found", key)
	}
	return value, exists
}

func (l *LocalMapService) Del(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	easyzap.Infof("DEL key: %s", key)
	delete(l.it, key)
}

func (l *LocalMapService) GetAll() map[string]interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()

	copyMap := make(map[string]interface{}, len(l.it))
	for k, v := range l.it {
		copyMap[k] = v
	}
	return copyMap
}

func (l *LocalMapService) Reset() {
	easyzap.Info("RESET map")
	l.it = make(map[string]interface{})
}
