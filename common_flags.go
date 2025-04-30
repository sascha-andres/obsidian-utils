package obsidianutils

import "github.com/sascha-andres/reuse/flag"

func AddCommonFlagPrefixes() {
	flag.SetEnvPrefixForFlag("daily-folder", "OBS_UTIL")
	flag.SetEnvPrefixForFlag("folder", "OBS_UTIL")
	flag.SetEnvPrefixForFlag("print-config", "OBS_UTIL")
}
