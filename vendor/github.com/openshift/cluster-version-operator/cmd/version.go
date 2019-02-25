package main

import (
	"flag"
	"fmt"

	"github.com/openshift/cluster-version-operator/pkg/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Cluster Version Operator",
		Long:  `All software has versions. This is Cluster Version Operator's.`,
		Run:   runVersionCmd,
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersionCmd(cmd *cobra.Command, args []string) {
	flag.Set("logtostderr", "true")
	flag.Parse()

	fmt.Println(version.String)
}
