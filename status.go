package gitstatus

import (
	"os/exec"
	"strconv"

	"github.com/go-git/go-git/v5"
	git2go "gopkg.in/libgit2/git2go.v24"
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
			g2gRepo, err := git2go.OpenRepository(k)
			if err != nil {
				defaultLogger.Error(err.Error())
				panic(err.Error())
				continue
			}
			workTree, err := repo.Worktree()
			if err != nil {
				defaultLogger.Debug("No worktree!")
				defaultLogger.Error(err.Error())
				panic(err.Error())
				continue
			}
			status, err := workTree.Status()
			if err != nil {
				defaultLogger.Error("Error getting status")
				defaultLogger.Error(err.Error())
				panic(err.Error())
			}
			if !status.IsClean() {
				defaultLogger.Debug("Not Clean")
				defaultLogger.Debug(status.String())
			}
			remotes, err := g2gRepo.Remotes.List()
			if err != nil {
				defaultLogger.Error("Error getting remotes list")
				defaultLogger.Error(err.Error())
				panic(err.Error())
			}
			if len(remotes) == 0 {
				defaultLogger.Debug("This repo has no remote!")
			} else {
				// There's remotes, lets update them
				for _, remote := range remotes {
					defaultLogger.Debug("Found a remote")
					defaultLogger.Debug(remote)

					cmd := exec.Command("git", "-C", k, "fetch", "--all")
					_, err = cmd.Output()
					if err != nil {
						defaultLogger.Error(err.Error())
						panic(err.Error())
						continue
					}
					if remote == "origin" {
						headRef, err := g2gRepo.Head()
						if err != nil {
							defaultLogger.Error(err.Error())
							panic(err.Error())
							continue
						}
						checkedOutBranch := headRef.Branch()
						upstreamRef, err := checkedOutBranch.Upstream()
						if err != nil {
							defaultLogger.Error(err.Error())
							panic(err.Error())
							continue
						}
						ahead, behind, err := g2gRepo.AheadBehind(headRef.Target(), upstreamRef.Target())
						if err != nil {
							defaultLogger.Error(err.Error())
							panic(err.Error())
							continue
						}
						defaultLogger.Debug(strconv.Itoa(ahead))
						defaultLogger.Debug(strconv.Itoa(behind))
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
