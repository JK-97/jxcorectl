package client

import (
	"github.com/c-bata/go-prompt"
)

// firstword
var rpcCommands = []prompt.Suggest{
	// {Text: "restart", Description: "shutdonw jxocre subprocess"},
	{Text: "status", Description: "Print the logs for a container in a pod."},
	{Text: "tail", Description: "tail the process stderr or stdout"},
	{Text: "stop", Description: "stop the program"},
	{Text: "start", Description: "start the program"},
	{Text: "pid", Description: "get pid by name"},
	{Text: "shutdown", Description: "shutdown all processes include jxcored"},
	{Text: "signal", Description: "start the program in which jxcore"},
	{Text: "reload", Description: "start the program in which jxcore"},
	// Custom command.

}

//secondword
var processTypes = []prompt.Suggest{
	{Text: "jxserving"}, // valid only for federation apiservers
	{Text: "db"},
	{Text: "mq"},
	{Text: "tsdb"},
	{Text: "telegraf"},
	{Text: "watchdog"},
	{Text: "all"},
	{Text: "test"},
}

var logLevel = []prompt.Suggest{
	{Text: "stdout"}, // valid only for federation apiservers
	{Text: "stderr"},
}

var configFile = []prompt.Suggest{
	{Text: "dnsmasqConf", Description: "configFile for dnsmasq"}, // valid only for federation apiservers
	{Text: "initFile", Description: "initFile build by jxcore bootstrap "},
	{Text: "dnsmasqHost", Description: "hostFile for dnsmasq"},
	{Text: "dnsmasqResolv", Description: "resolvFile for dnsmasq"},
}

// 将多个suggest队列融合
func Merge(firstGround []prompt.Suggest, otherGround ...[]prompt.Suggest) []prompt.Suggest {
	for _, Ground := range otherGround {
		for _, suggest := range Ground {
			firstGround = append(firstGround, suggest)
		}
	}
	return firstGround
}

// 提示参数主逻辑
func (c *customcompleter) argumentsCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(rpcCommands, args[0], true)
	}

	firstword := args[0]
	switch firstword {
	case "status", "tail", "stop", "start", "pid", "shutdown":
		secondword := args[1]
		if len(args) == 2 {
			return prompt.FilterHasPrefix(processTypes, secondword, true)
		}
		// thirdword := args[2]shutdown
		if len(args) == 3 {
			switch firstword {
			case "tail":
				return logLevel
			}
		}

	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}
