package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/openshift/cluster-version-operator/pkg/payload"
)

var (
	imageCmd = &cobra.Command{
		Use:     "image",
		Short:   "Returns image for requested short-name from UpdatePayload",
		Long:    "",
		Example: "%[1] image <short-name>",
		Run:     runImageCmd,
	}
)

func init() {
	rootCmd.AddCommand(imageCmd)
}

func runImageCmd(cmd *cobra.Command, args []string) {
	flag.Set("logtostderr", "true")
	flag.Parse()

	if len(args) == 0 {
		glog.Fatalf("missing command line argument short-name")
	}
	image, err := payload.ImageForShortName(args[0])
	if err != nil {
		glog.Fatalf("error: %v", err)
	}
	fmt.Printf(image)
}
