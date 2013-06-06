// Copyright 2013 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"encoding/json"
	"github.com/dotcloud/docker"
)

// ListContainersOptions specify parameters to the ListContainers function.
//
// See http://goo.gl/8IMr2 for more details.
type ListContainersOptions struct {
	All    bool
	Limit  int
	Since  string
	Before string
}

// ListContainers returns a slice of containers matching the given criteria.
//
// See http://goo.gl/8IMr2 for more details.
func (c *Client) ListContainers(opts *ListContainersOptions) ([]docker.ApiContainer, error) {
	path := "/containers/ps?" + queryString(opts)
	body, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var containers []docker.ApiContainer
	err = json.Unmarshal(body, &containers)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// InspectContainer returns information about a container by its ID.
//
// See http://goo.gl/g5tpG for more details.
func (c *Client) InspectContainer(id string) (*docker.Container, error) {
	path := "/containers/" + id + "/json"
	body, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var container docker.Container
	err = json.Unmarshal(body, &container)
	if err != nil {
		return nil, err
	}
	return &container, nil
}

// CreateContainer creates a new container, returning the container instance,
// or an error in case of failure.
//
// See http://goo.gl/lcR51 for more details.
func (c *Client) CreateContainer(config *docker.Config) (*docker.Container, error) {
	body, _, err := c.do("POST", "/containers/create", config)
	if err != nil {
		return nil, err
	}
	var container docker.Container
	err = json.Unmarshal(body, &container)
	if err != nil {
		return nil, err
	}
	return &container, nil
}

// KillContainer kills a container, returning an error in case of failure.
//
// See http://goo.gl/DfPJC for more details.
func (c *Client) KillContainer(id string) error {
	path := "/containers/" + id + "/kill"
	_, _, err := c.do("POST", path, nil)
	if err != nil {
		return err
	}
	return nil
}

// RemoveContainer removes a container, returning an error in case of failure.
//
// See http://goo.gl/vCybY for more details.
func (c *Client) RemoveContainer(id string) error {
	_, _, err := c.do("DELETE", "/containers/"+id, nil)
	if err != nil {
		return err
	}
	return nil
}