package cli

type PipelineGenerateCommandContext struct {
	TerminalContext
	ConfigContext

	Debug bool
}

func PipelineGenerateCommand(ctx PipelineGenerateCommandContext) error {
	return nil
}
