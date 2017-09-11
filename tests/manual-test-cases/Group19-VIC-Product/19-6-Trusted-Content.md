Test 19-6 Trusted Content
=======

# Purpose:
To verify that the trusted content feature of VIC product works across Admiral, Harbor, and Engine

# References:
[1- VIC Trusted Content Feature](TBD - waiting on official docs to link)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Install VIC OVA into the vSphere server
2. Register Harbor as a registry with Admiral
3. Populate the harbor instance with images that can be pulled
4. Install a VCH into the same VC environment via vic-machine with the harbor instance added as a secure registry but only docker hub whitelisted
5. Do a docker pull from docker hub that should succeed
6. Do a docker login/pull from the harbor instance that should fail
7. Add the VCH as a cluster to Admiral in the default project
8. Enable content trust in the default project of Admiral
9. Do a docker pull from docker hub that should succeed
10. Do a docker login/pull from the harbor instance that should now succeed
11. Disable content trust in the default project of Admiral
12. Do a docker pull from docker hub that should succeed
13. Do a docker login/pull from the harbor instance that should now fail again
14. Create a new project in Admiral called 'definitely-not-default'
15. Enable content trust in the new project
16. Remove the VCH cluster from the default project and add it into the new project
17. Do a docker pull from docker hub that should succeed
18. Do a docker login/pull from the harbor instance that should now succeed again
19. Remove the VCH from all projects that it is still in within Admiral
20. Do a docker pull from docker hub that should fail since whitelist mode is enabled and docker hub is still not explicitly added
21. Do a docker login/pull from the harbor instance that should still succeed

# Expected Outcome:
Admiral and Engine should work together to obey the enable content trust feature properly.  When the VCH is within a project that has content trust enabled then users should be able to successfully pull from the Harbor instance

# Possible Problems:
None
