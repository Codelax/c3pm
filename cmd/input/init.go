package input

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Masterminds/semver/v3"
	"github.com/gabrielcolson/c3pm/cli/config/manifest"
)

var InitSurvey = []*survey.Question{
	{
		Name:     "Name",
		Prompt:   &survey.Input{Message: "Project name"},
		Validate: survey.Required,
	},
	{
		Name: "Type",
		Prompt: &survey.Select{Message: "Project type",
			Options: []string{
				manifest.Executable.String(),
				manifest.Library.String(),
			},
		},
	},
	{
		Name:   "Description",
		Prompt: &survey.Input{Message: "Project description"},
	},
	{
		Name: "Version",
		Prompt: &survey.Input{
			Message: "Project version",
			Default: "1.0.0",
		},
		Validate: func(ans interface{}) error {
			_, err := semver.NewVersion(ans.(string))
			return err
		},
		Transform: func(ans interface{}) (newAns interface{}) {
			v, _ := semver.NewVersion(ans.(string))
			return v.String()
		},
	},
	{
		Name: "License",
		Prompt: &survey.Input{
			Message: "Project license",
			Default: "UNLICENSED",
			Help:    "You can read about code licenses on https://choosealicense.com/",
		},
	},
}

func Init() (manifest.Manifest, error) {
	man := manifest.New()
	err := survey.Ask(InitSurvey, &man, SurveyOptions...)
	return man, err
}
