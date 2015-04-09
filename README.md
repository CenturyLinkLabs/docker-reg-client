# docker-reg-client
[![Circle CI](https://circleci.com/gh/CenturyLinkLabs/docker-reg-client.svg?style=svg)](https://circleci.com/gh/CenturyLinkLabs/docker-reg-client)
[![GoDoc](http://godoc.org/github.com/CenturyLinkLabs/docker-reg-client/registry?status.png)](http://godoc.org/github.com/CenturyLinkLabs/docker-reg-client/registry)

docker-reg-client is an API wrapper for the [Docker Registry v1 API](https://docs.docker.com/reference/api/registry_api/) written in Go.

For detailed documentation, see the [GoDocs](http://godoc.org/github.com/CenturyLinkLabs/docker-reg-client/registry).

## Example

    package main

    import (
      "fmt"
      "github.com/CenturyLinkLabs/docker-reg-client/registry"
    )

    func main() {
      c := registry.NewClient()

      auth, err := c.Hub.GetReadToken("ubuntu")
      if err != nil {
        panic(err)
      }

      tags, err := c.Repository.ListTags("ubuntu", auth)
      if err != nil {
        panic(err)
      }

      fmt.Printf("%#v", tags)
    }

## Write token caveat

Please note that the `GetWriteToken()` only works on **non automated build**s.
