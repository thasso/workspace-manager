package cli

import (
	"fmt"
	"os"

	"github.com/thasso/wsm/internal/config"
	"github.com/thasso/wsm/internal/ui"
	"github.com/spf13/cobra"
)

var (
	initOrg     string
	initJiraURL string
	initForce   bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new workspace.json manifest",
	Long: `Initialize a new workspace by creating a workspace.json file. The --org
flag is required and sets the GitHub organization used for clone URLs.

If workspace.json already exists, use --force to overwrite it.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVar(&initOrg, "org", "", "GitHub organization (required)")
	initCmd.Flags().StringVar(&initJiraURL, "jira-url", "", "Jira base URL (optional)")
	initCmd.Flags().BoolVar(&initForce, "force", false, "overwrite existing workspace.json")
	initCmd.MarkFlagRequired("org")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(cfgFile); err == nil && !initForce {
		return fmt.Errorf("%s already exists (use --force to overwrite)", cfgFile)
	}

	ws := &config.Workspace{
		Org:      initOrg,
		JiraURL:  initJiraURL,
		Repos:    []config.Repo{},
		Projects: []config.Project{},
	}

	if err := ws.Save(cfgFile); err != nil {
		return err
	}

	ui.Success(fmt.Sprintf("Created %s (org: %s)", cfgFile, initOrg))
	return nil
}
