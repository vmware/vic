Test 6-13 - Verify proper volume store option behavior
=======

# Purpose:
Verify vic-machine create with volume store options behaves as we expect it to. 

# References:
* vic-machine-linux create -dvs <path> -vs <path:label>

# Environment:
This test requires that a vSphere server is running and available

# Test Cases


## Create with default volume store

Tests just the -dvs command for passing functionality

### Expected Outcome:

a successful vic-machine create should occur and there should be a volumestore tagged as default. 

## Create default store with volume store flag

Tests assigning the default store using the -vs for passing functionality. This is a test designed to protect backwards compatibility.

### Expected Outcome: 
a successful vic-machine create should occur and there should be a volumestore tagged as default. 

## Create default store using both flags with the same path

Tests assigning the default store using the -vs and -vsd with the same target path for passing functionality. This is a test designed to protect backwards compatibility.

### Expected Outcome: 


a successful vic-machine create should occur and there should be a volumestore tagged as default. 

## Create default store using both flags without the same path

Tests assigning the default volume store using the -vs and -vsd flags while assigning two different flags.

### Expected Outcome: 

vic-machine should return an error message stating that multiple path's were used when trying to set the volumes "default". vic-machine should also suggest the use of --default-volume-store(-vsd) as the preferred path. 

## Create with normal volume store

Tests just the -vs command for passing functionality

### Expected Outcome:

a successful vic-machine create should occur and there should be volumestores tagged as "test" and "cheap" 

## Create with both volume store flags

Tests the -vs and -dvs commands for passing functionality

### Expected Outcome:

a successful vic-machine create should occur and there should be volumestores tagged as "test" and "cheap"  as well as a default volume store

## Create with overlapping paths to single label

Tests the situation where a label has be targeted with two separate paths using the -vs flag. 

### Expected Outcome:

vic-machine create should fail and report that there was a label with two distinct paths

