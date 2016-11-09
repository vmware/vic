Test 6-01 - Verify Help
=======

# Purpose:
Verify vic-machine delete help

# References:
* vic-machine-linux delete -h

# Environment:
Standalone test requires nothing but vic-machine to be built

# Test Cases

## Delete help basic
1. Issue the following command:
```
* vic-machine-linux delete -h
```

### Expected Outcome:
* Command should output Usage of vic-machine delete -h:
