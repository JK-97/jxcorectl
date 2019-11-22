package client

import (
	"fmt"
	"jxcorectl/internal/debug"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
)

//启动client
func Run(version string) {
	c := NewCompleter()
	e := NewRpcExcutior(
		"http://127.0.0.1:9001",
		"root",
		"",
	)
	rpcc := e.createRpcClient()
	defer debug.Teardown()
	fmt.Printf("jxcorectl %s \n", version)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	defer fmt.Println("Bye!")

	e.status(rpcc, nil)
	p := prompt.New(
		e.Execute,
		c.completer,
		prompt.OptionTitle("jxcore-prompt: interactive jxcore client"),
		prompt.OptionPrefix("jxcorectl > "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}
