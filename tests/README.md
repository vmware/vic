VIC Integration & Functional Test Suite
=======

To run the deprecated tests:

1. Integration tests can be run by calling `make integration-tests` from the project's root directory.

To run these tests locally:

1. Create a secrets.yml file that includes:  
```
environment:  
  ESX_URL: <IP address of your ESX server>  
  ESX_USERNAME: <username you use to login to ESX server>  
  ESX_PASSWORD: <password you use to login to ESX server>  
  ESX_RESOURCE: <resource pool, e.g. /ha-datacenter/host/localhost.localdomain/Resources>  
```
2. Execute drone from the projects root directory:

  `drone exec --trusted -E "secrets.yml" --yaml ".drone-e2e.yml"`

Find the documentation for each of the tests here:
-
###[Test Suite Documentation](test-cases/TestGroups.md)
