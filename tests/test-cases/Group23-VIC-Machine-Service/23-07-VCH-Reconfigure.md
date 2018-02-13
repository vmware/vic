Test 23-07 - VCH Reconfigure
============================

#### Purpose:
To verify vic-machine-server correctly reconfigures VCHs

#### References:
1. [VIC Machine Service API Design Doc](../../../doc/design/vic-machine/service.md)

#### Environment:
This suite requires a vSphere system where VCHs can be deployed. Ideally, this suite would be executed against multiple such environments, including directly against ESX, against a VC, and against a VC with multiple datacenters configured.

**Note:** Multiple versions of tests must be implemented:
 * Using both PUT and PATCH.
 * With and without datacenter-scoped URLs.
 * Using both username/password- and session-based authentication.

It should not be necessary to implement all eight combinations for every test case, but full coverage should be provided for at least basic test cases.


Basic Validation
----------------

These tests are intended to verify the feature at a basic level.

###  1. Create a simple VCH using the API and then reconfigure its name

###  2. Create a simple VCH using the API and then reconfigure all of its mutable properties

###  3. Create a complex VCH using the API and then reconfigure all of its mutable properties


Negative Cases
--------------

These tests are intended to verify the behavior in various failure cases. We must gracefully handle invalid input and unexpected user behavior.

###  1. Attempt to reconfigure a non-existent VCH

###  2. Attempt to reconfigure a VCH which has been deleted

###  3. Attempt to reconfigure a VCH with a malformed body

###  4. Attempt to reconfigure a VCH in a different datacenter

###  5. Attempt to reconfigure each immutable property of a VCH *(one per request)*


Interoperability
----------------

These tests are intended to verify that the API and CLI can coexist without issue.

###  1. Use the API to reconfigure a VCH created via the CLI

###  2. Use the CLI to reconfigure a VCH created via the API


Workflow-based
--------------

These tests are designed to mimic realistic customer scenarios. These tests will usually duplicate coverage provided by a test above, but provide additional validation around specific important workflows.

###  1. Create a representative VCH, reconfigure its certificates, and connect using the new certificates to confirm the change took effect

###  2. Create a representative VCH, reconfigure its vSphere credentials, and deploy a new container to confirm the change took effect

###  3. Create a representative VCH, reconfigure it to add a volume store, and deploy a new container using the volume store to confirm the change took effect

###  4. Create a representative VCH, reconfigure it to change the debug level, and verify that the change took effect



