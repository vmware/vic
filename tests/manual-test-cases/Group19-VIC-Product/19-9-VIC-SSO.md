Test 19-9 VIC SSO
=======

# Purpose:
To verify that the SSO and session management features of VIC product works across Admiral, Harbor, and Engine

# References:
[1- VIC SSO Feature](TBD - waiting on official docs to link)

# Environment:
This test requires that a vCenter server with integrated PSC and with a discrete PSC is running and available

# Test Steps:
1. Download and install the VIC OVA from here:
[OVA Releases](https://console.cloud.google.com/storage/browser/vic-product-ova-releases/?project=eminent-nation-87317&authuser=1)
2. Navigate to the getting started page at https://<VIC>:9443
3. Login into the vCenter using incorrect credentials
4. Navigate to the admiral and demo VCH pages
5. Login into the vCenter using correct credentials
6. Navigate to the admiral and demo VCH pages
7. Wait 10 minutes
8. Navigate to the admiral and demo VCH pages
9. Wait another 50 minutes (60 minutes total)
10. Navigate to the admiral and demo VCH pages and attempt to make a change in admiral

# Expected Outcome:
* Steps 1,2,5-8 should all succeed without error
* Step 3 should result in error and prompt the user for the correct credentials again
* Step 4 should fail and the user should be redirected to the fileserver page in order to login correctly
* Step 6 and 8 should succeed and the user should not be prompted for login credentials to any of the pages
* Step 10 should fail and the user should be prompted for login credentials again

# Possible Problems:
None
