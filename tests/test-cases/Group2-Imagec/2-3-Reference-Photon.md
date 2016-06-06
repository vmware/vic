Test 2-3 - Reference Photon
=======

#Purpose:
To verify that when imagec is run with -reference photon then it should download library/photon image

#References:
* imagec --help

#Environment:
Standalone test requires nothing but imagec to be built

#Test Steps:
1. Issue the following command:
* imagec -standalone -reference photon

#Expected Outcome:
* Command should return success
* An images directory should be created
* imagec.log file should be created and it should not be empty
* All the checksums for each image layer should match the manifest file

#Possible Problems:
Make sure that you run imagec on the same hard drive partition as /tmp, otherwise you will receive a cross device link error.