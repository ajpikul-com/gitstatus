package gitstatus

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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
				defaultLogger.Debug("Not Clean")
				defaultLogger.Debug(status.String())
			}
			remotes, err := repo.Remotes()
			if err != nil {
				defaultLogger.Error(err.Error()) // error with repo.Config, probalby file system issue
				continue
			}
			if len(remotes) == 0 {
				defaultLogger.Debug("This repo has no remote!")
			} else {
				for _, remote := range remotes {
					err = remote.Fetch(&git.FetchOptions{})
					if err != nil && err == git.NoErrAlreadyUpToDate {
						defaultLogger.Error(err.Error())
						continue
					}
					if remote.Config().Name == "origin" {
						upstreamRef, err := repo.Reference(plumbing.NewRemoteHEADReferenceName("origin"), true)
						if err != nil {
							panic(err.Error())
						}
						options := &git.LogOptions{
							From:  upstreamRef.Hash(),
							Order: git.LogOrderCommitterTime,
						}
						iter, err := repo.Log(options)
						if err != nil {
							panic(err.Error())
						}
						err = iter.ForEach(func(commit *object.Commit) error {
							defaultLogger.Debug(commit.Hash.String()) // check to see if we are that house
							defaultLogger.Debug(commit.Message)
							return nil
						})
						if err != nil {
							panic(err.Error())
						}
					}
				}
				// It has remotes and we've fetched them
				// We now want to know if we're behind from upstream
				// But go-git provides no easy way to do that
				// We'd have to get a reference to the upstream head
			}
		}
	}
}
