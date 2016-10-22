package main

import (
	"fmt"
	ecsdeploy "github.com/dakatsuka/ecs-deploy"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var version string

func main() {
	var (
		cluster   = kingpin.Flag("cluster", "Set cluster name").Required().String()
		service   = kingpin.Flag("service", "Set service name").Required().String()
		container = kingpin.Flag("container", "Set container name").Required().String()
		image     = kingpin.Flag("image", "Set image").Required().String()
	)

	kingpin.Version(version)
	kingpin.Parse()

	err := ecsdeploy.Run(*cluster, *service, *container, *image)

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
