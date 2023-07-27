package cmd

import (
	"os"
	"strings"

	"github.com/appleboy/easyssh-proxy"
	"github.com/kevinburke/ssh_config"
	log "github.com/sirupsen/logrus"
)

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
