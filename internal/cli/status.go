package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thasso/wsm/internal/config"
	"github.com/thasso/wsm/internal/gitops"
	"github.com/thasso/wsm/internal/ui"
)

var statusJSON bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of all repos (branch, dirty state, ahead/behind)",
	Long: `Display a table showing the current branch, clean/dirty state, and
ahead/behind counts for each repository in the workspace. Repos that
have not been cloned yet are marked accordingly.

Use --json for machine-readable output.`,
	RunE: runStatus,
}

func init() {
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "output status as JSON")
	rootCmd.AddCommand(statusCmd)
}

type statusEntry struct {
	Name          string `json:"name"`
	CurrentBranch string `json:"current_branch"`
	DefaultBranch string `json:"default_branch"`
	IsCloned      bool   `json:"is_cloned"`
	IsDirty       bool   `json:"is_dirty"`
	Ahead         int    `json:"ahead"`
	Behind        int    `json:"behind"`
	HasTracking   bool   `json:"has_tracking"`
}

func runStatus(cmd *cobra.Command, args []string) error {
	ws, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	var entries []statusEntry

	for _, repo := range ws.Repos {
		if _, err := os.Stat(repo.Name + "/.git"); os.IsNotExist(err) {
			entries = append(entries, statusEntry{
				Name:          repo.Name,
				DefaultBranch: repo.Branch,
				IsCloned:      false,
			})
			continue
		}

		s, err := gitops.Status(repo.Name, repo.Branch)
		if err != nil {
			return fmt.Errorf("status for %s: %w", repo.Name, err)
		}

		entries = append(entries, statusEntry{
			Name:          s.Name,
			CurrentBranch: s.CurrentBranch,
			DefaultBranch: s.DefaultBranch,
			IsCloned:      true,
			IsDirty:       s.IsDirty,
			Ahead:         s.Ahead,
			Behind:        s.Behind,
			HasTracking:   s.HasTracking,
		})
	}

	if statusJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	// Table output
	fmt.Printf("%s %s %s %s\n",
		ui.Bold(padRight("REPO", 30)),
		ui.Bold(padRight("BRANCH", 20)),
		ui.Bold(padRight("STATE", 12)),
		ui.Bold("REMOTE"),
	)
	fmt.Printf("%s %s %s %s\n",
		padRight("----", 30),
		padRight("------", 20),
		padRight("-----", 12),
		"------",
	)

	for _, e := range entries {
		if !e.IsCloned {
			fmt.Printf("%s %s\n",
				ui.Red(padRight(e.Name, 30)),
				ui.Red(padRight("(not cloned)", 20)),
			)
			continue
		}

		branchDisplay := ui.Yellow(e.CurrentBranch)
		if e.CurrentBranch != e.DefaultBranch {
			branchDisplay = fmt.Sprintf("%s (default: %s)", ui.Yellow(e.CurrentBranch), e.DefaultBranch)
		}

		var state string
		if e.IsDirty {
			state = ui.Red("✗ dirty")
		} else {
			state = ui.Green("✓ clean")
		}

		var remote string
		if e.HasTracking {
			remote = fmt.Sprintf("↑%d ↓%d", e.Ahead, e.Behind)
		} else {
			remote = "(no tracking)"
		}

		fmt.Printf("%-30s %-20s %-12s %s\n", e.Name, branchDisplay, state, remote)
	}
	fmt.Println()
	return nil
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + fmt.Sprintf("%*s", width-len(s), "")
}
