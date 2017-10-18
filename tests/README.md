# VIC Engine Integration & Functional Test Suite

To run the integration tests locally:

## Automatic with defaults

Use ./local-integration-test.sh

## Manually configure local Drone

* Create a `test.secrets` file containing secrets in KEY=VALUE format which includes:

  ```
    GITHUB_AUTOMATION_API_KEY=<token from https://github.com/settings/tokens>
    TEST_BUILD_IMAGE=""
    TEST_URL=<IP address of your test server>
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
