VIC Integration & Functional Test Suite
=======

To run these tests locally:

1. Create a secrets.yml file that includes:  
```
environment:  
  VIC_ESX_TEST_URL: <user:password@IP address of your test server>  
  TEST_URL: <IP address of your test server>  
  TEST_USERNAME: <username you use to login to test server>  
  TEST_PASSWORD: <password you use to login to test server>  
  TEST_RESOURCE: <resource pool, e.g. /ha-datacenter/host/localhost.localdomain/Resources>  
```
2. Execute drone from the projects root directory:

  `drone exec --trusted -E "secrets.yml" --yaml ".drone-e2e.yml"`

Find the documentation for each of the tests here:
-
###[Test Suite Documentation](test-cases/TestGroups.md)
