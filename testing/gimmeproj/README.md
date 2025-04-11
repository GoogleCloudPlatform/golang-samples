## gimmeproj

gimmeproj manages a pool of projects and leases to those projects.

The meta project (specified by the `-project` flag) stores the metadata for the pool.

```
Usage:
  gimmeproj -project=[meta project ID] command

Commands:
  lease [duration]    Leases a project for a given duration. Prints the project ID to stdout.
  done [project ID]   Returns a project to the pool.

Administrative commands:
  pool-add [project ID]       Adds a project to the pool.
  pool-rm  [project ID]       Removes a project from the pool.
  status                      Displays the current status of the meta project.
```

### Example use in integration tests

```
set -e -x

curl https://storage.googleapis.com/gimme-proj/linux_amd64/gimmeproj > gimmeproj
chmod +x gimmeproj
./gimmeproj version

export TEST_PROJECT=$(./gimmeproj -project meta-project lease 15m)
trap "./gimmeproj -project meta-project done $TEST_PROJECT" EXIT

go test ....
```
