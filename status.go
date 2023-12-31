package gitstatus

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	git2go "gopkg.in/libgit2/git2go.v24"
)

type StateMap struct {
	states map[string]map[string]RepoState
	mutex  sync.Mutex
}

func NewStateMap() *StateMap {
	return &StateMap{states: make(map[string]map[string]RepoState)} // TODO
}

func (sm *StateMap) ClearClient(client string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if _, ok := sm.states[client]; ok {
		delete(sm.states, client)
	}
}

func (sm *StateMap) Update(client string, repoName string, state RepoState) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if _, ok := sm.states[client]; !ok {
		sm.states[client] = make(map[string]RepoState)
	}
	state.Client = client
	sm.states[client][repoName] = state
}
func (sm *StateMap) MarshalJSON() ([]byte, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return json.Marshal(sm.states)
}

type RepoState struct {
	Name   string
	Remote bool
	Dirty  string
	Ahead  int
	Behind int
	send   bool
	Client string
}

func VerifyRepos() map[string]RepoState { // From here, it's probably time to send them over to the server
	repoStates := make(map[string]RepoState)
	mutex.Lock()
	defer WriteDataStore()
	defer mutex.Unlock() // Last in first out!
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
			newState := RepoState{Name: k}
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
			for k, _ := range status {
				if strings.Contains(k, "node_modules") {
					delete(status, k)
				}
			}
			if !status.IsClean() {
				defaultLogger.Debug("Not Clean")
				defaultLogger.Debug(status.String())
				newState.Dirty = status.String()
				newState.send = true
			}
			remotes, err := g2gRepo.Remotes.List()
			if err != nil {
				defaultLogger.Error("Error getting remotes list")
				defaultLogger.Error(err.Error())
				panic(err.Error())
			}
			if len(remotes) == 0 {
				defaultLogger.Debug("This repo has no remote!")
				newState.Remote = false
				newState.send = true
			} else {
				defaultLogger.Debug("Setting repo states Remote to true")
				newState.Remote = true
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
						if headRef.IsBranch() {
							checkedOutBranch := headRef.Branch()
							upstreamRef, err := checkedOutBranch.Upstream()
							if err != nil {
								defaultLogger.Error(err.Error())
								// Likely not found
								newState.Remote = false
								newState.send = true // There are remotes, but no upstream!
								continue
							}
							ahead, behind, err := g2gRepo.AheadBehind(headRef.Target(), upstreamRef.Target())
							newState.Ahead = ahead
							newState.Behind = behind
							if ahead != 0 || behind != 0 {
								newState.send = true
							}
							if err != nil {
								defaultLogger.Error(err.Error())
								panic(err.Error())
								continue
							}
							defaultLogger.Debug(strconv.Itoa(ahead))
							defaultLogger.Debug(strconv.Itoa(behind))
						}
					}
				}
			}
			if newState.send {
				repoStates[newState.Name] = newState
			}
		}
	}
	return repoStates
}
