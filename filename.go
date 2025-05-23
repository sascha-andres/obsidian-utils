package obsidianutils

import (
	"fmt"
	"path"
	"strings"
	"time"
)

// replacements is a map containing character replacements for German umlauts and other special characters.
var replacements = map[string]string{
	"ä": "ae",
	"ö": "oe",
	"ü": "ue",
	"Ä": "Ae",
	"Ö": "Oe",
	"Ü": "Ue",
	"ß": "ss",
	":": "",
}

// CreateFileName generates a file name for a note based on the provided title and appointment time. It applies
// specified character replacements to the title and prefixes the file with the appointment date if the
// noDatePrefix flag is not set. The generated file name is returned as a string.
func CreateFileName(folder, localTitle string, noDatePrefix bool, timeForPrefix time.Time) (string, error) {
	fixed := localTitle
	for k, v := range replacements {
		fixed = strings.ReplaceAll(fixed, k, v)
	}
	fName := fmt.Sprintf("%s.md", fixed)
	if !noDatePrefix {
		fName = fmt.Sprintf("%s %s", timeForPrefix.Format("2006-01-02"), fName)
	}
	return ApplyDirectoryPlaceHolder(path.Join(folder, fName))
}
