package cli

import (
	"fmt"
	"os"

	"github.com/thasso/wsm/internal/config"
	"github.com/thasso/wsm/internal/gitops"
	"github.com/thasso/wsm/internal/ui"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull latest changes for all repos on their current branch",
	Long: `Run "git pull" in each cloned repository, then update submodules.
Repos that have not been cloned yet are skipped with a warning.`,
	RunE: runPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) error {
	ws, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	fmt.Println(ui.Bold("Pulling latest changes..."))
	fmt.Println()

	for _, repo := range ws.Repos {
		if _, err := os.Stat(repo.Name + "/.git"); os.IsNotExist(err) {
			ui.Fail(fmt.Sprintf("%s — not cloned (run 'wsm setup' first)", repo.Name))
			continue
		}

		ui.Spin(fmt.Sprintf("%s ...", repo.Name))

		updated, err := gitops.Pull(repo.Name)
		if err != nil {
			ui.Fail(fmt.Sprintf("%s: %v", repo.Name, err))
			continue
		}

		if updated {
			ui.Success(fmt.Sprintf("%s updated", repo.Name))
		} else {
			ui.Success(fmt.Sprintf("%s up to date", repo.Name))
		}

		if err := gitops.UpdateSubmodules(repo.Name); err != nil {
			ui.Fail(fmt.Sprintf("Submodule update failed for %s: %v", repo.Name, err))
		}
	}

	fmt.Println()
	return nil
}
