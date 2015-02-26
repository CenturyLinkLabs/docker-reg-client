/*
Uses the Docker Registry API to reverse engineer a Dockerfile from a named
image. The image name (including tag) should be specified as the first command
line argument when executing the program.
*/
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

const NOP_PREFIX = "#(nop) "

func main() {
	repo := "mysql"
	tag := "latest"

	// Check for image name as command line argument
	if len(os.Args) > 1 {
		// Separate repo and tag
		parts := strings.Split(os.Args[1], ":")
		repo = parts[0]

		if len(parts) > 1 {
			tag = parts[1]
		}
	}

	c := registry.NewClient()

	// First step is to retrieve a read-only token from the Hub
	auth, err := c.Hub.GetReadToken(repo)
	if err != nil {
		panic(err)
	}

	// Get image ID for tag
	imageID, err := c.Repository.GetImageID(repo, tag, auth)
	if err != nil {
		panic(err)
	}

	// Retrieve all ancestors of given layer ID
	layers, err := c.Image.GetAncestry(imageID, auth)
	if err != nil {
		panic(err)
	}

	var commands []string

	// Iterate over all ancestors
	for _, layerID := range layers {
		// Retrieve metadata for given layer ID
		m, err := c.Image.GetMetadata(layerID, auth)
		if err != nil {
			panic(err)
		}

		// Format the command string
		cmdParts := m.ContainerConfig.Cmd
		if len(cmdParts) == 3 {
			cmd := cmdParts[2]

			if strings.HasPrefix(cmd, NOP_PREFIX) {
				commands = append(commands, strings.Split(cmd, NOP_PREFIX)[1])
			} else {
				commands = append(commands, fmt.Sprintf("RUN %s", cmd))
			}
		}
	}

	// Iterate over commands in reverse order and print them
	for i := len(commands) - 1; i >= 0; i-- {
		fmt.Println(squeeze(commands[i]))
	}
}

// Replace consecutive spaces with a single space
func squeeze(s string) string {
	re := regexp.MustCompile("\\s+")
	return re.ReplaceAllString(s, " ")
}
