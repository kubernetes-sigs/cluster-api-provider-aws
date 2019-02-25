package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/openshift/suites"
)

func main() {
	flag.Parse()
	if err := suites.Run(); err != nil {
		glog.Fatal(err)
	}
}
