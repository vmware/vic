Test 2-10 - Verify Logfile
=======

#Purpose:
To verify that when imagec is run with -logfile flag, it should change the path of the installer log file (imagec.log)

#References:
* imagec --help

#Environment:
Standalone test requires nothing but imagec to be built

#Test Steps:
1. Issue the following command:
* imagec -standalone -reference photon -logfile foo.log

#Expected Outcome:
* Command should return success
* foo.log should exist and not be empty

#Possible Problems:
Make sure that you run imagec on the same hard drive partition as /tmp, otherwise you will receive a cross device link error.