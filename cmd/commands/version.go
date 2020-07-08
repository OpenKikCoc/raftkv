package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/OpenKikCoc/raftkv/version"
)

var (
	// VersionCmd the version command
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version Command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s-%s-%s, CommitHash: %s\n", version.Maj, version.Min, version.Fix, version.GitCommit)
		},
	}
)
