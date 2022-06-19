#### Related Issues
Fixes:   #  
Towards: #  

<!--
Thank you for submitting a pull request!

Here's a checklist you might find useful.
[  ] There is an associated issue that is labelled
[  ] Code is up-to-date with the `master` branch
[  ] You've successfully run `make test` locally
[  ] There are new or updated unit tests validating the change

Refer to CONTRIBUTING.MD for more details.
  https://github.com/vmware/vic/blob/master/.github/CONTRIBUTING.md
-->


<details><summary>CI Options</summary><p>

<!--  DO NOT EDIT BELOW. The markdown below is used to provide some coarse options for how CI runs against this PR -->
------
- [ ] <!-- directive:fast-fail --> fail fast on error 
- [ ] <!-- directive:ops-user --> use ops-user instead of admin
- [ ] <!-- directive:parallel-jobs-inline --> parallel-jobs: `6` - number of parallel test jobs to use
- [ ] <!-- directive:shared-datastore-inline --> shared-datastore: `vsanDatastore` - name of a shared datastore to use
------
<!-- Default: run regression bucket only -->
- [ ] <!-- directive:skip-unit --> skip unit tests 
- [ ] <!-- directive:focused-unit --> focused unit tests 
------
<!-- Default: run all unit tests -->
- [ ] <!-- directive:skip-integration --> skip integration tests
- [ ] <!-- directive:all-integration --> all integration tests
- [ ] <!-- directive:specific-integration-begin --> specific integration tests:
```
Group1-Docker-Commands
Group0-Bugs/Group0-Bugs.4817
```
<!-- directive:specific-integration-end -->
------
<!-- Default: do not run scenario -->
- [ ] <!-- directive:all-scenario --> all scenario tests (slow!)
- [ ] <!-- directive:dirty-scenario --> scenario tests use existing testbed if present
- [ ] <!-- directive:specific-scenario-begin --> specific scenario tests:
```
Group5-Functional-Tests/5-11-Multiple-Cluster
```
<!-- directive:specific-scenario-end -->
------

</p></details>
<p/>
