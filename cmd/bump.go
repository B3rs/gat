package cmd

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func makeBumpCmd(bumpType string) *cobra.Command {
	return &cobra.Command{
		Use:   bumpType,
		Short: "Bumps the " + bumpType + " version of your software",
		Long:  "Bumps the " + bumpType + " version of your software",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("tag " + bumpType + " version")

			repo, err := git.PlainOpen(".")
			handleError(err)

			commit, latestTag, err := getLatestTagFromGit(repo)
			handleError(err)

			head, err := isHead(repo, commit)
			handleError(err)
			if head {
				fmt.Println("nothing to tag: latest commit is already tagged.")
			}

			latestVersion := latestTag.Name().Short()
			newVersion := bumpVersion(latestVersion, bumpType)

			fmt.Printf("%s => %s\n", latestVersion, newVersion)

			if !dryrun {
				fmt.Printf("tagging repo version %s...\n", newVersion)
				ref, err := tag(repo, newVersion)
				handleError(err)

				fmt.Printf("pushing %s to %s...\n", ref.Name(), remote)
				err = push(repo, ref, remote, sshFile)
				handleError(err)

				fmt.Println("done")
			}
		},
	}
}
