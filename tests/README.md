# VIC Engine Integration & Functional Test Suite

To run the integration tests locally:

## Automatic with defaults

Use ./local-integration-test.sh

## Manually configure local Drone

* Create a `test.secrets` file containing secrets in KEY=VALUE format which includes:

  ```
    GITHUB_AUTOMATION_API_KEY=<token from https://github.com/settings/tokens>
    TEST_BUILD_IMAGE=""
    TEST_URL_ARRAY=<IP address of your test server>
    TEST_USERNAME=<username you use to login to test server>
    TEST_PASSWORD=<password you use to login to test server>
    TEST_RESOURCE=<resource pool, e.g. /ha-datacenter/host/localhost.localdomain/Resources>
    TEST_DATASTORE=<datastore name, e.g. datastore1>
    TEST_TIMEOUT=60s
    VIC_ESX_TEST_DATASTORE=<datastore path, e.g. /ha-datacenter/datastore/datastore1>
    VIC_ESX_TEST_URL=<user:password@IP address of your test server>
    DOMAIN=<domain for TLS cert generation, may be blank>
  ```

  If you are using a vSAN environment or non-default ESX install, then you can also specify the two networks to use with the following command (make sure to add them to the yaml file in Step 2 below as well):

  ```
    BRIDGE_NETWORK=bridge
    PUBLIC_NETWORK=public
  ```

  If you want to use an existing VCH to run a test (e.g. any of the group 1 tests) on, add the following secret to the secrets file:

  ```
    TARGET_VCH=<name of an existing VCH>
  ```

  The above TARGET_VCH is best used for tests where you do not want to exercise vic-machine's create/delete operations.  The Group 1 tests is a great example.  Their main goal is to test docker commands.

  If TARGET_VCH is not specified, and you have a group initializer and cleanup file (see the group 1 tests), there is another variable to control whether use a shared VCH.

  ```
    MULTI_VCH=<1 for enable>
  ```

  Enabling MULTI_VCH forces each suite to install a new VCH and cleans it up at the end of the test.  If the test is in 'single vch' mode, it will respect the group initializer and cleanup file.  If the initializer creates the shared VCH, then all tests will use that shared VCH.  If TARGET_VCH exist, MULTI_VCH is ignored.

  ```
    DEBUG_VCH=<1 to enable>
  ```

  Enabling DEBUG_VCH will log existing docker images and containers on a VCH at the start of a test suite.


* Execute Drone from the project root directory:

  Drone will run based on `.drone.local.yml` - defaults should be fine, edit as needed

  *  To run only the regression tests:
     ```
     drone exec --repo.trusted --secrets-file "test.secrets"  .drone.local.yml
     ```

  * To run the full suite:
     ```
		 drone exec --repo.trusted --repo.branch=master --repo.fullname="vmware/vic"  --secrets-file "test.secrets"  .drone.local.yml
     ```

## Test a specific .robot file

* Set environment in robot.sh
* Run robot.sh with the desired .robot file

  From the project root directory:
  ```
  ./tests/robot.sh tests/test-cases/Group6-VIC-Machine/6-04-Create-Basic.robot
  ```

## Find the documentation for each of the tests here:

* [Automated Test Suite Documentation](test-cases/TestGroups.md)
* [Manual Test Suite Documentation](manual-test-cases/TestGroups.md)

## Tips on running tests more efficiently

Here are some recommendations that will make running tests more effective.

1. If a group of tests do not need an independent VCH to run on, there is a facility to use a single VCH for the entire group.  The Group 1 tests utilizes this facility.  To utilize this in a group (a folder of robot files),
    - Add an __init__.robot file as the first robot file in your group.  This special init file should install the VCH and save the VCH-NAME to environment variable REUSE-VCH.  The bootrap file also needs to save the VCH to the removal exception list.
    - Every robot file should neither assume a group-wide VCH.  It should install and remove a VCH for it's own use.  This allows the single robot file to be properly targeted for testing as a single test or as part of a group of test (with group-wide VCH).  When a group wide VCH is in use, the exception list will bypass the per-robot file VCH install and removal.
    - Write individual tests within a robot file with NO assumption of a standalone VCH.  Assume a shared VCH.  This will allow the
    tests to run in either shared VCH or standalone VCH mode.
    - Add a cleanup.robot file that handles cleaning up the group-wide VCH.  It needs to remove the group-wide VCH-NAME from the cleanup exception list.
2. Write all tests within robot file with the assumption that the VCH is in shared mode.  Don't assume there are no previously created containers and images.  If a robot file needs this precondition, make sure the suite setup cleans out the VCH before running any test.
3. If there is an existing VCH available, it is possible to bypass the VCH installation/deletion by adding a TARGET_VCH into the list of test secrets.