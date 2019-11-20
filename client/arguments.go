package client

import ()
import "github.com/c-bata/go-prompt"

var commands = []prompt.Suggest{
	{Text: "get", Description: "Display one or many resources"},
	{Text: "describe", Description: "Show details of a specific resource or group of resources"},
	{Text: "logs", Description: "Print the logs for a container in a pod."},
	{Text: "start", Description: "start the program in jxcore"},
	// Custom command.
	{Text: "exit", Description: "Exit this program"},
}

var resourceTypes = []prompt.Suggest{
	{Text: "jxserving"}, // valid only for federation apiservers
	{Text: "db"},
	{Text: "mq"},
	{Text: "telegraf"},
}
