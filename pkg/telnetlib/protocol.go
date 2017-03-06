// Copyright 2016 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telnetlib

const (
	/* Telnet Protocol Characters */
	IAC  byte = 255 // "Interpret As Command"
	DONT byte = 254
	DO   byte = 253
	WONT byte = 252
	WILL byte = 251
	NULL byte = 0

	SE  byte = 240 // Subnegotiation End
	NOP byte = 241 // No Operation
	DM  byte = 242 // Data Mark
	BRK byte = 243 // Break
	IP  byte = 244 // Interrupt process
	AO  byte = 245 // Abort output
	AYT byte = 246 // Are You There
	EC  byte = 247 // Erase Character
	EL  byte = 248 // Erase Line
	GA  byte = 249 // Go Ahead
	SB  byte = 250 // Subnegotiation Begin

	/* Telnet Options */
	BINARY              byte = 0   // 8-bit data path
	ECHO                byte = 1   // echo
	RCP                 byte = 2   // prepare to reconnect
	SGA                 byte = 3   // suppress go ahead
	NAMS                byte = 4   // approximate message size
	STATUS              byte = 5   // give status
	TM                  byte = 6   // timing mark
	RCTE                byte = 7   // remote controlled transmission and echo
	NAOL                byte = 8   // negotiate about output line width
	NAOP                byte = 9   // negotiate about output page size
	NAOCRD              byte = 10  // negotiate about CR disposition
	NAOHTS              byte = 11  // negotiate about horizontal tabstops
	NAOHTD              byte = 12  // negotiate about horizontal tab disposition
	NAOFFD              byte = 13  // negotiate about formfeed disposition
	NAOVTS              byte = 14  // negotiate about vertical tab stops
	NAOVTD              byte = 15  // negotiate about vertical tab disposition
	NAOLFD              byte = 16  // negotiate about output LF disposition
	XASCII              byte = 17  // extended ascii character set
	LOGOUT              byte = 18  // force logout
	BM                  byte = 19  // byte macro
	DET                 byte = 20  // data entry terminal
	SUPDUP              byte = 21  // supdup protocol
	SUPDUPOUTPUT        byte = 22  // supdup output
	SNDLOC              byte = 23  // send location
	TTYPE               byte = 24  // terminal type
	EOR                 byte = 25  // end or record
	TUID                byte = 26  // TACACS user identification
	OUTMRK              byte = 27  // output marking
	TTYLOC              byte = 28  // terminal location number
	VT3270REGIME        byte = 29  // 3270 regime
	X3PAD               byte = 30  // X.3 PAD
	NAWS                byte = 31  // window size
	TSPEED              byte = 32  // terminal speed
	LFLOW               byte = 33  // remote flow control
	LINEMODE            byte = 34  // Linemode option
	XDISPLOC            byte = 35  // X Display Location
	OLD_ENVIRON         byte = 36  // Old - Environment variables
	AUTHENTICATION      byte = 37  // Authenticate
	ENCRYPT             byte = 38  // Encryption option
	NEW_ENVIRON         byte = 39  // New - Environment variables
	TN3270E             byte = 40  // TN3270E
	XAUTH               byte = 41  // XAUTH
	CHARSET             byte = 42  // CHARSET
	RSP                 byte = 43  // Telnet Remote Serial Port
	COM_PORT_OPTION     byte = 44  // Com Port Control Option
	SUPPRESS_LOCAL_ECHO byte = 45  // Telnet Suppress Local Echo
	TLS                 byte = 46  // Telnet Start TLS
	KERMIT              byte = 47  // KERMIT
	SEND_URL            byte = 48  // SEND-URL
	FORWARD_X           byte = 49  // FORWARD_X
	PRAGMA_LOGON        byte = 138 // TELOPT PRAGMA LOGON
	SSPI_LOGON          byte = 139 // TELOPT SSPI LOGON
	PRAGMA_HEARTBEAT    byte = 140 // TELOPT PRAGMA HEARTBEAT
	EXOPL               byte = 255 // Extended-Options-List
	NOOPT               byte = 0
)
