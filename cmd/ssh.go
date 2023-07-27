package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/appleboy/easyssh-proxy"
	log "github.com/sirupsen/logrus"
)

func executeCommands(targets []string, commands []string, outputFile string, format string) {
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
				result := RunAndParse(target, command, ssh)
				if err != nil {
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
		}(target)
	}
	wg.Wait()
}

func RunAndParse(target string, command string, ssh *easyssh.MakeConfig) Result {
	stdout, stderr, _, err := ssh.Run(command)
	return Result{
		Target:   target,
		Hostname: ssh.Server,
		Stdout:   stdout,
		Stderr:   stderr,
		Error:    err,
	}
}
