package main

import (
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rancher/go-rancher-metadata/metadata"
	fastping "github.com/tatsushid/go-fastping"
)

var wg sync.WaitGroup

func pingContainers(retContainers []metadata.Container) {

	// Added a reset step to ensure the metrics are reset before values are set //
	// This is needed because of a case where rancher has rescheduled a container //
	// with a new index value, as a result the old index is always going to be marked down. //
	// this ends up generating a lot of alerts for containers which no longer exist //
	currentState.Reset()

	// Process the metadata to update metrics //
	
	for _, retContainer := range retContainers {
		//   Loop and check each container //
		// Ignore stopped containers, ones without ip's (generally sidekicks)
		if len(retContainer.PrimaryIp) != 0 && retContainer.State == "running" {
			wg.Add(1)
			go ping(retContainer.PrimaryIp, retContainer.Name)
		}

	}
	wg.Wait()

}

func ping(ip string, name string) {

	defer wg.Done()
	counter := 0 // Counter to add retry interval //
	transientState := 0

	p := fastping.NewPinger()

	// Check env variable to see if we want to use UDP //
	if MODE == "udp" {
		p.Network("udp") // Use udp to allow container to run as non-privledged one //
	}

	// Setup MaxRTT for the pinger //
	p.MaxRTT = 10 * time.Second

	ra, err := net.ResolveIPAddr("ip4", ip)

	checkError(err)

	p.AddIPAddr(ra)

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		// Define what to do when ip is reached //
		transientState = 1
	}

	p.OnIdle = func() {
		// Define what to do after maxRTT has exceeded //
		// Increment transientState for each failure //
	}

	// Run the container check in a loop //
	for counter < 1 {
		err = p.Run()
		counter++
	}

	currentState.WithLabelValues(strings.Replace(strings.Replace(name, "-", "_", -1), " ", "_", -1)).Set(float64(transientState))


	checkError(err)

}
