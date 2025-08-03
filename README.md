# rbxbin
[pkg.go.dev]:     https://pkg.go.dev/github.com/sewnie/rbxbin
[pkg.go.dev_img]: https://img.shields.io/badge/%E2%80%8B-reference-007d9c?logo=go&logoColor=white&style=flat-square

[![Godoc Reference][pkg.go.dev_img]][pkg.go.dev]

Go package providing routines to deploy, install and fetch Roblox programs (referred to as Binary)

Binaries WindowsStudio, MacPlayer and MacStudio do not have guranteed support, as this package
is directed towards Windows support. Feel free to contribute.

### Packages fetcher example

```go
package main

import (
	"log"

	"github.com/sewnie/rbxbin"
	cs "github.com/sewnie/rbxweb/clientsettings"
)

func main() {
	d, err := rbxbin.GetDeployment(cs.WindowsPlayer, "LIVE")
	if err != nil {
		log.Fatalln("failed to get deployment:", err)
	}

	ps, err := rbxbin.DefaultMirror.GetPackages(d)
	if err != nil {
		log.Fatalln("failed to get packages:", err)
	}

	for _, p := range ps {
		log.Println("Package:", p.Name)
	}
}
```
