package main

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/divolgin/kurl-aid/pkg/checks"
	"github.com/divolgin/kurl-aid/pkg/log"
	"github.com/divolgin/kurl-aid/pkg/system"
	"github.com/manifoldco/promptui"
)

func main() {
	if canRunAsRoot() {
		prompt := promptui.Prompt{
			Label:     "This program will run system commands as root. Do you want to proceed",
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			os.Exit(1)
		}
	} else {
		prompt := promptui.Prompt{
			Label:     "Running without root access limits what ations can be taken. Do you want to proceed anyway",
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			os.Exit(1)
		}
	}

	checks := []checks.Check{
		checks.Path{},
		checks.KubeConfig{},
		checks.FirewallD{},
		checks.IPTables{},
		checks.Kubectl{},
		checks.Kubelet{},
		checks.SELinux{},
	}

	ctx := context.TODO()

	logFile, err := ioutil.TempFile("", "kurl-aid-")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	postCheckCommands := []string{}
	for _, check := range checks {
		log.Printf("Running %s check", check.Name())

		log.IndentMore()

		resutl := check.Run(ctx, logFile)
		for _, err := range resutl.Errors {
			log.Printf("Failed: %v", err)
		}
		postCheckCommands = append(postCheckCommands, resutl.PostCheck...)

		log.IndentLess()
	}

	if len(postCheckCommands) > 0 {
		log.Printf("The following commands may need to be executed in the current bash session")
		log.IndentMore()
		for _, cmd := range postCheckCommands {
			log.Printf("%s", cmd)
		}
		log.IndentLess()
	}

	log.Printf("Commands executed in this session are logged in %s", logFile.Name())
}

func canRunAsRoot() bool {
	if os.Geteuid() == 0 {
		return true
	}

	stdout, _, err := system.Exec("sudo", "echo", "test")
	if err != nil {
		return false
	}

	if string(stdout) == "test\n" { // echo adds a line break
		return true
	}

	return false
}
