Test 2-2 - No Arguments
=======

#Purpose:
To verify that when imagec is run without arguments that it fails.

#References:
* imagec --help

#Environment:
Standalone test requires nothing but imagec to be built

#Test Steps:
1. Issue the following command:
* imagec -standalone

#Expected Outcome:
* Command should return failure

#Possible Problems:
Make sure that you run imagec on the same hard drive partition as /tmp, otherwise you will receive a cross device link error.