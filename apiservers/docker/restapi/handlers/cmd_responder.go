package handlers

import (
    "io"
    "log"
   	"os/exec"
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"
)

type CmdResponder struct {
    writer      io.Writer
    flusher     http.Flusher
    cmdPath     string
    cmdArgs     []string
}

func NewCmdResponder(path string, args []string) *CmdResponder {
	responder := &CmdResponder{cmdPath: path, cmdArgs: args}

	return responder
}

func (self *CmdResponder) Write(data []byte) (int, error) {
	n, err := self.writer.Write(data)

	self.flusher.Flush()

	return n, err
}

// WriteResponse to the client
func (self *CmdResponder) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {
    var exist bool
    
	rw.Header().Set("Content-Type", "application/json")
	
	self.flusher, exist = rw.(http.Flusher)

	if exist {
        self.writer = rw
        
        cmd := exec.Command(self.cmdPath, self.cmdArgs...)
        cmd.Stdout = self
        cmd.Stderr = self

        // Execute
        err := cmd.Start()

        if err != nil {
            log.Printf("Error starting %s - %s\n", self.cmdPath, err)
            rw.WriteHeader(http.StatusInternalServerError)
        }

        // Wait for the fetcher to finish.  Should effectively close the stdio pipe above.
        err = cmd.Wait()
        
        if err != nil {
           log.Println("imageC exit code:", err)
        }
        
        rw.WriteHeader(http.StatusOK)
	} else {
		log.Println("ChunkResponder failed to get the HTTP flusher")
	}
}
