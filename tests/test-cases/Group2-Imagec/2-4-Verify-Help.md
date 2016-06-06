Test 2-4 - Verify Help
=======

#Purpose:
To verify that when imagec is run with --help, then it provides the usage output to the user

#References:
* imagec --help

#Environment:
Standalone test requires nothing but imagec to be built

#Test Steps:
1. Issue the following command:
* imagec --help

#Expected Outcome:
* Command should return error code
* Command should output Usage of imagec:

#Possible Problems:
Make sure that you run imagec on the same hard drive partition as /tmp, otherwise you will receive a cross device link error.