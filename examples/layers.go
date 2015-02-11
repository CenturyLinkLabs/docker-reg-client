package main

import (
	"fmt"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

var (
	c *registry.Client
)

func main() {
	repo := "centurylink/image-graph"
	c = registry.NewClient()

	// First step is to retrieve a read-only token from the Hub
	auth, err := c.Hub.GetReadToken(repo)
	if err != nil {
		panic(err)
	}

	dumpRepo(repo, auth)
}

func dumpRepo(repo string, auth registry.Authenticator) {
	// List the tags available in the named repo
	tags, err := c.Repository.ListTags(repo, auth)
	if err != nil {
		panic(err)
	}

	for tag, id := range tags {
		fmt.Printf("# %s\n", tag)
		dumpImageLayers(id, auth)
	}
}

func dumpImageLayers(id string, auth registry.Authenticator) {
	// Retrieve all ancestors of given layer ID
	layers, err := c.Image.GetAncestry(id, auth)
	if err != nil {
		panic(err)
	}

	for _, layerID := range layers {
		dumpImageLayer(layerID, auth)
	}
}

func dumpImageLayer(id string, auth registry.Authenticator) {
	// Retrieve metadata for given layer ID
	m, err := c.Image.GetMetadata(id, auth)
	if err != nil {
		panic(err)
	}

	fmt.Printf("  %s\n", m.ID)
	fmt.Printf("    - %s\n", m.ContainerConfig.Cmd)
}
