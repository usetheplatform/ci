# CI System

A back-end server of a CI system, that is built purely for educational purposes.

## Architecture Overview

It consists of 4 components, that live in separate processes:
- a repository observer
- a job dispatcher
- a test runner
- a reporting system.

It is built as a distributed system, where every component can live in a different
environment. However, for testing, it could be setup to work on a local machine.
Communication between the components is done via websockets.

### Repository Observer

Repository Observer observes events from a given repository and notifies the job dispatcher, that a new commit is available for test.

> Please note, that only git/github VCS is supported.

### Job Dispatcher

Job Dispatcher holds a queue of commits, that need to be tested.
It communicates with the registered test runners, allocating jobs for them.

### Test Runners

Test Runners run a job for a given commit and report back to the job dispatcher and to the reporting system.

### Reporting System

Reporting System collects information about the job runs and display it in a different ways.


> @TODO: Write a diagram explaining the relationship between the components
> @TODO: Write a guide how to run it