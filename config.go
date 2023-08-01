package gitstatus

import (
	"encoding/json"
	"os"
	"sync"
)

var globalRepos map[string]bool
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
		globalRepos = make(map[string]bool)
		return
	}
	defer repoStore.Close()
	reposDecoder := json.NewDecoder(repoStore)
	if err != nil {
		panic(err.Error())
	}
	err = reposDecoder.Decode(&globalRepos)
	if err != nil {
		panic(err.Error())
	}
}

func WriteDataStore() {
	mutex.Lock()
	defer mutex.Unlock()
	repoStore, err := os.Create(dataStore)
	if err != nil {
		panic(err.Error())
	}
	defer repoStore.Close()
	reposEncoder := json.NewEncoder(repoStore)
	if err != nil {
		panic(err.Error())
	}
	reposEncoder.SetIndent(" ", "\t")
	err = reposEncoder.Encode(&globalRepos)
	if err != nil {
		panic(err.Error())
	}
}
