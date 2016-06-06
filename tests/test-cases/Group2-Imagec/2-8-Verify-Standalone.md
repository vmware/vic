Test 2-8 - Verify Standalone
=======

#Purpose:
To verify that when imagec is run with the -standalone option, then imagec can be run without portlayer API running

#References:
* imagec --help

#Environment:
Standalone test requires nothing but imagec to be built

#Test Steps:
1. Issue the following command:
* imagec -standalone -reference photon

#Expected Outcome:
* Command should return success
* All the checksums for each image layer should match the manifest file

#Possible Problems:
Make sure that you run imagec on the same hard drive partition as /tmp, otherwise you will receive a cross device link error.