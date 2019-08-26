package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	var printers = dashboard.Printers

	var wg sync.WaitGroup
	wg.Add(len(printers))
	for i := 0; i < len(printers); i++ {
		fmt.Printf("%+v\n", printers[i])
		go func(p *Printer) {
			p.getSettings()
			p.getConnectionInfo()
			p.getTemperatureInfo()
			wg.Done()
		}(&printers[i])
	}
	wg.Wait()
	// Load Job info if Printing
	for i := 0; i < len(printers); i++ {
		fillDashboardPorts(&printers[i])
		if printers[i].Connection.State == "Printing" {
			printers[i].getJobInfo()
		}
	}

	json.NewEncoder(w).Encode(dashboard)
}

func fillDashboardPorts(p *Printer) {
	for i := 0; i < len(p.Connection.AvailablePorts); i++ {
		var found = -1
		for j := 0; j < len(dashboard.Ports); j++ {
			if dashboard.Ports[j].Name == p.Connection.AvailablePorts[i] && dashboard.Ports[j].HostKey == p.HostKey {
				found = j
			}
		}
		if found == -1 {
			dashboard.Ports = append(dashboard.Ports, Port{true, p.Connection.AvailablePorts[i], p.HostKey})
			found = 0
		}
		if dashboard.Ports[found].Name == p.Connection.Port {
			dashboard.Ports[found].Available = false
		}
	}
}