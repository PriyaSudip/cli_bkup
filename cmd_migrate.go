package cli

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var providers map[string]string = map[string]string{
	"CircleCI":       "circleci.yml",
	"GitHub Actions": "github.yml",
}

type PipelineMigrateCommandContext struct {
	TerminalContext
	ConfigContext

	MigrateAPIEndpoint string

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
	file, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(input))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ctx.MigrateAPIEndpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "text/yaml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resBody := &bytes.Buffer{}
	_, err = resBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	fmt.Println(resBody)

	err = os.WriteFile(output, resBody.Bytes(), 0644)
	if err != nil {
		return nil, err
	}

	return os.Open(output)
}

func (ctx *PipelineMigrateCommandContext) providerPrompt() survey.Prompt {
	var options []string
	for k := range providers {
		options = append(options, k)
	}

	return &survey.Select{
		Message: "Which CI provider are you migrating from?",
		Options: options,
	}
}

func (ctx *PipelineMigrateCommandContext) inputPrompt(provider string) survey.Prompt {
	defaultValue, ok := providers[provider]
	if !ok {
		defaultValue = ""
	}

	return &survey.Input{
		Message: "Input:",
		Default: defaultValue,
	}
}
