package internal

import "github.com/manifoldco/promptui"

// PromptText runs a textual prompt
func PromptText(label, defaultValue string, val func(string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	if nil != val {
		prompt.Validate = val
	}
	return prompt.Run()
}
