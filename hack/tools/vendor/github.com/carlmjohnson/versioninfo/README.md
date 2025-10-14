# versioninfo [![GoDoc](https://godoc.org/github.com/carlmjohnson/versioninfo?status.svg)](https://godoc.org/github.com/carlmjohnson/versioninfo) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/versioninfo)](https://goreportcard.com/report/github.com/carlmjohnson/versioninfo)

Importable package that parses `debug.ReadBuildInfo()` for inclusion in your Go application.

Requires Go 1.18+ for Git revision information, but compatible with prior versions of Go.

## Examples

```go
package main

import (
    "fmt"

    "github.com/carlmjohnson/versioninfo"
)

func main() {
    fmt.Println("Version:", versioninfo.Version)
    fmt.Println("Revision:", versioninfo.Revision)
    fmt.Println("DirtyBuild:", versioninfo.DirtyBuild)
    fmt.Println("LastCommit:", versioninfo.LastCommit)
}
```

You may use the concatenated information provided by `versioninfo.Short()`:

```go
package main

import (
    "fmt"

    "github.com/carlmjohnson/versioninfo"
)

func main() {
    fmt.Println("ShortInfo:", versioninfo.Short())
}
```

Add the `-v` and `-version` flags:

```go
package main

import (
    "flag"
    "fmt"

    "github.com/carlmjohnson/versioninfo"
)

func main() {
    versioninfo.AddFlag(nil)
    flag.Parse()
    fmt.Println("done")
}
```
