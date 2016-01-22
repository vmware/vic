package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"

	"enatai-gerrit.eng.vmware.com/bonneville-container/tether"
)

func (ch *GlobalHandler) StartConnectionManager(conn *ssh.ServerConn) {
	log.Println("Registering fork trigger signal handler")
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in StartConnectionManager", r)
		}
	}()

	var incoming = make(chan os.Signal, 1)
	signal.Notify(incoming, syscall.SIGABRT)

	log.Println("SIGABRT handling initialized for fork support")
	for _ = range incoming {
		// validate that this is a fork trigger and not just a random signal from
		// container processes
		log.Println("Received SIGABRT - preparing to transition to fork parent")
		break
	}

	// tell client that we're disconnecting
	if ok, _, err := conn.SendRequest("fork", true, nil); !ok || err != nil || !ok {
		if err != nil {
			log.Printf("Unable to inform remote about fork (channel error): %s\n", err)
		} else {
			log.Println("Unable to register fork with remote - remote error")
		}
	} else {

		log.Println("Closing control connections")

		// regardless of errors we have to continue if externally driven
		conn.Close()

		// TODO: do we need to rebind session executions stdio to /dev/null or to files?
		// run the /.tether/vmfork.sh script
		log.Println("Running vmfork.sh")
		cmd := exec.Command("/.tether/vmfork.sh")
		// FORK HAPPENS DURING CALL, BEFORE RETURN FROM COMBINEDOUTPUT
		out, err := cmd.CombinedOutput()
		log.Printf("vmfork:%s\n%s\n", err, string(out))

		return
	}

	log.Println("Closing control connections")

	// regardless of errors we have to continue if externally driven
	conn.Close()

	// the StartTether loop will now exit and we'll fall back into waiting for SIGHUP in main
}

func (ch *GlobalHandler) ContainerId() string {
	return ch.id
}

func (ch *GlobalHandler) StaticIPAddress(cidr, gateway string) error {
	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	cmd := exec.Command("/Windows/System32/netsh.exe", "interface", "ipv4", "set", "address", "name", "Local Area Connection", "source=static", "address="+ip.String(), "mask="+net.Mask.String(), "gateway="+gateway)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error setting IP address: %s\n%s", err, out)
	}
	return err

	/*
		link, err := netlink.LinkByName("eth0")
		if err != nil {
			detail := fmt.Sprintf("failed to get link for eth0: %s", err)
			log.Print(detail)
			return errors.New(detail)
		}

		// set the ip address for our interface
		addr, err := netlink.ParseAddr(cidr)
		if err != nil {
			detail := fmt.Sprintf("failed to parse address for eth0: %s", err)
			log.Print(detail)
			return errors.New(detail)
		}

		existingAddr, err := netlink.AddrList(link, syscall.AF_UNSPEC)
		if err != nil {
			detail := fmt.Sprintf("failed to list existing address for eth0: %s", err)
			log.Print(detail)
			return errors.New(detail)
		}

		for _, oldAddr := range existingAddr {
			err = netlink.AddrDel(link, &oldAddr)
			if err != nil {
				detail := fmt.Sprintf("failed to del existing address for eth0: %s", err)
				log.Print(detail)
			}
		}

		if err := netlink.AddrAdd(link, addr); err != nil {
			detail := fmt.Sprintf("failed to add address to eth0: %s", err)
			log.Print(detail)
			return errors.New(detail)
		}

		// bring the interface up
		if err = netlink.LinkSetUp(link); err != nil {
			detail := fmt.Sprintf("failed to bring up eth0: %s", err)
			log.Print(detail)
			return errors.New(detail)
		}

		// add a default route
		//route := netlink.Route{LinkIndex: link.Attrs().Index, Dst: ipnet, Src: ip, Gw: net.ParseIP(gateway)}
		_, defaultNet, _ := net.ParseCIDR("0.0.0.0/0")
		route := netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet, Gw: net.ParseIP(gateway)}
		err = netlink.RouteAdd(&route)
		if err != nil {
			detail := fmt.Sprintf("failed to add route for gateway: %s", err)
			log.Print(detail)
			return errors.New(detail)
		}

		// TODO: even if we want to keep an entry in the local /etc/hosts to speed resolution
		// we should do it in a maner that doesn't result in it being present in a diff with the parent
		f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Print(err.Error())
			return err
		}
		defer f.Close()

		hostname, err := os.Hostname()
		if err != nil {
			log.Print(err.Error())
			return err
		}

		if _, err = f.WriteString(fmt.Sprintf("%s %s\n", ip.String(), hostname)); err != nil {
			log.Fatal(err.Error())
			return err
		}
	*/

	return nil
}

func (ch *GlobalHandler) DynamicIPAddress() (string, error) {
	return "", nil
}

func (c *GlobalHandler) MountLabel(label, target string) error {
	/*
		if err := os.MkdirAll(target, 0600); err != nil {
			return fmt.Errorf("Unable to create mount point %s: %s", target, err)
		}

		//volumes := "/.tether/volumes"
		volumes := "/dev/disk/by-label"

		source := volumes + "/" + label

		// wait for mount source to appear or timeout
		for start := time.Now(); time.Since(start) < (5 * time.Second); {
			_, err := os.Stat(source)
			if err == nil || !os.IsNotExist(err) {
				break
			}
		}

		if err := syscall.Mount(source, target, "ext4", syscall.MS_NOATIME, ""); err != nil {
			detail := fmt.Sprintf("Unable to mount %s: %s", source, err)
			log.Print(detail)
			// for debug purposes, dump the directory listing of volumes and /dev/disk/by-label
			for _, dir := range []string{volumes, "/dev/disk/by-label"} {
				files, err := ioutil.ReadDir(dir)
				if err != nil {
					log.Printf("unable to read listing for %s: %s\n", dir, err)
					continue
				}

				log.Printf("%s/\n", dir)
				for _, fi := range files {
					log.Printf("\t%s\n", fi.Name())
				}
			}

			return errors.New(detail)
		}
	*/
	return nil
}

func (c *GlobalHandler) Sync() {
	/*
		syscall.Sync()
	*/
}

func (c *GlobalHandler) NewSessionContext() tether.SessionContext {
	return &SessionHandler{
		GlobalHandler: c,
		env:           make(map[string]string),
	}
}
