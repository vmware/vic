package main

import (
	"log"
	"time"

	"golang.org/x/sys/windows/svc"
)

type tetherservice struct{}

func (m *tetherservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			default:
				log.Printf("unexpected control request #%d\n", c)
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string, isDebug bool) {
	var err error

	log.Printf("starting %s service\n", name)
	err = svc.Run(name, &tetherservice{})
	if err != nil {
		log.Printf("%s service failed: %v\n", name, err)
		return
	}
	log.Printf("%s service stopped\n", name)
}
