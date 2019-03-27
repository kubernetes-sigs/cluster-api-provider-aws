package main

import (
	"flag"
	"fmt"

	"github.com/openshift/machine-api-operator/pkg/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Machine API Operator",
		Long:  `All software has versions. This is Machine API Operator's.`,
		Run:   runVersionCmd,
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersionCmd(cmd *cobra.Command, args []string) {
	flag.Set("logtostderr", "true")
	flag.Parse()

	program := "MachineAPIOperator"
	version := "v" + version.Version.String()

	fmt.Println(program, version)
}
