package client

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

type customcompleter struct {
}

// NewCompleter 返customcompleter 对象
func NewCompleter() *customcompleter {
	return &customcompleter{}
}

// 用于补全的借口
func (c *customcompleter) completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return rpcCommands
	}
	args := strings.Split(d.TextBeforeCursor(), " ")
	_ = d.GetWordBeforeCursor()
	// If word before the cursor starts with "-", returns CLI flag options.

	return c.argumentsCompleter(args)
}
