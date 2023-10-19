package cli

import (
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

const (
	CircleCI      string = "CircleCI"
	GitHubActions        = "GitHub Actions"
)

type PipelineMigrateCommandContext struct {
	TerminalContext
	ConfigContext

	Debug bool
}

func PipelineMigrateCommand(ctx PipelineMigrateCommandContext) error {
	var provider string
	survey.AskOne(ctx.providerPrompt(), &provider, survey.WithValidator(survey.Required))

	var input string
	survey.AskOne(ctx.inputPrompt(provider), &input)

	output := ""
	outputPrompt := &survey.Input{
		Message: "Output:",
		Default: ".buildkite/pipeline.yml",
	}
	survey.AskOne(outputPrompt, &output)

	// TODO: insert amazing buildkite spinner here.
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	pipeline, err := ctx.transform(input, output)
	if err != nil {
		return err
	}
	s.Stop()

	ctx.Println(color.GreenString("✅ %s", pipeline.Name()))

	build := false
	buildPrompt := &survey.Confirm{
		Default: true,
		Message: "Run local build?",
	}
	survey.AskOne(buildPrompt, &build)

	if build {
		err := LocalRunCommand(LocalRunCommandContext{File: pipeline})
		if err != nil {
			return err
		}
	}

	return nil
}

func (ctx *PipelineMigrateCommandContext) transform(input, output string) (*os.File, error) {
	// FIXME: faking something happening...
	time.Sleep(1 * time.Second)

	// return nil, fmt.Errorf("oww my bones")

	// TODO: Read input file.
	// TODO: call API and write out file.

	return os.Open(output)
}

func (ctx *PipelineMigrateCommandContext) providerPrompt() survey.Prompt {
	return &survey.Select{
		Message: "Which CI provider are you migrating from?",
		Options: []string{
			CircleCI,
			GitHubActions,
		},
	}
}

func (ctx *PipelineMigrateCommandContext) inputPrompt(provider string) survey.Prompt {
	var defaultValue string
	switch provider {
	case CircleCI:
		defaultValue = ".circleci/config.yml"
	case GitHubActions:
		defaultValue = ".github/workflows/main.yml"
	}

	return &survey.Input{
		Message: "Input:",
		Default: defaultValue,
	}
}
