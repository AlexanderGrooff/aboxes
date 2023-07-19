package cmd

import (
	"log"

	"github.com/appleboy/easyssh-proxy"
	"github.com/kevinburke/ssh_config"
)

func executeCommands(targets []string, commands []string, outputFile string) {
	for _, target := range targets {
		ssh := getConfigForHost(target)
		for _, command := range commands {
			stdout, stderr, _, err := ssh.Run(command)
			if err != nil {
				panic(err)
			}
			log.Println(stdout, stderr)
		}
	}
}

func getConfigForHost(target string) *easyssh.MakeConfig {
	ssh := &easyssh.MakeConfig{
		User:   ssh_config.Get(target, "User"),
		Server: ssh_config.Get(target, "HostName"),
		Port:   ssh_config.Get(target, "Port"),
	}
	return ssh
}
