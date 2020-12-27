package main

import (
	"flag"
	"fmt"
	"github.com/eyedeekay/go-i2pcontrol"
	"log"
	"strings"
	"time"
)

var usage = `i2p-control
===========

Terminal interface to monitor and manage I2P router service. Basically, an
terminal i2pcontrol client.

        -host default:"127.0.0.1"
        -port default:"7657"
        -path default:"jsonrpc"
        -password default:"itoopie"
        -method default:"echo"
        -block default:false
        -verbose default:false

Installation with go get

        go get -u github.com/eyedeekay/i2p-control

The methods that have been implemented are

        echo              : i2pcontrol:Echo
        stat              : i2pcontrol:RouterInfo:i2p.router.status
        netstat           : i2pcontrol:RouterInfo:i2p.router.net.router.status
        tunstat           : i2pcontrol:RouterInfo:i2p.router.net.tunnels.participating
        restart           : i2pcontrol:Restart
        graceful-restart  : i2pcontrol:RestartGraceful
        shutdown          : i2pcontrol:Shutdown
        graceful-shutdown : i2pcontrol:ShutdownGraceful
        update            : i2pcontrol:Update
        find-update       : i2pcontrol:FindUpdate

So, for instance, to initiate a graceful shutdown and block until the router is
shut down, use the command:

        i2p-control -block -method=graceful-shutdown

`

var (
	host     = flag.String("host", "localhost", "Host of the i2pcontrol interface")
	port     = flag.String("port", "7657", "Port of the i2pcontrol interface")
	path     = flag.String("path", "jsonrpc", "Path to the i2pcontrol interface")
	password = flag.String("password", "itoopie", "Password for the i2pcontrol interface")
	command  = flag.String("method", "echo", "Method call to invoke")
	shelp    = flag.Bool("h", false, "Show the help message")
	sverbose = flag.Bool("v", false, "Verbosely update participating tunnel count while running.")
	lverbose = flag.Bool("verbose", false, "Verbosely update participating tunnel count while running.")
	lhelp    = flag.Bool("help", false, "Show the help message")
	block    = flag.Bool("block", false, "Block the terminal until the router is completely shut down")
)

func main() {
	flag.Parse()
	if *shelp || *lhelp {
		fmt.Printf(usage)
		return
	}
	shuttingdown := false
	i2pcontrol.Initialize(*host, *port, *path)
	_, err := i2pcontrol.Authenticate(*password)
	if err != nil {
		log.Fatal(err)
	}
	verbose := false
	if *sverbose {
		verbose = true
	}
	if *lverbose {
		verbose = true
	}
	if verbose {
		log.Println(*command)
	}
	switch *command {
	case "echo":
		message, err := i2pcontrol.Echo(strings.Join(flag.Args(), " "))
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
	case "restart":
		message, err := i2pcontrol.Restart()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
		shuttingdown = true
	case "graceful-restart":
		message, err := i2pcontrol.RestartGraceful()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
		shuttingdown = true
	case "shutdown":
		message, err := i2pcontrol.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
		shuttingdown = true
	case "graceful-shutdown":
		message, err := i2pcontrol.ShutdownGraceful()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
		shuttingdown = true
	case "update":
		message, err := i2pcontrol.FindUpdates()
		if err != nil {
			log.Fatal(err)
		}
		if message {
			log.Println("You need an update")
			message, err := i2pcontrol.Update()
			if err != nil {
				log.Fatal(err)
			}
			log.Println(message)
		}
		log.Println("You don't need an update")
	case "find-update":
		message, err := i2pcontrol.FindUpdates()
		if err != nil {
			log.Fatal(err)
		}
		if message {
			log.Println("You need an update")
			return
		}
		log.Println("You don't need an update")
	case "stat":
		message, err := i2pcontrol.Status()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
	case "netstat":
		message, err := i2pcontrol.NetStatus()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
	case "reseedstat":
		message, err := i2pcontrol.Reseeding()
		if err != nil {
			log.Fatal(err)
		}
		if message {
			log.Println("Router is reseeding")
		} else {
			log.Println("Router is not reseeding")
		}
	case "tunstat":
		message, err := i2pcontrol.ParticipatingTunnels()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
	}
	if shuttingdown {
		lastParticipatingTunnels, err := i2pcontrol.ParticipatingTunnels()
		if err != nil {
			log.Fatal(err)
		}
		for *block {
			participatingTunnels, err := i2pcontrol.ParticipatingTunnels()
			if err != nil {
				log.Fatal(err)
			}
			if participatingTunnels != lastParticipatingTunnels {
				log.Println("Waiting for expiration of:", participatingTunnels, "participating tunnels.")
				lastParticipatingTunnels = participatingTunnels
			}
			time.Sleep(time.Duration(time.Second * 1))
			if participatingTunnels < 1 {
				*block = false
				break
			}
		}
	}
}
