package obsidianutils

import "github.com/sascha-andres/reuse/flag"

// AddCommonFlagPrefixes sets environment variable prefixes for common flags used across multiple configurations.
func AddCommonFlagPrefixes() {
	flag.SetEnvPrefixForFlag("daily-folder", "OBS_UTIL")
	flag.SetEnvPrefixForFlag("folder", "OBS_UTIL")
	flag.SetEnvPrefixForFlag("print-config", "OBS_UTIL")
}
