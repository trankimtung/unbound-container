package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	anchorPath   = "/opt/unbound/etc/unbound/var/root.key"
	confPath     = "/opt/unbound/etc/unbound/unbound.conf"
	templatePath = "/opt/unbound/etc/unbound/unbound.conf.template"
)

func main() {
	port := os.Getenv("UNBOUND_PORT")
	if port == "" {
		port = "5335"
	}

	// Initialize or refresh the DNSSEC root trust anchor (RFC 5011).
	cmd := exec.Command("unbound-anchor", "-a", anchorPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok || exitErr.ExitCode() > 1 {
			fmt.Fprintf(os.Stderr, "failed to run unbound-anchor. dnssec may not work correctly: %v\n", err)
		}
	}

	// Render config from template substituting ${UNBOUND_PORT} only if it doesn't exist.
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		template, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read config template %s: %v\n", templatePath, err)
			os.Exit(1)
		}
		rendered := strings.ReplaceAll(string(template), "${UNBOUND_PORT}", port)
		if err := os.WriteFile(confPath, []byte(rendered), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write config %s: %v\n", confPath, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Starting unbound on port %s\n", port)

	// Exec into unbound, replacing this process so unbound becomes PID 1.
	unboundBin, err := exec.LookPath("unbound")
	if err != nil {
		fmt.Fprintf(os.Stderr, "unbound not found in PATH: %v\n", err)
		os.Exit(1)
	}
	if err := syscall.Exec(unboundBin, []string{"unbound", "-d", "-c", confPath}, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "exec unbound: %v\n", err)
		os.Exit(1)
	}
}
