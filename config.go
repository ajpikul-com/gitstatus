package gitstatus

import (
	"encoding/json"
	"os"
	"sync"
)

// b, err := json.Marshal(instance of config)
var repos map[string]bool
var dataStore string
var mutex sync.Mutex

func InitDataStore(ds string) {
	dataStore = ds
	readDataStore()
}

func readDataStore() {
	mutex.Lock()
	defer mutex.Unlock()
	repoStore, err := os.Open(dataStore)
	if err != nil {
		panic(err.Error())
	}
	defer repoStore.Close()
	reposDecoder := json.NewDecoder(repoStore)
	if err != nil {
		panic(err.Error())
	}
	err = reposDecoder.Decode(&repos)
	if err != nil {
		panic(err.Error())
	}
}

func WriteDataStore() {
	mutex.Lock()
	defer mutex.Unlock()
	repoStore, err := os.Open(dataStore)
	if err != nil {
		panic(err.Error())
	}
	defer repoStore.Close()
	reposEncoder := json.NewEncoder(repoStore)
	if err != nil {
		panic(err.Error())
	}
	err = reposEncoder.Encode(&repos)
	if err != nil {
		panic(err.Error())
	}
}
