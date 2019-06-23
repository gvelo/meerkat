package build

import (
	"fmt"
	"runtime"
)

// Build information. Populated at build-time.
var (
	Version   string
	Branch    string
	Tag       string
	BuildUser string
	BuildDate string
	Commit    string
	GoVersion = runtime.Version()
	Platform  = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
)
