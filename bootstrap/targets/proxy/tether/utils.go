// Basic helper functions for interacting with DOS
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// Output from a cmd call to MS-DOS consists of an echo of the command and the output, separated by CRLF
// We want just the second line without the CRLF
func StripCommandOutput(output string) string {
	reader := strings.NewReader(output)
	reader2 := bufio.NewReader(reader)
	_, _ = reader2.ReadString('\n') // throw away the echo line
	line2, err := reader2.ReadString('\n')
	if err == nil && len(line2) > 2 {
		return line2[:len(line2)-2] // Assume line ends with CRLF
	}
	return ""
}

// Reads from the dos connection until it sees a command prompt
// We're presuming that means end of output
func (ch *GlobalHandler) discardUntilPrompt() error {
	tmp := make([]byte, 16)
	for {
		n, err := ch.dosconn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
				return err
			}
			break
		}
		/* consume prompt text */
		if tmp[n-1] == '>' {
			break
		}

	}

	return nil
}

// Reads from the dos connection until it sees a command prompt
// We're presuming that means end of output
func (ch *GlobalHandler) readUntilPrompt() (string, error) {
	buf := make([]byte, 0, 512)
	tmp := make([]byte, 16)

	for {
		n, err := ch.dosconn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
				return string(buf), err
			}
			break
		}

		buf = append(buf, tmp[:n]...)

		if tmp[n-1] == '>' {
			break
		}

	}

	// TODO: back track to last \n and strip entire prompt
	return string(buf), nil
}

// Run the command and return the output as a string
func (ch *GlobalHandler) cmdCombinedOutput(cmd string) (string, error) {
	_, _ = ch.dosconn.Write([]byte(cmd + "\r\n"))

	return ch.readUntilPrompt()
}

// Run the command and return
func (ch *GlobalHandler) cmdStart(cmd string) error {
	_, _ = ch.dosconn.Write([]byte(cmd + "\r\n"))

	// check for "command not found" without taking
	// it off the stream

	return nil
}

// Reads from the dos connection until it sees a command prompt
// We're presuming that means end of output
// Copies the data to the handler.chennel
func (ch *SessionHandler) copyUntilPrompt() error {
	// wait for the exec reply to be sent
	ch.waitGroup.Add(1)
	var fatalerr error = nil
	log.Println("copying the data from ssh to dosconn and vice versa")
	go func() {
		buf := make([]byte, 16)
		sshChannel := *ch.channel
		exitTrigger := magicPrompt + ">" // We're monitoring the output looking for this string
		triggerPos := 0
	ReadLoop:
		for {
			nr, er := ch.dosconn.Read(buf)
			// log.Printf(">\"%s\"", string(buf[:]))
			if nr > 0 {
				for i := 0; i < nr; i++ {
					// log.Printf("comparing %c with %c: %d, %d", exitTrigger[triggerPos], buf[i], triggerPos, i)
					if exitTrigger[triggerPos] == buf[i] {
						triggerPos++
						if triggerPos == len(exitTrigger) {
							sshChannel.Write([]byte{13, 10}) // crlf on the end
							break ReadLoop
						}
					} else {
						triggerPos = 0
					}
				}
				nw, ew := sshChannel.Write(buf[0:nr])
				if ew != nil {
					fatalerr = fmt.Errorf("Error writing to ssh channel: %s", ew)
					break
				}
				if nw != nr {
					fatalerr = fmt.Errorf("Short write! %n != %n", nr, nw)
					break
				}
			}
			if er == io.EOF {
				break
			} else if er != nil {
				if e, ok := er.(*net.OpError); ok { // We can't distinguish between different types of OpError, but most are fatal
					log.Printf("Lost read connection: %s. Exiting", e)
					fatalerr = er
					break
				}
			}
		}
		ch.waitGroup.Done()
	}()

	// we shouldn't wait on stdin to close or we'll be here forever
	go func() { log.Println("Copying from daemon to dosconn"); io.Copy(ch.dosconn, *ch.channel) }()

	log.Println("Waiting for I/O streams to close")
	ch.waitGroup.Wait()
	if fatalerr != nil {
		log.Println("Returning read/write error: ", fatalerr)
	}
	return fatalerr
}

// Run the command and copy the output until it's done
func (ch *SessionHandler) cmdRun(cmd string) error {
	_, _ = ch.dosconn.Write([]byte(cmd + "\r\n"))

	return ch.copyUntilPrompt()
}
