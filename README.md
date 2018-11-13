# chester [![Build Status](https://travis-ci.org/markhalonen/chester.svg?branch=master)](https://travis-ci.org/markhalonen/chester) [![Go Report Card](https://goreportcard.com/badge/github.com/markhalonen/chester)](https://goreportcard.com/report/github.com/markhalonen/chester) [![Coverage Status](https://coveralls.io/repos/github/markhalonen/chester/badge.svg?branch=master)](https://coveralls.io/github/markhalonen/chester?branch=master)
---
A low-effort testing framework.

chester works by creating snapshot tests against any command line output. REST APIs or anything that can be invoked from command line can use chester to watch for changes and easily update tests with one click. You can use any language that can be called from the command line (aka any language).

### Install
1. Download the latest release from the [Github Releases Tab](https://github.com/markhalonen/chester/releases)

### Usage - Minimal Example
`./chester init` Creates the `__chester__` directory where all the commands and snapshots will be stored

`./chester create 'echo "Hello world"'` Create your first chester test

`./chester test` Run your new chester test

### Test a REST API with Python
`./chester init`

`mkdir my_test` We will use this `my_test` folder to create the test. It must contain `command.sh`.

Create these files in `my_test`:
```
my_test/
├── command.sh
└── test.py
```

Where `command.sh` contains:
```bash
python test.py
``` 
and `test.py` contains:
```python
import json
import random

# Simulated server response
def json_endpoint():
    return json.dumps({"name": "Mark Halonen", "age": 23, "timestamp": random.randint(0,10000)})


# Call our "endpoint"
response = json_endpoint()

# timestamp is not expected to remain the same, so let's remove it
response_json = json.loads(response)
del response_json["timestamp"]

# Now we print so that chester can capture the output
print json.dumps(response_json, indent=4)
```

You should be able to run `./command.sh` and see json output. You may have to `chmod 777 command.sh`

Now, let's use chester to create a test from `my_test`

`chester create my_test` Will run the test and confirm the output. Select create.

`chester test` Can then be used to run this new test that uses Python to ignore the timestamp field, because we expect it to change.

## Develop Chester
1. Follow [golang.org](https://golang.org/doc/install) to get your go environment up and running
2. `go get https://github.com/markhalonen/chester`
3. `cd` into the cloned repo
4. `go install` after you make changes

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

chester create 'curl -X GET http:localhost:8080/todos' // Runs the command and displays the result. Y/N to save to disk.
chester test // Runs all the tests. For each failing test, show a diff of the output and Y/N if the snapshot should be updated.
