package client

import (
	"fmt"
	"jxcorectl/rpc"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/ochinchina/supervisord/types"
)

type RpcExector struct {
	ServerUrl string `short:"s" long:"serverurl" description:"URL on which supervisord server is listening"`
	User      string `short:"u" long:"user" description:"the user name"`
	Password  string `short:"P" long:"password" description:"the password"`
	Verbose   bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
}

type StatusCommand struct {
}

type StartCommand struct {
}

type StopCommand struct {
}

type RestartCommand struct {
}

type ShutdownCommand struct {
}

type ReloadCommand struct {
}

type PidCommand struct {
}

type SignalCommand struct {
}

type TailCommand struct {
}

var rpcExector RpcExector
var statusCommand StatusCommand
var startCommand StartCommand
var stopCommand StopCommand
var restartCommand RestartCommand
var shutdownCommand ShutdownCommand
var reloadCommand ReloadCommand
var pidCommand PidCommand
var signalCommand SignalCommand
var tailCommand TailCommand

func NewRpcExcutior(serveurl, user, password string) *RpcExector {
	return &RpcExector{
		ServerUrl: serveurl,
		User:      user,
		Password:  password,
	}
}

func (x *RpcExector) getServerUrl() string {
	return x.ServerUrl
}
func (x *RpcExector) getUser() string {
	return x.User
}
func (x *RpcExector) getPassword() string {

	return x.Password
}

func (x *RpcExector) createRpcClient() *rpc.XmlRPCClient {
	rpcc := rpc.NewXmlRPCClient(x.getServerUrl(), x.Verbose)
	rpcc.SetUser(x.getUser())
	rpcc.SetPassword(x.getPassword())
	return rpcc
}

func (x *RpcExector) Execute(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	} else if s == "quit" || s == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}
	args := strings.Split(s, " ")
	length := len(args)
	rpcc := x.createRpcClient()
	firstword := args[0]

	switch firstword {
	case "status":
		x.status(rpcc, args[1:])
	case "start", "stop":
		if length >= 1 {
			x.startStopProcesses(rpcc, firstword, args[1:])
			return
		}
	case "shutdown":
		x.shutdown(rpcc)
	case "reload":
		x.reload(rpcc)
	case "signal":
		sig_name, processes := args[1], args[2:]
		x.signal(rpcc, sig_name, processes)
	case "tail":
		if length == 3 {
			thirdword := args[2]
			// secondword := args[1]
			switch thirdword {
			case "stderr":
				x.tailProcessStderrLog(rpcc, args[1:])
			case "stdout":
				x.tailProcessStdoutLog(rpcc, args[1:])
			}
		}
	case "pid":
		if length == 2 {
			x.getPid(rpcc, args[1])
			return
		}

	default:
		fmt.Println("unknown command")
	}

	return
}

// get the status of processes
func (x *RpcExector) status(rpcc *rpc.XmlRPCClient, processes []string) {
	processesMap := make(map[string]bool)
	for _, process := range processes {
		processesMap[process] = true
	}
	if reply, err := rpcc.GetAllProcessInfo(); err == nil {
		x.showProcessInfo(&reply, processesMap)
	} else {
		os.Exit(1)
	}
}

// start or stop the processes
// verb must be: start or stop
func (x *RpcExector) startStopProcesses(rpcc *rpc.XmlRPCClient, verb string, processes []string) {
	state := map[string]string{
		"start": "started",
		"stop":  "stopped",
	}
	x._startStopProcesses(rpcc, verb, processes, state[verb], true)
}

func (x *RpcExector) _startStopProcesses(rpcc *rpc.XmlRPCClient, verb string, processes []string, state string, showProcessInfo bool) {
	if len(processes) <= 0 {
		fmt.Printf("Please specify process for %s\n", verb)
	}
	for _, pname := range processes {
		if pname == "all" {
			reply, err := rpcc.ChangeAllProcessState(verb)
			if err == nil {
				if showProcessInfo {
					x.showProcessInfo(&reply, make(map[string]bool))
				}
			} else {
				fmt.Printf("Fail to change all process state to %s", state)
			}
		} else {
			if reply, err := rpcc.ChangeProcessState(verb, pname); err == nil {
				if showProcessInfo {
					fmt.Printf("%s: ", pname)
					if !reply.Value {
						fmt.Printf("not ")
					}
					fmt.Printf("%s\n", state)
				}
			} else {
				fmt.Printf("%s: failed [%v]\n", pname, err)
				os.Exit(1)
			}
		}
	}
}

