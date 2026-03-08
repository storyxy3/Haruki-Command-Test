package sekai

import (
	"Haruki-Command-Parser/internal/parser"
)

func makeResolvedCmd(ctx SekaiHandlerContext, module parser.TargetModule, mode string) *parser.ResolvedCommand {
	return &parser.ResolvedCommand{
		Module:    module,
		Mode:      mode,
		Query:     ctx.GetArgs(),
		Region:    ctx.Region().Id(),
		IsHelp:    ctx.Flags()["is_help"],
		IsVerbose: ctx.Flags()["is_verbose"],
		IsPreview: ctx.Flags()["is_preview"],
	}
}
