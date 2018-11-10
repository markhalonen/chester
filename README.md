# snapper


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
  4.1 Specify the query
  4.2 Inspect the result
  4.3 Save snapshot
4. Quickly compare failing snapshot tests and update so they pass again
  4.1 Update in bulk

My first thought was a UI, but in reality a command line tool is probably best.

The command line tool would look something like:

snapper create 'curl -X GET http:localhost:8080/todos' // Runs the query and displays the result. Y/N to save to disk. curl command is the key for the db (db probably json)
snapper test // Runs all the tests. For each failing test, show a diff of the output and Y/N if the snapshot should be updated.
