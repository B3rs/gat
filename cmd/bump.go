package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func makeBumpCmd(bump string) *cobra.Command {
	return &cobra.Command{
		Use:   bump,
		Short: "Bumps the " + bump + " version of your software",
		Long:  "Bumps the " + bump + " version of your software",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("tag " + bump + " version")

			latestVersion, err := getLatestVersionFromGit(".")
			handleError(err)

			newVersion := bumpVersion(latestVersion, bump)
			printAction(latestVersion, newVersion)

			if !dryrun {
				tagAndPush(".", remote, sshFile, newVersion)
			}
		},
	}
}
