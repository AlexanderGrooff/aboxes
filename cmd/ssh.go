package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/appleboy/easyssh-proxy"
	"github.com/kevinburke/ssh_config"
	log "github.com/sirupsen/logrus"
)

func executeCommands(targets []string, commands []string, outputFile string) {
	// Open output file in append mode
	var file *os.File
	var err error
	if outputFile != "" {
		file, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open output file: %s", err)
		}
		defer file.Close()
	}

	// Run commands in parallel for hosts
	var wg sync.WaitGroup
	for _, target := range targets {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			ssh := getConfigForHost(target)
			// TODO: Use same SSH connection for all commands
			for _, command := range commands {
				stdout, stderr, _, err := ssh.Run(command)
				if err != nil {
					log.Printf("Error running command on %s: %s\n", target, err)
					continue
				}
				// Identify host with output
				log.Printf("%s: %s\n", target, stdout)
				if stderr != "" {
					log.Printf("%s: %s\n", target, stderr)
				}
				// Write output to file if given
				if file != nil {
					if _, err := file.WriteString(fmt.Sprintf("%s: %s\n%s\n", target, command, stdout)); err != nil {
						log.Printf("Error writing to file: %s\n", err)
					}
					if stderr != "" {
						if _, err := file.WriteString(fmt.Sprintf("%s: %s\n%s\n", target, command, stderr)); err != nil {
							log.Printf("Error writing to file: %s\n", err)
						}
					}
				}
			}
		}(target)
	}
	wg.Wait()
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
