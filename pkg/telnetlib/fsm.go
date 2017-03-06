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

import log "github.com/Sirupsen/logrus"

type state int

const (
	dataState state = iota
	optionNegotiationState
	cmdState
	subnegState
	subnegEndState
	errorState
)

type telnetFSM struct {
	curState state
	tc       *TelnetConn
	doneCh   chan chan struct{}
}

func newTelnetFSM() *telnetFSM {
	f := &telnetFSM{
		doneCh:   make(chan chan struct{}),
		curState: dataState,
	}
	return f
}

func (fsm *telnetFSM) start() {
	defer func() {
		log.Infof("FSM closed")
	}()
	for {
		select {
		case ch := <-fsm.tc.fsmInputCh:
			//log.Infof("FSM state is %d", fsm.curState)
			ns := fsm.nextState(ch)
			fsm.curState = ns
		case ch := <-fsm.doneCh:
			ch <- struct{}{}
			return
		}
	}
}

// this function returns what the next state is and performs the appropriate action
func (fsm *telnetFSM) nextState(ch byte) state {
	var nextState state
	b := []byte{ch}
	switch fsm.curState {
	case dataState:
		if ch != IAC {
			fsm.tc.dataRW.Write(b)
			fsm.tc.dataWrittenCh <- true
			nextState = dataState
		} else {
			nextState = cmdState
		}

	case cmdState:
		if ch == IAC { // this is an escaping of IAC to send it as data
			fsm.tc.dataRW.Write(b)
			fsm.tc.dataWrittenCh <- true
			nextState = dataState
		} else if ch == DO || ch == DONT || ch == WILL || ch == WONT {
			fsm.tc.cmdBuffer.WriteByte(ch)
			nextState = optionNegotiationState
		} else if ch == SB {
			fsm.tc.cmdBuffer.WriteByte(ch)
			nextState = subnegState
		} else { // anything else
			fsm.tc.cmdBuffer.WriteByte(ch)
			fsm.tc.cmdHandlerWrapper(fsm.tc.handlerWriter, &fsm.tc.cmdBuffer)
			fsm.tc.cmdBuffer.Reset()
			nextState = dataState
		}
	case optionNegotiationState:
		fsm.tc.cmdBuffer.WriteByte(ch)
		opt := ch
		cmd := fsm.tc.cmdBuffer.Bytes()[0]
		fsm.tc.optionCallback(cmd, opt)
		fsm.tc.cmdBuffer.Reset()
		nextState = dataState
	case subnegState:
		if ch == IAC {
			nextState = subnegEndState
		} else {
			nextState = subnegState
			fsm.tc.cmdBuffer.WriteByte(ch)
		}
	case subnegEndState:
		if ch == SE {
			fsm.tc.cmdBuffer.WriteByte(ch)
			fsm.tc.cmdHandlerWrapper(fsm.tc.handlerWriter, &fsm.tc.cmdBuffer)
			fsm.tc.cmdBuffer.Reset()
			nextState = dataState
		} else if ch == IAC { // escaping IAC
			nextState = subnegState
			fsm.tc.cmdBuffer.WriteByte(ch)
		} else {
			nextState = errorState
		}
	case errorState:
		nextState = dataState
		log.Infof("Finite state machine is in an error state. This should not happen for correct telnel protocol syntax")
	}
	return nextState
}
