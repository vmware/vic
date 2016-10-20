Test 6-10 - Verify vic-machine ls
=======

#Purpose:
Verify vic-machine ls functions

#References:
* vic-machine-linux ls -h

#Environment:
This test requires that a vSphere server is running and available

#Test Steps
1. Create VCH
3. Run ls to query out VCH
4. Run inspect to verify VCH id is correct
5. Run inspect to verify VCH compute path and name are correct
