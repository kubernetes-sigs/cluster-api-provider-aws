//go:build go1.12

package versioninfo

import "runtime/debug"

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	if info.Main.Version != "" {
		Version = info.Main.Version
	}
}
