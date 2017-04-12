Test 1-38 - Docker Stats
=======

#Purpose:
To verify that `docker stats` is supported and works as expected.

#Environment:
This test requires that a vSphere server is running and available


#Test Steps:
1. Run a busybox container and create a busybox container
2. Run Stats with no-stream which will return stats for any running container
3. Run Stats with no-stream  all which will return stats for all containers
4. Verify the API memory output against govc
5. Verify the API CPU output against govc


#Expected Outcome:
1. Fails if two containers are not created
2. Return stats for a single container and validate memory -- will fail if there's too
   much variation in the memroy
3. Return stats for all containers -- will fail if output is missing either container
4. Compare API results vs. govc result for memory accuracy -- will fail if large variation
5. Compare API results vs. govc result for CPU accuracy -- will fail if API value not present in past
   six govc readings



#Possible Problems:
Stats are created by the ESXi host every 20s -- if there are long pauses between calls
in a single test the results could be incorrect and a failure could occur.
