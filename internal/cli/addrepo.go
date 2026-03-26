package cli

import (
	"fmt"

	"github.com/thasso/wsm/internal/config"
	"github.com/thasso/wsm/internal/ui"
	"github.com/spf13/cobra"
)

var (
	addName        string
	addRepo        string
	addBranch      string
	addJiraKey     string
	addDescription string
)

var addRepoCmd = &cobra.Command{
	Use:   "add-repo",
	Short: "Add a repository to the workspace manifest",
	Long: `Add a new repository entry to workspace.json. The --name and --repo flags
are required. After adding, run "wsm setup" to clone the new repo.

Example:
  wsm add-repo --name prestoplay-ios --repo prestoplay-ios --branch main --jira-key IOS`,
	RunE: runAddRepo,
}

func init() {
	addRepoCmd.Flags().StringVar(&addName, "name", "", "local directory name (required)")
	addRepoCmd.Flags().StringVar(&addRepo, "repo", "", "GitHub repository name under the org (required)")
	addRepoCmd.Flags().StringVar(&addBranch, "branch", "main", "default branch")
	addRepoCmd.Flags().StringVar(&addJiraKey, "jira-key", "", "Jira project key")
	addRepoCmd.Flags().StringVar(&addDescription, "description", "", "brief description")
	addRepoCmd.MarkFlagRequired("name")
	addRepoCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(addRepoCmd)
}

func runAddRepo(cmd *cobra.Command, args []string) error {
	ws, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	r := config.Repo{
		Name:        addName,
		RepoName:    addRepo,
		Branch:      addBranch,
		JiraKey:     addJiraKey,
		Description: addDescription,
	}

	if err := ws.AddRepo(r); err != nil {
		return err
	}

	if err := ws.Save(cfgFile); err != nil {
		return err
	}

	ui.Success(fmt.Sprintf("Added %s (%s/%s, branch: %s)", r.Name, ws.Org, r.RepoName, r.Branch))
	fmt.Println("Run 'wsm setup' to clone the new repo.")
	return nil
}
