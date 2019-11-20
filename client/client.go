package client

import (
	"fmt"
	"jxcorectl/internal/debug"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
)

func customcompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Run(version string) {

	defer debug.Teardown()
	fmt.Printf("kube-prompt %s \n", version)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	defer fmt.Println("Bye!")
	p := prompt.New(
		Executor,
		customcompleter,
		prompt.OptionTitle("jxcore-prompt: interactive jxcore client"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}
