Test 2-5 - Verify Debug And Stdout
=======

#Purpose:
To verify that when imagec is run with -debug and -stdout flags, then it returns additional debug output to the screen.

#References:
* imagec --help

#Environment:
Standalone test requires nothing but imagec to be built

#Test Steps:
1. Issue the following command:
* imagec -standalone -reference photon -stdout -debug

#Expected Outcome:
* Output from the command should contain at least one line of level=debug
* Command should return success
* All the checksums for each image layer should match the manifest file

#Possible Problems:
Make sure that you run imagec on the same hard drive partition as /tmp, otherwise you will receive a cross device link error.