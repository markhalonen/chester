# snapper [![Build Status](https://travis-ci.org/markhalonen/snapper.svg?branch=master)](https://travis-ci.org/markhalonen/snapper) [![Go Report Card](https://goreportcard.com/badge/github.com/markhalonen/snapper)](https://goreportcard.com/report/github.com/markhalonen/snapper) [![Coverage Status](https://coveralls.io/repos/github/markhalonen/snapper/badge.svg?branch=master)](https://coveralls.io/github/markhalonen/snapper?branch=master)
---
Snapper aims to be a low-effort testing framework. Life is too short to write tests manually.

Snapper works by creating snapshot tests against any command line output. REST APIs or anything that can be invoked from command line can use snapper to watch for changes and easily update tests with one click.

### Install
1. Download the latest release from the [Github Releases Tab](https://github.com/markhalonen/snapper/releases)

### Usage
`./snapper init` Creates the `__snapper__` directory where all the commands and snapshots will be stored

`./snapper create 'echo "Hello world"'` Create your first snapper test

`./snapper test` Run your new snapper test

## Motivation
API Snapshot Testing Tool

"Snapshot Testing" is a type of testing where you basically just watch for changes. You take a "snapshot" of the response from a system and save it. Then, you continually test the system by comparing the system response to the snapshot saved on disk. So it really is just watching for changes.

The benefits of Snapshot Testing:
1. Trivial to "write tests"
2. See #1. Writing tests manually is so boring.

There should be very limited code involved in "writing" snapshot tests. But that's currently not the case. Developers are expected to manually write the same snapshot code over and over.

So I propose a new API Snapshot Testing Tool with the following goals:

1. Works for any system that gives output from a command line instruction (it can test itself!)
2. Minimize writing code
3. Quickly create snapshot tests

   3.1. Specify the command
   
   3.2. Inspect the response
   
   3.3. Save snapshot
   
4. Quickly compare failing snapshot tests and update so they pass again

The command line tool would look something like:

snapper create 'curl -X GET http:localhost:8080/todos' // Runs the command and displays the result. Y/N to save to disk.
snapper test // Runs all the tests. For each failing test, show a diff of the output and Y/N if the snapshot should be updated.
