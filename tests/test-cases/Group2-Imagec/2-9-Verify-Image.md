Test 2-9 - Verify Image
=======

#Purpose:
To verify that when imagec is run with -image flag, then it will download a specific image correctly

#References:
* imagec --help

#Environment:
Standalone test requires nothing but imagec to be built

#Test Steps:
1. Issue the following command:
* imagec -standalone -reference tatsushid/tinycore

#Expected Outcome:
* Command should return success
* All the checksums for each image layer should match the manifest file

#Possible Problems:
Make sure that you run imagec on the same hard drive partition as /tmp, otherwise you will receive a cross device link error.