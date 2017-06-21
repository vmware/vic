// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package toolbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestProcessFunction(t *testing.T) {
	m := NewProcessManager()
	var pids []int64

	for i := 0; i <= 2; i++ {
		r := &VixMsgStartProgramRequest{
			ProgramPath: "test",
			Arguments:   strconv.Itoa(i),
		}

		pid, _ := m.Start(r, NewProcessFunc(func(_ context.Context, arg string) error {
			rc, _ := strconv.Atoi(arg)
			if rc == 0 {
				return nil

			}
			return &ProcessError{Err: errors.New("fail"), ExitCode: int32(rc)}
		}))

		if pid == 0 {
			t.Fatalf("no pid")
		}

		pids = append(pids, pid)
	}

	m.wg.Wait()

	_ = m.ListProcesses(pids)

	for i, pid := range pids {
		p := m.entries[pid]
		if p.ExitCode != int32(i) {
			t.Errorf("%d: %d != %d", pid, p.ExitCode, i)
		}
	}
}

func TestProcessCommand(t *testing.T) {
	m := NewProcessManager()
	var pids []int64

	for i := 0; i <= 2; i++ {
		r := &VixMsgStartProgramRequest{
			ProgramPath: "/bin/bash",
			Arguments:   fmt.Sprintf(`-c "exit %d"`, i),
		}

		pid, _ := m.Start(r, NewProcess())
		pids = append(pids, pid)
	}

	m.wg.Wait()

	_ = m.ListProcesses(nil)

	for i, pid := range pids {
		p := m.entries[pid]
		if p.ExitCode != int32(i) {
			t.Errorf("%d: %d != %d", pid, p.ExitCode, i)
		}
	}

	r := &VixMsgStartProgramRequest{
		ProgramPath: shell,
	}

	shell = "/enoent/enoent"
	_, err := m.Start(r, NewProcess())
	if err == nil {
		t.Error("expected error")
	}
	shell = r.ProgramPath

	r.ProgramPath = "/enoent/enoent"
	_, err = m.Start(r, NewProcess())
	if err == nil {
		t.Error("expected error")
	}
}

func TestProcessKill(t *testing.T) {
	m := NewProcessManager()
	var pids []int64

	procs := []struct {
		r *VixMsgStartProgramRequest
		p *Process
	}{
		{
			&VixMsgStartProgramRequest{
				ProgramPath: "test",
				Arguments:   "none",
			},
			NewProcessFunc(func(ctx context.Context, _ string) error {
				select {
				case <-ctx.Done():
					return &ProcessError{Err: ctx.Err(), ExitCode: 42}
				case <-time.After(time.Minute):
				}

				return nil
			}),
		},
		{
			&VixMsgStartProgramRequest{
				ProgramPath: "/bin/bash",
				Arguments:   fmt.Sprintf(`-c "while true; do sleep 1; done"`),
			},
			NewProcess(),
		},
	}

	for _, test := range procs {
		pid, err := m.Start(test.r, test.p)
		if err != nil {
			t.Fatal(err)
		}

		pids = append(pids, pid)
	}

	for {
		b := m.ListProcesses(pids)
		if bytes.Count(b, []byte("<proc>")) == len(pids) {
			break
		}

		<-time.After(time.Millisecond * 100)
	}

	for _, pid := range pids {
		if !m.Kill(pid) {
			t.Errorf("kill %d", pid)
		}
	}

	m.wg.Wait()

	for _, pid := range pids {
		p := m.entries[pid]

		if p.ExitCode == 0 {
			t.Errorf("%s: exit=%d", p.Name, p.ExitCode)
		}
	}

	if m.Kill(-1) {
		t.Error("kill -1")
	}
}

func TestProcessRemove(t *testing.T) {
	m := NewProcessManager()

	m.expire = time.Millisecond

	r := &VixMsgStartProgramRequest{
		ProgramPath: "test",
	}

	pid, _ := m.Start(r, NewProcessFunc(func(_ context.Context, arg string) error {
		return nil
	}))

	m.wg.Wait()

	<-time.After(m.expire * 20)
	// pid should be removed by now
	b := m.ListProcesses([]int64{pid})
	if len(b) != 0 {
		t.Error("expected 0 processes")
	}
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{`echo "foo bar" > /dev/null`, "echo %22foo bar%22 %3E /dev/null"},
	}

	for i, test := range tests {
		e := xmlEscape.Replace(test.in)
		if e != test.out {
			t.Errorf("%d: %s != %s", i, e, test.out)
		}
	}
}

func TestProcessError(t *testing.T) {
	fault := errors.New("fail")
	var err error

	err = &ProcessError{Err: fault}

	if err.Error() != fault.Error() {
		t.Fatal()
	}
}
