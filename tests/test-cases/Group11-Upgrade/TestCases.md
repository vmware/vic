Group 11 - Upgrade
=======

#Purpose:
To test vic-machine upgrade

[Test 11-01 - Upgrade](11-01-Upgrade.md)

#Environment:
Set up a VCH using an old, known version from bintray

#Test Steps:
1. Verify that Upgrade is present
2. Perform an upgrade with a short timeout to test the automated rollback from failed upgrades
3. Perform an upgrade and check that it is successful
4. Manually roll back to the last version and check that it is successful 
5. Upgrade again to make sure that rollback is also reversible.

#Expected Outome:
1. Upgrade is present
2. Upgrade times out and rolls back to the pre-upgraded, functioning VCH
3. Upgrade, rollback, upgrade again routine works as expected
-