func (x *RpcExector) restartProcesses(rpcc *rpc.XmlRPCClient, processes []string) {
	x._startStopProcesses(rpcc, "stop", processes, "stopped", false)
	x._startStopProcesses(rpcc, "start", processes, "restarted", true)
}

// shutdown the supervisord
func (x *RpcExector) shutdown(rpcc *rpc.XmlRPCClient) {
	if reply, err := rpcc.Shutdown(); err == nil {
		if reply.Value {
			fmt.Printf("Shut Down\n")
		} else {
			fmt.Printf("Hmmm! Something gone wrong?!\n")
		}
	} else {
		os.Exit(1)
	}
}

// reload all the programs in the supervisord
func (x *RpcExector) reload(rpcc *rpc.XmlRPCClient) {
	if reply, err := rpcc.ReloadConfig(); err == nil {

		if len(reply.AddedGroup) > 0 {
			fmt.Printf("Added Groups: %s\n", strings.Join(reply.AddedGroup, ","))
		}
		if len(reply.ChangedGroup) > 0 {
			fmt.Printf("Changed Groups: %s\n", strings.Join(reply.ChangedGroup, ","))
		}
		if len(reply.RemovedGroup) > 0 {
			fmt.Printf("Removed Groups: %s\n", strings.Join(reply.RemovedGroup, ","))
		}
	} else {
		os.Exit(1)
	}
}

// send signal to one or more processes
func (x *RpcExector) signal(rpcc *rpc.XmlRPCClient, sig_name string, processes []string) {
	for _, process := range processes {
		if process == "all" {
			reply, err := rpcc.SignalAll(process)
			if err == nil {
				x.showProcessInfo(&reply, make(map[string]bool))
			} else {
				fmt.Printf("Fail to send signal %s to all process", sig_name)
				os.Exit(1)
			}
		} else {
			reply, err := rpcc.SignalProcess(sig_name, process)
			if err == nil && reply.Success {
				fmt.Printf("Succeed to send signal %s to process %s\n", sig_name, process)
			} else {
				fmt.Printf("Fail to send signal %s to process %s\n", sig_name, process)
				os.Exit(1)
			}
		}
	}
}

// get the pid of running program
func (x *RpcExector) getPid(rpcc *rpc.XmlRPCClient, process string) {
	procInfo, err := rpcc.GetProcessInfo(process)
	if err != nil {
		fmt.Printf("program '%s' not found\n", process)
		os.Exit(1)
	} else {
		fmt.Printf("%d\n", procInfo.Pid)
	}
}

// check if group name should be displayed
func (x *RpcExector) showGroupName() bool {
	val, ok := os.LookupEnv("SUPERVISOR_GROUP_DISPLAY")
	if !ok {
		return false
	}

	val = strings.ToLower(val)
	return val == "yes" || val == "true" || val == "y" || val == "t" || val == "1"
}

// tail the process stderr log
func (x *RpcExector) tailProcessStderrLog(rpcc *rpc.XmlRPCClient, args []string) {
	process := args[0]

	processLogReadInfo := rpc.ProcessLogReadInfo{
		Name:   process,
		Offset: 0,
		Length: 100,
	}
	if reply, err := rpcc.TailProcessStderrLog(processLogReadInfo); err == nil {
		x.showLogFlow(&reply)
	} else {
		fmt.Println(err)
		os.Exit(1)
	}

}

