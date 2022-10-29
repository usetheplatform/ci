package main

import (
	"time"

	. "github.com/usetheplatform/ci-system/pkg/common"
	"github.com/usetheplatform/ci-system/pkg/observer"
)

func main() {
	options, err := observer.GetOptions()
	CheckIfError(err)

	repoObserver := observer.NewRepoObserver(options)

	// start polling in a separate goroutine

	// @TODO: Replace it with github events
	_ = NewPoller(func() {
		repoObserver.DetectChanges()
	}, time.Minute*5)

}
