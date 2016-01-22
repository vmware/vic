// +build linux
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"syscall"
)

var (
	ws syscall.WaitStatus = 0
)

// echoOff turns off the terminal echo.
func echoOff(fd []uintptr) (int, error) {
	pid, err := syscall.ForkExec(sttyArg0, sttyArgvEOff, &syscall.ProcAttr{Dir: "", Files: fd})
	if err != nil {
		return 0, fmt.Errorf("failed turning off console echo for password entry:\n\t%s", err)
	}
	return pid, nil
}

// echoOn turns back on the terminal echo.
func echoOn(fd []uintptr) {
	// Turn on the terminal echo.
	pid, e := syscall.ForkExec(sttyArg0, sttyArgvEOn, &syscall.ProcAttr{Dir: "", Files: fd})
	if e == nil {
		syscall.Wait4(pid, &ws, 0, nil)
	}
}

func testWriteConnection(conn net.Conn) {
	//var buffer []byte
	buffer := make([]byte, 0, 256)
	cmd_reader := bufio.NewReader(os.Stdin)

	std_fd := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}

	pid, err := echoOff(std_fd)

	if err != nil {
		return
	}

	defer echoOn(std_fd)

	for {

		text, _ := cmd_reader.ReadByte()

		if text == 3 {
			return
		}
		buffer = append(buffer, text)
		//buffer[0] = text
		conn.Write(buffer)
		buffer = buffer[:0]
		//fmt.Println(text)
		//testReadConnection(conn)
	}
	syscall.Wait4(pid, &ws, 0, nil)
}

func testReadConnection(conn net.Conn) {
	//var buffer []byte
	buffer := make([]byte, 256)
	writer := bufio.NewWriter(os.Stdout)
	for {
		read_size, _ := conn.Read(buffer)

		//fmt.Print(string(buffer[:read_size]))
		writer.Write(buffer[:read_size])
		writer.Flush()

	}
}
