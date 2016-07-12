# VIC Integration & Functional Test Suite

To run the integration tests locally:

1. Create a `test_secrets.yml` file that includes:

  ```
  environment:
    GITHUB_AUTOMATION_API_KEY: <token from https://github.com/settings/tokens>
    TEST_BUILD_IMAGE:
    TEST_URL_ARRAY: <IP address of your test server>
    TEST_USERNAME: <username you use to login to test server>
    TEST_PASSWORD: <password you use to login to test server>
    TEST_RESOURCE: <resource pool, e.g. /ha-datacenter/host/localhost.localdomain/Resources>
    TEST_DATASTORE: <datastore name, e.g. datastore1>
    TEST_TIMEOUT: 60s
    VIC_ESX_TEST_DATASTORE: <datastore path, e.g. /ha-datacenter/datastore/datastore1>
    VIC_ESX_TEST_URL: <user:password@IP address of your test server>
  ```

2. Create a `.drone.local.yml` file that includes:

  ```
  ---
  clone:
    path: github.com/vmware/vic
    tags: true

  build:
    integration-test:
      image: $${TEST_BUILD_IMAGE=vmware-docker-ci-repo.bintray.io/integration/vic-test:1.1}
      pull: true
      environment:
        BIN: bin
        GOPATH: /drone
        SHELL: /bin/bash
        DOCKER_API_VERSION: "1.21"
        VIC_ESX_TEST_URL: $$VIC_ESX_TEST_URL
        LOG_TEMP_DIR: install-logs
        DRONE_SERVER:  $$DRONE_SERVER
        GITHUB_AUTOMATION_API_KEY:  $$GITHUB_AUTOMATION_API_KEY
        DRONE_TOKEN:  $$DRONE_TOKEN
        TEST_URL_ARRAY:  $$TEST_URL_ARRAY
        TEST_USERNAME:  $$TEST_USERNAME
        TEST_PASSWORD:  $$TEST_PASSWORD
        TEST_DATASTORE: $$TEST_DATASTORE
        TEST_TIMEOUT: $$TEST_TIMEOUT
        GOVC_INSECURE: true
        GOVC_USERNAME:  $$TEST_USERNAME
        GOVC_PASSWORD:  $$TEST_PASSWORD
        GOVC_RESOURCE_POOL:  $$TEST_RESOURCE
        GOVC_DATASTORE: $$TEST_DATASTORE
        GS_PROJECT_ID: $$GS_PROJECT_ID
        GS_CLIENT_EMAIL: $$GS_CLIENT_EMAIL
        GS_PRIVATE_KEY: $$GS_PRIVATE_KEY
      commands:
        - tests/integration-test.sh
        #- pybot tests/test-cases/Group1-Docker-Commands/1-5-Docker-Start.robot
  ```

3. Execute drone from the projects root directory:

  `drone exec --trusted -E "test_secrets.yml" --yaml ".drone.local.yml"`


## Find the documentation for each of the tests here:

* [Automated Test Suite Documentation](test-cases/TestGroups.md)
* [Manual Test Suite Documentation](manual-test-cases/TestGroups.md)
