package version

import (
	"fmt"
	"runtime"
	"time"
)

var (
	Version   = "dev"
	BuildTime = time.Now().Format("2006-01-02_15:04:05_UTC")
	GoVersion = runtime.Version()
)

func GetVersion() string {
	return fmt.Sprintf("komiko %s (build: %s, go: %s)",
		Version, BuildTime, GoVersion)
}

func GetShortVersion() string {
	return Version
}
