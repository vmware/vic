# Container VM output logging in VIC

When a container VM process produces output, that output is sent down through a serial connection into the vSphere backend where it is written to a logfile. When a user wants to see this output, they use the `docker logs` command, which reads from this logfile and displays the log entries on the command line in order from oldest to newest.

However, in order for users to be able to use the `--since` option with `docker logs`, we must also couple timestamps to these log entries so that the user can selectively filter log messages to be displayed based on when those entries occurred.

### Requirements

The container VM logging mechanism must:
  1. Wrap each log entry in a JSON struct containing both the log message(s) in that entry, and the timestamp in which that entry occurred. 
 2. Allow these entries to be read and unwrapped at a later time, starting with the first entry occurring at or beyond the timestamp suppled by the user with the `--since` option to `docker logs`, or starting with the first entry in the logfile if `--since` was not used. 

### Implementation

A package for ioutils in VIC is added as `github.com/vmware/vic/lib/ioutil`. In this package are two files, `log_writer.go` and `log_reader.go`. These files contain an implementation of Go's `io.Writer` and `io.Reader` interface, respectively.

The responsibilities of the `LogReader` are:
 1. To instantiate a `LogEntry` struct, and add the log message and timestamp to the `LogEntry`'s `Log` and `Time` fields. 
 2. To serialize and write this struct data to the serial port associated with the containerVM logfile on the backend, and follow that write with a single newline character so that the logfile can be more easily consumed by humans for debugging purposes. 

The responsibilities of the `LogWriter` are:
 1. To scan in lines from the containerVM logfile one at a time, and unmarshal them into `LogEntry` structs. 
 2. To copy the `LogEntry`'s `Log` field into the underlying `Read` stream's `[]byte` slice. 
 3. To preserve unwritten bytes in a call to `Read` in memory so that they may be written during the next call, in the case where the supplied `[]byte` slice was smaller than the log message we are trying to write. 
