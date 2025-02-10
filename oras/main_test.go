package main

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

// Block all incoming network traffic
func blockAllIncomingNetwork() error {
	cmd := exec.Command("sudo", "iptables", "-I", "INPUT", "1", "-j", "DROP")
	return cmd.Run()
}

// Unblock all incoming network traffic
func unblockAllIncomingNetwork() error {
	cmd := exec.Command("sudo", "iptables", "-D", "INPUT", "-j", "DROP")
	return cmd.Run()
}

// This test must be run as sudo
func TestOrasNetworkInterruption(t *testing.T) {
	t.SkipNow()
	// Start a goroutine to block and then unblock the network
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Blocking incoming network traffic...")
		if err := blockAllIncomingNetwork(); err != nil {
			t.Errorf("failed to block incoming network: %v", err)
			return
		}

		time.Sleep(5 * time.Second) // Keep the network blocked for 5 seconds
		fmt.Println("Unblocking incoming network traffic...")
		if err := unblockAllIncomingNetwork(); err != nil {
			t.Errorf("failed to unblock incoming network: %v", err)
		}
	}()

	// Start the ORAS pull operation
	if err := doOras(); err != nil {
		t.Fatalf("doOras failed: %v", err)
	}
}
