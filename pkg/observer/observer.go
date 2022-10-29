package observer

import (
	"fmt"
	"os"

	. "github.com/usetheplatform/ci-system/pkg/common"
	"github.com/usetheplatform/ci-system/pkg/connection"
	"gopkg.in/src-d/go-git.v4"
)

type RepoObserver struct {
	config     *Options
	repo       *git.Repository
	store      Store
	connection connection.Connection
}

func NewRepoObserver(config *Options) RepoObserver {
	repo, err := cloneRepo(config.Path)
	CheckIfError(err)

	store := NewStore(".commit-sha") // @TODO: Make it dynamic
	CheckIfError(err)

	connection := connection.Connect(config.DispatcherAddress)

	return RepoObserver{
		repo:       repo,
		store:      store,
		config:     config,
		connection: connection,
	}
}

func (observer *RepoObserver) DetectChanges() {
	cachedCommitHash, err := observer.store.Read()
	CheckIfError(err)

	// @CLARIFY: poll is called infinitely and we want to make sure that
	// it doesn't pick up already tested commit
	err = observer.store.Clear()
	CheckIfError(err)

	observer.sync()
	CheckIfError(err)

	ref, err := observer.repo.Head()
	CheckIfError(err)

	// Detect changes

	currentCommit, err := observer.repo.CommitObject(ref.Hash())
	CheckIfError(err)
	currentCommitHash := currentCommit.Hash.String()

	if cachedCommitHash == nil || *cachedCommitHash != currentCommitHash {
		Info("Found new changes...%v", currentCommitHash)
		observer.store.Write(currentCommitHash)
	}

	if exists, err := observer.store.Exists(); err == nil && exists == true {
		// notify dispatcher by asking for status
		err = observer.connection.Notify("status")
		CheckIfError(err)
		// notify dispatcher: send dispatcher address & payload - repo and commit
		err = observer.connection.Notify(fmt.Sprintf("dispatch:%s", currentCommitHash))
		CheckIfError(err)
		Info("%s has been dispatched", currentCommitHash)
	}
	CheckIfError(err)
}

func (observer *RepoObserver) Destroy() {
	observer.connection.Close()
	// destroy poller, repo and close connection
}

func (observer *RepoObserver) sync() error {
	workTree, err := observer.repo.Worktree()
	CheckIfError(err)

	Info("git reset --hard HEAD")
	err = workTree.Reset(&git.ResetOptions{
		Mode: git.HardReset,
	})
	CheckIfError(err)

	Info("git pull origin")

	// @TODO: Add support for different branches
	return workTree.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
}

func cloneRepo(repoPath string) (*git.Repository, error) {
	Info("git clone %v", repoPath)

	clonePath := "/tmp/observee" // @TODO: Make it dynamic

	repo, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:      repoPath,
		Progress: os.Stdout,
	})

	if err == nil {
		Info("Repository %s has been cloned to %s", repoPath, clonePath)
	}

	return repo, err
}
