package internal

import (
	stdflag "flag"
	"fmt"

	"go.livingit.de/reuse/flag"
)

const commonFlagPrefix = "OBS_UTIL"

// AddCommonFlagPrefixes sets a common environment variable prefix for specific flags in the application configuration.
func AddCommonFlagPrefixes() {
	flag.SetEnvPrefixForFlag("daily-folder", commonFlagPrefix)
	flag.SetEnvPrefixForFlag("folder", commonFlagPrefix)
	flag.SetEnvPrefixForFlag("print-config", commonFlagPrefix)
	flag.SetEnvPrefixForFlag("log-level", commonFlagPrefix)
}

// PrintFlags prints all parsed flags and their values.
func PrintFlags() {
	flag.VisitAll(func(f *stdflag.Flag) {
		fmt.Printf("%s=%s\n", f.Name, f.Value.String())
	})
}
