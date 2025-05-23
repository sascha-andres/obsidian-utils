package internal

import "github.com/sascha-andres/reuse/flag"

const commonFlagPrefix = "OBS_UTIL"

// AddCommonFlagPrefixes sets a common environment variable prefix for specific flags in the application configuration.
func AddCommonFlagPrefixes() {
	flag.SetEnvPrefixForFlag("daily-folder", commonFlagPrefix)
	flag.SetEnvPrefixForFlag("folder", commonFlagPrefix)
	flag.SetEnvPrefixForFlag("print-config", commonFlagPrefix)
	flag.SetEnvPrefixForFlag("log-level", commonFlagPrefix)
}
