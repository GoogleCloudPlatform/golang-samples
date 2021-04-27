# Sample port TODOs

## Milestone 0: Getting Started

**Time Spent:** 15 minutes

* [x] Find issue/review notes
* [x] Orientation & Git Wrangling

## Milestone 1: Pets topic MVP

**Time Spent:** 50 minutes

* [x] Add Cloud Run Button integration
* [x] Update SQL queries/table setup
* [x] Switch to lib/pq since we don't need non-module compat.
* [x] Update HTML template
* [ ] Use bash scripts to deploy & manual test

Change Flags for PR:

* Add `set -eu` to setup.sh and postcreate.sh
  * -e should always be used in a script
  * -u errors on unset variables, an alt to specific checks.


## Milestone 2: Service Upgrade

* [ ] Add graceful termination to server
* [ ] Create leveled logger with tracing via logrus

## Milestone 3: Add Login flow

* [ ] Extract & validate UID
* [ ] Add login form to UI
* [ ] Use Firebase Admin SDK to validate JWT
* [ ] Test IDP working

## Milestone 4: Add tests

* [ ] Test evaluation
* [ ] Test writing

## Milestone 5: Cleanup

* [ ] Load secrets from env, no library
* [ ] Parameterize db table
* [ ] Fix up region tags
* [ ] Copyright review
* [ ] Final run-through with Cloud Run Button
