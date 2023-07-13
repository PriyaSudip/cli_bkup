package cli

import (
	"github.com/AlecAivazis/survey/v2"
)

type PipelineGenerateCommandContext struct {
	TerminalContext
	ConfigContext

	Debug bool
}

func PipelineGenerateCommand(ctx PipelineGenerateCommandContext) error {
	var project string
	if projectPrompt, err := ctx.projectPrompt(); err == nil {
		survey.AskOne(projectPrompt, &project, survey.WithValidator(survey.Required))
	}

	var language string
	if languagePrompt, err := ctx.languagePrompt(project); err == nil {
		survey.AskOne(languagePrompt, &language, survey.WithValidator(survey.Required))
	}

	var framework string
	if frameworkPrompt, err := ctx.frameworkPrompt(project, language); err == nil {
		survey.AskOne(frameworkPrompt, &framework, survey.WithValidator(survey.Required))
	}

	var hosting string
	if hostingPrompt, err := ctx.hostingPrompt(project, language, framework); err == nil {
		survey.AskOne(hostingPrompt, &hosting, survey.WithValidator(survey.Required))
	}

	file := ""
	prompt := &survey.Input{
		Message: "Save output to:",
		Default: ".buildkite/pipeline.yml",
	}
	survey.AskOne(prompt, &file)

	// todo: generate pipeline.yml

	return nil
}

func (ctx *PipelineGenerateCommandContext) projectPrompt() (survey.Prompt, error) {
	return &survey.Select{
		Message: "What type of project are you working on?",
		Options: []string{
			"Web application",
			"Mobile application",
			"ML model",
			"Other",
		},
	}, nil
}

func (ctx *PipelineGenerateCommandContext) languagePrompt(project string) (survey.Prompt, error) {
	// todo: fetch inputs from api.

	return &survey.Select{
		Message: "What language are you using?",
		Options: []string{
			"Node.js",
			"Ruby",
			"Python",
			"Go",
		},
	}, nil
}

func (ctx *PipelineGenerateCommandContext) frameworkPrompt(project string, language string) (survey.Prompt, error) {
	// todo: fetch inputs from api.

	return &survey.Select{
		Message: "What framework are you using?",
		Options: []string{
			"Ruby on Rails",
			"React",
			"Vue",
			"Angular",
			"Ember",
		},
	}, nil
}

func (ctx *PipelineGenerateCommandContext) hostingPrompt(project string, language string, framework string) (survey.Prompt, error) {
	// todo: fetch inputs from api.

	return &survey.Select{
		Message: "Where will you be deploying to?",
		Options: []string{
			"Digital Ocean",
			"AWS",
			"Heroku",
			"Netlify",
			"Vercel",
			"Other",
		},
	}, nil
}
