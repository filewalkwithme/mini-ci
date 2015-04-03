package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func main() {
	minici, _ := ioutil.ReadFile(".mini-ci.yml")
	fmt.Println("---minideploy---")

	lines := strings.Split(string(minici), "\n")

	section := ""
	var hostDeploy string
	var folderDeploy string
	var files []string
	var commands []string

	if strings.HasSuffix(lines[0], "deploy-in-this-host:") {
		for _, v := range lines {
			//which section are we?
			if v == "deploy-in-this-host:" {
				section = "host"
				continue
			}

			if v == "these-files:" {
				section = "files"
				continue
			}

			if v == "then-run-these-commands:" {
				section = "commands"
				continue
			}

			if len(strings.Trim(v, "")) == 0 {
				continue
			}

			if section == "host" {
				hostDeploy = strings.Split(v, ":")[0]
				folderDeploy = strings.Split(v, ":")[1]
			}

			if section == "files" {
				files = append(files, v)
			}

			if section == "commands" {
				commands = append(commands, v)
			}

		}

		cmdMkdir := "mkdir " + folderDeploy
		out, err := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-i", "/root/.ssh/id_rsa", hostDeploy, cmdMkdir).CombinedOutput()
		if err == nil {
			for _, file := range files {
				out, err = exec.Command("scp", "-r", "-o", "StrictHostKeyChecking=no", "-i", "/root/.ssh/id_rsa", file, hostDeploy+":"+folderDeploy+"/"+file).CombinedOutput()
				if err != nil {
					fmt.Printf("err: %s [%s]\n", string(out), err)
					return
				}
			}
		} else {
			fmt.Printf("err: %s [%s]\n", string(out), err)
			return
		}

		script := ""
		for _, cmd := range commands {
			script = script + cmd + ";"
		}

		fmt.Printf("%s\n", script)
		out, err = exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-i", "/root/.ssh/id_rsa", hostDeploy, script).CombinedOutput()
		if err != nil {
			fmt.Printf("err: %s [%s]\n", string(out), err)
			return
		}
	}
	fmt.Println("---end minideploy---")
}
