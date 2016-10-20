Test 6-09 - Verify vic-machine inspect
=======

#Purpose:
Verify vic-machine inspect functions

#References:
* vic-machine-linux inspect -h

#Environment:
This test requires that a vSphere server is running and available

#Test Steps
1. Create VCH1
2. Create VCH2
3. Run inspect for VCH1
4. Using inspect result to run docker command
5. Verify docker VM is created under correct VCH resource pool or Virtual App through govc
