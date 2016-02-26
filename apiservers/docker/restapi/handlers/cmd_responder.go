package handlers

import (
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/go-swagger/go-swagger/httpkit"
)

type CmdResponder struct {
	writer  io.Writer
	flusher http.Flusher
	cmdPath string
	cmdArgs []string
}

func NewCmdResponder(path string, args []string) *CmdResponder {
	responder := &CmdResponder{cmdPath: path, cmdArgs: args}

	return responder
}

func (cr *CmdResponder) Write(data []byte) (int, error) {
	n, err := cr.writer.Write(data)

	cr.flusher.Flush()

	return n, err
}

// WriteResponse to the client
func (cr *CmdResponder) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {
	var exist bool

	rw.Header().Set("Content-Type", "application/json")

	cr.flusher, exist = rw.(http.Flusher)

	if exist {
		cr.writer = rw

		cmd := exec.Command(cr.cmdPath, cr.cmdArgs...)
		cmd.Stdout = cr
		cmd.Stderr = cr

		// Execute
		err := cmd.Start()

		if err != nil {
			log.Printf("Error starting %s - %s\n", cr.cmdPath, err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Wait for the fetcher to finish.  Should effectively close the stdio pipe above.
		err = cmd.Wait()

		if err != nil {
			log.Println("imagec exit code:", err)
		}

		rw.WriteHeader(http.StatusOK)
		return
	} else {
		log.Println("CmdResponder failed to get the HTTP flusher")
	}

	rw.WriteHeader(http.StatusInternalServerError)
}
