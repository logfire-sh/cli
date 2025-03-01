package views_list

import (
	"fmt"
	"github.com/logfire-sh/cli/pkg/cmdutil/helpers"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type ViewListOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
}

func NewViewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ViewListOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all views",
		Long:  "list all views",
		Args:  cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire views list

			# start argument setup
			$ logfire views list --team-name <team-name>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			viewsListRun(opts)
		},
	}
	cmd.Flags().StringVar(&opts.TeamId, "team-name", "", "Team name to be deleted.")
	return cmd
}

func viewsListRun(opts *ViewListOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	client := http.Client{}

	if opts.TeamId != "" {
		teamId := helpers.TeamNameToTeamId(&client, cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

		if teamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s no team with name: %s found.\n", cs.FailureIcon(), opts.TeamId)
			return
		}

		opts.TeamId = teamId
	}

	if opts.Interactive && opts.TeamId == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId

			println("opts.TeamId: ", opts.TeamId)
		}
	}

	list, err := APICalls.ListView(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to list view\n", cs.FailureIcon())
	} else if len(list) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "View-Id"})

		for _, i2 := range list {
			table.Append([]string{i2.Name, i2.Id})
		}

		table.Render()
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s No views created. Please create a view\n", cs.FailureIcon())
		os.Exit(0)
	}
}
