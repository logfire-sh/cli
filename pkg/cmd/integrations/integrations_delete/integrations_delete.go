package integrations_delete

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/helpers"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type DeleteIntegrationOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive   bool
	TeamId        string
	IntegrationId string
}

func NewDeleteIntegrationCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteIntegrationOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete integration",
		Long:  "delete integration",
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire integrations delete

			# start argument setup
			$ logfire integrations delete --team-name <team-name> --integration-id <integration-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			DeleteIntegrationRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-name", "t", "", "Team name from which integrationId is to be deleted.")
	cmd.Flags().StringVarP(&opts.IntegrationId, "integration-id", "i", "", "Integration to be deleted.")
	return cmd
}

func DeleteIntegrationRun(opts *DeleteIntegrationOptions) {
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

	if opts.Interactive {
		if opts.TeamId == "" && opts.IntegrationId == "" {
			opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

			opts.IntegrationId, _ = pre_defined_prompters.AskIntegrationId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)
		}
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId
		}

		if opts.IntegrationId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Integration id is required.\n", cs.FailureIcon())
			os.Exit(0)
		}
	}

	err = APICalls.DeleteIntegration(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId,
		opts.IntegrationId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Integration deleted successfully!\n", cs.SuccessIcon())
	}
}
