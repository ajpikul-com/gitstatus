package gitstatus

import (
	"github.com/go-git/go-git/v5"
)

type repoState struct {
}

func VerifyRepos() {
	for k, v := range globalRepos {
		if v {
			repo, err := git.PlainOpen(k)
			if err == git.ErrRepositoryNotExists {
				delete(globalRepos, k)
				continue
			} else if err != nil {
				defaultLogger.Error(err.Error())
				continue
			}
			defaultLogger.Debug("We found a good repo")
			defaultLogger.Debug(k)
			workTree, err := repo.Worktree()
			if err != nil {
				defaultLogger.Debug("No worktree!")
				continue
			}
			status, err := workTree.Status()
			if err != nil {
				defaultLogger.Error("Error getting status")
			}
			if !status.IsClean() {
				defaultLogger.Debug("There are changes here to be commited")
			}
			remotes, err := repo.Remotes()
			if err != nil {
				defaultLogger.Error(err.Error()) // error with repo.Config, probalby file system issue
				continue
			}
			if len(remotes) == 0 {
				defaultLogger.Debug("This repo has no remote!")
			} // else fetch
		}
	}
}
