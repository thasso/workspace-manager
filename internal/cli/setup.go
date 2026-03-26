package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thasso/wsm/internal/config"
	"github.com/thasso/wsm/internal/gitignore"
	"github.com/thasso/wsm/internal/gitops"
	"github.com/thasso/wsm/internal/ui"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Clone missing repos and update .gitignore",
	Long: `Clone all repositories listed in workspace.json that are not yet present
in the workspace directory. After cloning, git submodules are initialized
recursively. Finally, the .gitignore file is updated with a managed section
listing all repo directories.`,
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) error {
	ws, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	fmt.Println(ui.Bold("Setting up workspace..."))
	fmt.Println()

	for _, repo := range ws.Repos {
		if _, err := os.Stat(repo.Name); err == nil {
			ui.Success(fmt.Sprintf("%s (already cloned)", repo.Name))
			continue
		}

		url := repo.CloneURL(ws.Org)
		ui.Spin(fmt.Sprintf("Cloning %s from %s (branch: %s)...", repo.Name, url, repo.Branch))

		if err := gitops.Clone(url, repo.Branch, repo.Name); err != nil {
			ui.Fail(fmt.Sprintf("Failed to clone %s: %v", repo.Name, err))
			continue
		}
		ui.Success(fmt.Sprintf("%s cloned", repo.Name))

		ui.Spin(fmt.Sprintf("Initializing submodules for %s...", repo.Name))
		if err := gitops.InitSubmodules(repo.Name); err != nil {
			ui.Fail(fmt.Sprintf("Submodule init failed for %s: %v", repo.Name, err))
		}
	}

	fmt.Println()
	fmt.Println("Updating .gitignore...")
	if err := gitignore.Update(".gitignore", ws.RepoNames()); err != nil {
		return fmt.Errorf("updating .gitignore: %w", err)
	}
	ui.Success(".gitignore updated")

	fmt.Println()
	fmt.Println(ui.Green(ui.Bold("Workspace setup complete.")))
	return nil
}
