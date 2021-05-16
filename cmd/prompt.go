package cmd

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/manifoldco/promptui"
)

func promptModule() (string, error) {
	moduleValidate :=
		func(input string) error {
			return validation.Validate(input,
				validation.Required,
			)
		}

	modulePrompt := promptui.Prompt{
		Label:    "module name",
		Validate: moduleValidate,
	}

	return modulePrompt.Run()
}

func promptGoVersion() (string, error) {
	goVersionPrompt := promptui.Select{
		Label: "go version",
		Items: []string{"1.14", "1.15", "1.16"},
	}

	_, r, err := goVersionPrompt.Run()
	return r, err
}