// tail the process stdout log
func (x *RpcExector) tailProcessStdoutLog(rpcc *rpc.XmlRPCClient, args []string) {
	process := args[0]

	processLogReadInfo := rpc.ProcessLogReadInfo{
		Name:   process,
		Offset: 0,
		Length: 100,
	}
	if reply, err := rpcc.TailProcessStdoutLog(processLogReadInfo); err == nil {
		x.showLogFlow(&reply)
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (x *RpcExector) showLogFlow(processTailLog *rpc.ProcessTailLog) {
	if processTailLog.LogData != "" {
		fmt.Println(processTailLog.LogData)
	}

}
func (x *RpcExector) showProcessInfo(reply *rpc.AllProcessInfoReply, processesMap map[string]bool) {
	for _, pinfo := range reply.Value {
		description := pinfo.Description
		if strings.ToLower(description) == "<string></string>" {
			description = ""
		}
		if x.inProcessMap(&pinfo, processesMap) {
			processName := pinfo.GetFullName()
			if !x.showGroupName() {
				processName = pinfo.Name
			}
			fmt.Printf("%s%-33s%-10s%s%s\n", x.getANSIColor(pinfo.Statename), processName, pinfo.Statename, description, "\x1b[0m")
		}
	}
}

func (x *RpcExector) inProcessMap(procInfo *types.ProcessInfo, processesMap map[string]bool) bool {
	if len(processesMap) <= 0 {
		return true
	}
	for procName, _ := range processesMap {
		if procName == procInfo.Name || procName == procInfo.GetFullName() {
			return true
		}

		// check the wildcast '*'
		pos := strings.Index(procName, ":")
		if pos != -1 {
			groupName := procName[0:pos]
			programName := procName[pos+1:]
			if programName == "*" && groupName == procInfo.Group {
				return true
			}
		}
	}
	return false
}

func (x *RpcExector) getANSIColor(statename string) string {
	if statename == "RUNNING" {
		// green
		return "\x1b[0;32m"
	} else if statename == "BACKOFF" || statename == "FATAL" {
		// red
		return "\x1b[0;31m"
	} else {
		// yellow
		return "\x1b[1;33m"
	}
}

func (sc *StatusCommand) Execute(args []string) error {
	rpcExector.status(rpcExector.createRpcClient(), args)
	return nil
}

func (sc *StartCommand) Execute(args []string) error {
	rpcExector.startStopProcesses(rpcExector.createRpcClient(), "start", args)
	return nil
}

func (sc *StopCommand) Execute(args []string) error {
	rpcExector.startStopProcesses(rpcExector.createRpcClient(), "stop", args)
	return nil
}

func (rc *RestartCommand) Execute(args []string) error {
	rpcExector.restartProcesses(rpcExector.createRpcClient(), args)
	return nil
}

func (sc *ShutdownCommand) Execute(args []string) error {
	rpcExector.shutdown(rpcExector.createRpcClient())
	return nil
}

func (rc *ReloadCommand) Execute(args []string) error {
	rpcExector.reload(rpcExector.createRpcClient())
	return nil
}

func (rc *SignalCommand) Execute(args []string) error {
	sig_name, processes := args[0], args[1:]
	rpcExector.signal(rpcExector.createRpcClient(), sig_name, processes)
	return nil
}

func (pc *PidCommand) Execute(args []string) error {
	rpcExector.getPid(rpcExector.createRpcClient(), args[0])
	return nil
}

func init() {
	parser := flags.NewParser(nil, flags.Default & ^flags.PrintErrors)
	ctlCmd, _ := parser.AddCommand("ctl",
		"Control a running daemon",
		"The ctl subcommand resembles supervisorctl command of original daemon.",
		&rpcExector)
	ctlCmd.AddCommand("status",
		"show program status",
		"show all or some program status",
		&statusCommand)
	ctlCmd.AddCommand("start",
		"start programs",
		"start one or more programs",
		&startCommand)
	ctlCmd.AddCommand("stop",
		"stop programs",
		"stop one or more programs",
		&stopCommand)
	ctlCmd.AddCommand("restart",
		"restart programs",
		"restart one or more programs",
		&restartCommand)
	ctlCmd.AddCommand("shutdown",
		"shutdown supervisord",
		"shutdown supervisord",
		&shutdownCommand)
	ctlCmd.AddCommand("reload",
		"reload the programs",
		"reload the programs",
		&reloadCommand)
	ctlCmd.AddCommand("signal",
		"send signal to program",
		"send signal to program",
		&signalCommand)
	ctlCmd.AddCommand("pid",
		"get the pid of specified program",
		"get the pid of specified program",
		&pidCommand)

}
