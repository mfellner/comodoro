// Copyright 2014 CoreOS, Inc.
// Copyright 2015 Maximilian Fellner <https://github.com/mfellner>
//
// Includes modified parts of https://github.com/coreos/fleet.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/fleet/client"
	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/pkg"
	"github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"
)

// FleetClient encapsulates the fleet client.
type FleetClient struct {
	api client.API
}

// NewFleetClient returns a new fleet client.
func NewFleetClient(endpoint string) *FleetClient {
	httpClient, err := getHTTPClient(endpoint)
	if err != nil {
		log.Fatalf("Unable to initialize fleet client: %v", err)
	}

	log.WithFields(log.Fields{
		"endpoint": endpoint,
	}).Info("Initialized fleet client")

	fleetClient := &FleetClient{
		api: httpClient,
	}

	return fleetClient
}

func getHTTPClient(endpoint string) (client.API, error) {
	ep, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if len(ep.Scheme) == 0 {
		return nil, errors.New("URL scheme undefined")
	}

	dialFunc := net.Dial

	if ep.Scheme == "unix" || ep.Scheme == "file" {
		// This commonly happens if the user misses the leading slash after the scheme.
		// For example, "unix://var/run/fleet.sock" would be parsed as host "var".
		if len(ep.Host) > 0 {
			return nil, fmt.Errorf("unable to connect to host %q with scheme %q", ep.Host, ep.Scheme)
		}

		// The Path field is only used for dialing and should not be used when
		// building any further HTTP requests.
		sockPath := ep.Path
		ep.Path = ""

		// http.Client will dial the unix socket directly but it does
		// not natively support dialing a unix domain socket, so the
		// dial function must be overridden.
		dialFunc = func(string, string) (net.Conn, error) {
			return net.Dial("unix", sockPath)
		}

		// http.Client doesn't support the schemes "unix" or "file", but it
		// is safe to use "http" as dialFunc ignores it anyway.
		ep.Scheme = "http"

		// The Host field is not used for dialing, but will be exposed in debug logs.
		ep.Host = "domain-sock"
	}

	// TODO: add TLS support. Params: CAFile, CertFile, KeyFile.
	tlsConfig, err := pkg.ReadTLSConfigFiles("", "", "")
	if err != nil {
		return nil, err
	}

	// Use a regular, non-logging http.Transport.
	trans := http.Transport{
		Dial:            dialFunc,
		TLSClientConfig: tlsConfig,
	}

	hc := http.Client{
		Transport: &trans,
	}

	return client.NewHTTPClient(&hc, *ep)
}

// CreateUnit submits a new fleet unit.
func (fc *FleetClient) CreateUnit(name string, uf *unit.UnitFile) (*schema.Unit, error) {
	if uf == nil {
		return nil, fmt.Errorf("nil unit provided")
	}
	u := schema.Unit{
		Name:    name,
		Options: schema.MapUnitFileToSchemaUnitOptions(uf),
	}

	j := &job.Job{Unit: *uf}
	if err := j.ValidateRequirements(); err != nil {
		log.Warning("Unit %s: %v", name, err)
	}
	err := fc.api.CreateUnit(&u)
	if err != nil {
		return nil, fmt.Errorf("failed creating unit %s: %v", name, err)
	}

	log.Debug("Created Unit(%s) in Registry", name)
	return &u, nil
}

// ListMachines lists the active CoreOS nodes in the fleet.
func (fc *FleetClient) ListMachines() ([]machine.MachineState, error) {
	return fc.api.Machines()
}
