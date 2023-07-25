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
					log.Warn("Error running command on ", target, ": ", err)
					continue
				}
				// Identify host with output
				log.Info(target, ": ", stdout)
				if stderr != "" {
					log.Info(target, ": ", stderr)
				}
				// Write output to file if given
				if file != nil {
					if _, err := file.WriteString(fmt.Sprintf("%s: %s\n%s\n", target, command, stdout)); err != nil {
						log.Warn("Error writing to file: ", err)
					}
					if stderr != "" {
						if _, err := file.WriteString(fmt.Sprintf("%s: %s\n%s\n", target, command, stderr)); err != nil {
							log.Warn("Error writing to file: ", err)
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
	user := ssh_config.Get(target, "User")
	if user == "" {
		user = os.Getenv("USER")
	}
	port := ssh_config.Get(target, "Port")
	if port == "" {
		port = "22"
	}
	hostname := fillSSHConfigHostname(target, ssh_config.Get(target, "HostName"))
	if hostname == "" {
		hostname = target
	}

	log.Debug("Creating SSH config for target ", target, " with hostname ", hostname, " and proxy ", proxy)

	if proxy != "" {
		log.Debug("Using proxy ", proxy, " for target ", target)
		return &easyssh.MakeConfig{
			User:   user,
			Server: hostname,
			Port:   port,
			Proxy:  *makeConfigToDefaultConfig(getConfigForHost(proxy)),
		}
	} else {
		log.Debug("No proxy found for target ", target)
		return &easyssh.MakeConfig{
			User:   user,
			Server: hostname,
			Port:   port,
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
