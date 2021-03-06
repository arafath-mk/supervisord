package main

import (
	"fmt"
	"strings"
)

type CtlCommand struct {
	ServerUrl string `short:"s" long:"serverurl" description:"URL on which supervisord server is listening" default:"http://localhost:9001"`
}

var ctlCommand CtlCommand

func (x *CtlCommand) Execute(args []string) error {
	if len(args) == 0 {
		return nil
	}

	rpcc := NewXmlRPCClient(x.ServerUrl)

	verb, processes := args[0], args[1:]
	hasProcesses := len(processes) > 0
	processesMap := make(map[string]bool)
	for _, process := range processes {
		processesMap[strings.ToLower(process)] = true
	}

	switch verb {

	////////////////////////////////////////////////////////////////////////////////
	// STATUS
	////////////////////////////////////////////////////////////////////////////////
	case "status":
		if reply, err := rpcc.GetAllProcessInfo(); err == nil {
			for _, pinfo := range reply.Value {
				name := strings.ToLower(pinfo.Name)
				description := pinfo.Description
				if strings.ToLower(description) == "<string></string>" {
					description = ""
				}
				if !hasProcesses || processesMap[name] {
					fmt.Printf("%-33s%-10s%s\n", name, pinfo.Statename, description)
				}
			}
		}

	////////////////////////////////////////////////////////////////////////////////
	// START or STOP
	////////////////////////////////////////////////////////////////////////////////
	case "start", "stop":
		state := map[string]string{
			"start": "started",
			"stop":  "stopped",
		}
		for _, pname := range processes {
			if reply, err := rpcc.ChangeProcessState(verb, pname); err == nil {
				fmt.Printf("%s: ", pname)
				if !reply.Value {
					fmt.Printf("not ")
				}
				fmt.Printf("%s\n", state[verb])
			} else {
				fmt.Printf("%s: failed [%v]\n", pname, err)
			}
		}

	////////////////////////////////////////////////////////////////////////////////
	// SHUTDOWN
	////////////////////////////////////////////////////////////////////////////////
	case "shutdown":
		if reply, err := rpcc.Shutdown(); err == nil {
			if reply.Value {
				fmt.Printf("Shut Down\n")
			} else {
				fmt.Printf("Hmmm! Something gone wrong?!\n")
			}
		}

	default:
		fmt.Println("unknown command")
	}

	return nil
}

func init() {
	parser.AddCommand("ctl",
		"Control a running daemon",
		"The ctl subcommand resembles supervisorctl command of original daemon.",
		&ctlCommand)
}
