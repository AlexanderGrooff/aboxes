package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/appleboy/easyssh-proxy"
	log "github.com/sirupsen/logrus"
)

func executeCommands(targets []string, commands []string, outputFile string, format string, files []string) {
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
				result := runAndParse(target, command, ssh)
				if result.Error != nil {
					log.Warn("Error running command on ", target, ": ", err)
					continue
				}
				output := result.toString(format)
				log.Info(output)
				// Write output to file if given
				if file != nil {
					if _, err := file.WriteString(fmt.Sprintf("%s\n", output)); err != nil {
						log.Warn("Error writing to file: ", err)
					}
				}
			}
			for _, script := range files {
				log.Debug("Running script ", script, " on ", target)
				result := runScriptOnHost(target, script, ssh)
				if result.Error != nil {
					log.Warn("Error running script on ", target, ": ", err)
					continue
				}
				output := result.toString(format)
				log.Info(output)
				// Write output to file if given
				if file != nil {
					if _, err := file.WriteString(fmt.Sprintf("%s\n", output)); err != nil {
						log.Warn("Error writing to file: ", err)
					}
				}
			}
		}(target)
	}
	wg.Wait()
}

func runAndParse(target string, command string, ssh *easyssh.MakeConfig) Result {
	stdout, stderr, _, err := ssh.Run(command)
	return Result{
		Target:   target,
		Hostname: ssh.Server,
		Stdout:   stdout,
		Stderr:   stderr,
		Error:    err,
	}
}

func runScriptOnHost(target string, script string, ssh *easyssh.MakeConfig) Result {
	// Only keep the filename without the directory
	basename := script[strings.LastIndex(script, "/")+1:]

	remotePath := fmt.Sprintf("/tmp/%s", basename)
	if err := ssh.Scp(script, remotePath); err != nil {
		log.Warn("Error copying script to host: ", err)
		return Result{
			Target:   target,
			Hostname: ssh.Server,
			Error:    err,
		}
	}
	cmd := fmt.Sprintf("bash -s < %s", remotePath)
	result := runAndParse(target, cmd, ssh)

	// Remove script from host in background
	go func() {
		if _, _, _, err := ssh.Run(fmt.Sprintf("rm %s", remotePath)); err != nil {
			log.Warn("Error removing script from host: ", err)
		}
	}()
	return result
}
