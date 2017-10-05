Test 23-05 - VCH Logs
=======

# Purpose:
To verify vic-machine-server can provide logs for a VCH host when available

# References:
[1 - VIC Machine Service API Design Doc - VCH Certificate](../../../doc/design/vic-machine/service.md)

# Environment:
This test runs an external service binary that exposes the vic-machine API.

# Test Steps:
1. Verify that the creation log is unavailable before creating a VCH
2. Deloy a VCH into the test environment
3. Verify that the creation log is available after the VCH is created using the vic-machine-service
4. Verify that the creation log is available for its particular datacenter using the vic-machine-service

# Expected Outcome:
* Step 1 should error with 404 (not found) as no log file exists
* Step 3-4 should succeed and output should contain log message that the creation is completed successfully

# Possible Problems:
None
