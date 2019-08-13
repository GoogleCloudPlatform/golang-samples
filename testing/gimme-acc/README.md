## gimme-acc: Service account pool for GCS SA HMAC system test

Credits: [gimmeproj](https://github.com/GoogleCloudPlatform/golang-samples/tree/master/testing/gimmeproj)

gimme-acc manages a pool of service accounts in a project and leases one when a system test needs one.

```
Usage:
	gimme-acc -project=[meta project ID] command
	gimme-acc -project=[meta project ID] -output=list status

Commands:
	lease [duration]               Leases a service account for a given duration. Prints the service account email to stdout.
	done [service account email]   Returns a service account to the pool.
	version                        Prints the version of gimme-acc.

Administrative commands:
	pool-add [service account email]   Adds a service account to the pool.
	pool-rm  [service account email]   Removes a service account from the pool.
	status                             Displays the current status of the meta project. Respects -output.
```

### Example use in integration tests

```
set -e -x

curl https://storage.googleapis.com/gimme-acc/linux_amd64/gimme-acc > gimme-acc
chmod +x gimme-acc
./gimme-acc version

export HMAC_SERVICE_ACCOUNT=$(./gimme-acc -project gimme-acc lease 15m)
trap "./gimme-acc -project gimme-acc done $HMAC_SERVICE_ACCOUNT" EXIT

go test ....
```
