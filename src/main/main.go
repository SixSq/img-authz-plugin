// Docker Image Authorization Plugin.
// Allows docker images to be fetched from a list of authorized registries only.
// AUTHOR: Chaitanya Prakash N <cpdevws@gmail.com>
package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/docker/go-plugins-helpers/authorization"
)

const (
	defaultDockerHost = "unix:///var/run/docker.sock"
	pluginSocket      = "/run/docker/plugins/img-authz-plugin.sock"
)

var (
	flDockerHost         = flag.String("host", defaultDockerHost, "Specifies the host where docker daemon is running")
	authorizedRegistries stringslice
	authorizedNotary     = flag.String("notary", "", "Specifies the authorized image notary")
	Version              string
	Build                string
)

func main() {

	log.Println("Plugin Version:", Version, "Build: ", Build)

	// Fetch the registry from env
	authorizedRegistries := os.Getenv("REGISTRIES")

	// Convert authorized registries into a map for efficient lookup
	registries := make(map[string]bool)
	for _, registry := range strings.Split(authorizedRegistries, ",") {
		log.Println("Authorized registry:", registry)
		registries[registry] = true
	}
	log.Println("No. of authorized registries: ", len(registries))

	log.Println("Authorized notary: ", *authorizedNotary)

	// Create image authorization plugin
	plugin, err := newPlugin(*flDockerHost, registries, *authorizedNotary)
	if err != nil {
		log.Fatal(err)
	}

	// Start service handler on the local sock
	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)
	handler := authorization.NewHandler(plugin)
	if err := handler.ServeUnix(pluginSocket, gid); err != nil {
		log.Fatal(err)
	}
}
