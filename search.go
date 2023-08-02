package gitstatus

import (
	"io/fs"
	"os"
	"path/filepath"
)

func FindRepos() []string {
	dir := os.Getenv("HOME")
	if dir == "" {
		panic("No $HOME")
	}

	repos := make([]string, 0)
	isRepo := func(path string, d fs.DirEntry, err error) error {
		if path == dir+"/go" || d.Name() == "node_modules" {
			return filepath.SkipDir
		}
		if err != nil {
			return filepath.SkipDir
		}
		possipath := path + "/.git"
		file, err := os.Open(possipath)
		if err != nil {
			return nil
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			return nil
		}
		if fileInfo.IsDir() {
			repos = append(repos, path)
			return filepath.SkipDir
		}
		return nil
	}

	err := filepath.WalkDir(dir, isRepo)
	if err != nil {
		panic(err.Error())
	}
	return repos
}

// Automatically add all new repos
func CompareRepos(repos []string) {
	mutex.Lock()
	for _, v := range repos {
		if _, ok := globalRepos[v]; !ok {
			globalRepos[v] = true
		}
	}
	mutex.Unlock()
}

func DumpRepos() {
	defaultLogger.Debug("Dumping Repos")
	for k, _ := range globalRepos {
		defaultLogger.Debug(k)
	}
}

func UpdateRepos() {
	allRepos := FindRepos()
	CompareRepos(allRepos)
	WriteDataStore()
}
