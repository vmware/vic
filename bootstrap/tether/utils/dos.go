// Basic helper functions for interacting with DOS
package utils

import (
	"bufio"
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
