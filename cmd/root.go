package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var (
	dryrun          = false
	remote          = ""
	sshFile         = ""
	sshFilePassword = ""
)

func init() {
	usr, err := user.Current()
	handleError(err)

	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "avoid to touch current git repository")
	rootCmd.PersistentFlags().StringVar(&remote, "remote", "origin", "origin where push to")
	rootCmd.PersistentFlags().StringVar(&sshFile, "sshfile", usr.HomeDir+"/.ssh/id_rsa", "ssh file used to authenticate on git remote")
	rootCmd.PersistentFlags().StringVar(&sshFilePassword, "sshpwd", "", "ssh file password")

	rootCmd.AddCommand(makeBumpCmd("patch"))
	rootCmd.AddCommand(makeBumpCmd("minor"))
	rootCmd.AddCommand(makeBumpCmd("major"))
}

var rootCmd = &cobra.Command{
	Use:   "gat",
	Short: "Gat is a tagging tool for git.",
	Long:  `Gat is a tagging tool for git. It gets last version and tags it for you automatically`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	if err == nil {
		return
	}
	if err == git.NoErrAlreadyUpToDate {
		fmt.Println("origin remote was up to date, no push done")
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
