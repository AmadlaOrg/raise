package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	sshProvider string

	sshLookPath = exec.LookPath

	// SSHCmd opens an SSH session via a raise plugin.
	SSHCmd = &cobra.Command{
		Use:   "ssh [name]",
		Short: "SSH into infrastructure using a raise plugin",
		Long:  "Opens an interactive SSH session by replacing the current process with the raise-* plugin's ssh command (--provider flag).",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runSSH,
	}
)

func init() {
	SSHCmd.Flags().StringVar(&sshProvider, "provider", "", "Provider plugin name (e.g. libvirt, virtualbox, aws)")
	_ = SSHCmd.MarkFlagRequired("provider")
}

func runSSH(cmd *cobra.Command, args []string) error {
	pluginName := "raise-" + sshProvider

	pluginPath, err := sshLookPath(pluginName)
	if err != nil {
		return fmt.Errorf("plugin %s not found in PATH: %w", pluginName, err)
	}

	execArgs := []string{pluginName, "ssh"}
	if len(args) > 0 {
		execArgs = append(execArgs, args[0])
	}

	// Run the plugin for interactive SSH.
	sshCmd := exec.Command(pluginPath, execArgs[1:]...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr
	return sshCmd.Run()
}
