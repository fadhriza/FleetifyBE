package migration

import (
	"sync"
)

var (
	seedRegistry  = make(map[string]func() interface{})
	registryMutex sync.RWMutex
)

func RegisterSeeder(modelName string, seedFunc func() interface{}) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	seedRegistry[modelName] = seedFunc
}

func getSeedDataFromRegistry(modelName string) interface{} {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	if seedFunc, exists := seedRegistry[modelName]; exists {
		return seedFunc()
	}
	return nil
}
