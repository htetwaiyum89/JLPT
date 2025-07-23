package main

import (
        "crypto/x509/pkix"
        "encoding/base64"
        "errors"
        "flag"
        "fmt"
        "io/ioutil"
        "os"
        "strings"
        "time"

        "github.com/masterzen/winrm"
        "github.com/mattn/go-isatty"
)
func main() {
	var hostname, username, password, command string
	var port int
	var useHTTPS, insecure, showHelp bool

	// Define CLI flags
	flag.StringVar(&hostname, "hostname", "", "Target Windows host (IP or hostname)")
	flag.StringVar(&username, "username", "", "Username for authentication")
	flag.StringVar(&password, "password", "", "Password for authentication")
	flag.IntVar(&port, "port", 5985, "WinRM port (5985 for HTTP, 5986 for HTTPS)")
	flag.BoolVar(&useHTTPS, "https", false, "Use HTTPS instead of HTTP")
	flag.BoolVar(&insecure, "insecure", false, "Skip SSL verification")
	flag.StringVar(&command, "cmd", "", "Command to run on the remote Windows host")
	flag.BoolVar(&showHelp, "h", false, "Show help")

	flag.Parse()

	if showHelp || hostname == "" || username == "" || password == "" || (command == "" && flag.NArg() == 0) {
		fmt.Printf("Usage:\n")
		fmt.Printf("  %s -hostname <IP> -username <user> -password <pass> -cmd \"command\"\n", os.Args[0])
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Fallback to positional args if -cmd not used
	if command == "" {
		command = flag.Args()[0]
	}

	// Optional args to command
	cmdArgs := flag.Args()[1:]

	endpoint := winrm.NewEndpoint(hostname, port, useHTTPS, insecure, nil, nil, nil, 0)
	params := winrm.DefaultParameters
	client, err := winrm.NewClientWithParameters(endpoint, username, password, params)
	if err != nil {
		log.Fatalf("Failed to create WinRM client: %s", err)
	}

	fullCommand := command
	if len(cmdArgs) > 0 {
		fullCommand = fullCommand + " " + strings.Join(cmdArgs, " ")
	}

	exitCode, err := client.Run(fullCommand, os.Stdout, os.Stderr)
	if err != nil {
		log.Fatalf("Command execution failed: %s", err)
	}
	os.Exit(exitCode)
}

