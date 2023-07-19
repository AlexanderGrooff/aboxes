package cmd

import (
	"strings"

	"github.com/appleboy/easyssh-proxy"
	"github.com/kevinburke/ssh_config"
	log "github.com/sirupsen/logrus"
)

func executeCommands(targets []string, commands []string, outputFile string) {
	for _, target := range targets {
		ssh := getConfigForHost(target)
		// TODO: Use same SSH connection for all commands
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
	proxy := ssh_config.Get(target, "ProxyJump")
	if proxy != "" {
		log.Info("Using proxy ", proxy, " for target ", target)
		return &easyssh.MakeConfig{
			User:   ssh_config.Get(target, "User"),
			Server: fillSSHConfigHostname(target, ssh_config.Get(target, "HostName")),
			Port:   ssh_config.Get(target, "Port"),
			Proxy:  *makeConfigToDefaultConfig(getConfigForHost(proxy)),
		}
	} else {
		log.Info("No proxy found for target ", target)
		return &easyssh.MakeConfig{
			User:   ssh_config.Get(target, "User"),
			Server: fillSSHConfigHostname(target, ssh_config.Get(target, "HostName")),
			Port:   ssh_config.Get(target, "Port"),
		}
	}
}

func fillSSHConfigHostname(target string, configHostName string) string {
	// Check if "%h" is in the config name
	if strings.Contains(configHostName, "%h") {
		// Replace %h with the target
		return strings.ReplaceAll(configHostName, "%h", target)
	}

	// Nothing to replace
	return configHostName
}

// Note: this is necessary for generating proxy configs
func makeConfigToDefaultConfig(makeConfig *easyssh.MakeConfig) *easyssh.DefaultConfig {
	return &easyssh.DefaultConfig{
		User:     makeConfig.User,
		Server:   makeConfig.Server,
		Port:     makeConfig.Port,
		Password: makeConfig.Password,
		KeyPath:  makeConfig.KeyPath,
		Timeout:  makeConfig.Timeout,
	}
}
