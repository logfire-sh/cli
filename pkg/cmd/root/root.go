package root

import (
	"fmt"
	"os"

	"github.com/logfire-sh/cli/pkg/cmd/delete_profile"
	"github.com/logfire-sh/cli/pkg/cmd/set_password"

	"github.com/logfire-sh/cli/pkg/cmd/settings"

	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/alerts"
	"github.com/logfire-sh/cli/pkg/cmd/bootstrap"
	"github.com/logfire-sh/cli/pkg/cmd/check_endpoint"
	"github.com/logfire-sh/cli/pkg/cmd/integrations"
	"github.com/logfire-sh/cli/pkg/cmd/reset_password"
	"github.com/logfire-sh/cli/pkg/cmd/roundtrip"
	"github.com/logfire-sh/cli/pkg/cmd/sql"
	"github.com/logfire-sh/cli/pkg/cmd/update_profile"
	"github.com/logfire-sh/cli/pkg/cmd/views"
	"github.com/logfire-sh/cli/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/pkg/cmd/login"
	"github.com/logfire-sh/cli/pkg/cmd/logout"
	"github.com/logfire-sh/cli/pkg/cmd/signup"
	"github.com/logfire-sh/cli/pkg/cmd/sources"
	"github.com/logfire-sh/cli/pkg/cmd/stream"
	"github.com/logfire-sh/cli/pkg/cmd/tail"
	"github.com/logfire-sh/cli/pkg/cmd/teams"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type PromptRootOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Interactive       bool
	Choice            string
	NotLoggedInChoice string
	LoggedIn          bool
}

var choices = []string{"Reset password", "Logout", "Sources", "Teams",
	"Tail", "Start Stream", "Views", "Alerts", "Integrations", "SQL", "Update profile", "Settings", "Round trip", "Delete profile", "Exit"}

var NotLoggedInChoices = []string{"Signup", "Login"}

func NewCmdRoot(f *cmdutil.Factory, cmdCh chan bool) (*cobra.Command, error) {
	opts := &PromptRootOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
	}

	cfg, err := f.Config()
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %s", err)
	}

	cmd := &cobra.Command{
		Use: "logfire <command> <subcommand> [flags]",
		//Args:  cobra.ExactArgs(1),
		Short: "Logfire CLI",
		Long:  `Work seamlessly with logfire log management system from the command line.`,
		Example: heredoc.Doc(`
			$ logfire login
			$ logfire stream livetail

		`),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// require that the user is authenticated before running most commands
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if cmdutil.IsAuthCheckEnabled(cmd) && !cmdutil.CheckAuth(cfg) && cmd.Name() != "bootstrap" {
				if opts.Interactive {
					fmt.Printf("You are not logged in\nTo get started with Logfire CLI, please choose below to Login or Signup:\n")

					NotLoggedInPromptRun(opts)

					switch opts.NotLoggedInChoice {
					case NotLoggedInChoices[0]:
						signup.NewSignupCmd(f).Run(cmd, []string{})
					case NotLoggedInChoices[1]:
						login.NewLoginCmd(f).Run(cmd, []string{})
					default:
						break
					}
				} else {
					fmt.Fprint(opts.IO.ErrOut, authHelp())
				}
				opts.LoggedIn = false
				return nil
			}

			opts.LoggedIn = true

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if opts.LoggedIn {
				PromptRootRun(opts)

				switch opts.Choice {
				case choices[0]:
					reset_password.NewResetPasswordCmd(f).Run(cmd, []string{})
				case choices[1]:
					logout.NewLogoutCmd(f).Run(cmd, []string{})
				case choices[2]:
					sources.NewCmdSource(f).Run(cmd, []string{})
				case choices[3]:
					teams.NewCmdTeam(f).Run(cmd, []string{})
				case choices[4]:
					tail.NewTailCmd(f).Run(cmd, []string{})
				case choices[5]:
					stream.NewCmdStream(f).Run(cmd, []string{})
				case choices[6]:
					views.NewCmdViews(f).Run(cmd, []string{})
				case choices[7]:
					alerts.NewCmdAlerts(f).Run(cmd, []string{})
				case choices[8]:
					integrations.NewCmdIntegrations(f).Run(cmd, []string{})
				case choices[9]:
					sql.NewCmdSql(f).Run(cmd, []string{})
				case choices[10]:
					update_profile.UpdateProfileCmd(f).Run(cmd, []string{})
				case choices[11]:
					settings.SettingsCmd(f).Run(cmd, []string{})
				case choices[12]:
					roundtrip.NewCmdRoundTrip(f).Run(cmd, []string{})
				case choices[13]:
					delete_profile.DeleteProfileCmd(f).Run(cmd, []string{})

				case "Exit":
					os.Exit(0)
				default:
					break
				}
			}
		},
	}

	cmd.PersistentFlags().Bool("help", false, "Show help for command")

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		rootHelpFunc(f, c, args)
	})

	cmd.SetUsageFunc(func(c *cobra.Command) error {
		return rootUsageFunc(f.IOStreams.ErrOut, c)
	})

	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	cmd.AddGroup(&cobra.Group{
		ID:    "core",
		Title: "Core commands",
	})

	cmd.AddCommand(signup.NewSignupCmd(f))
	cmd.AddCommand(login.NewLoginCmd(f))
	cmd.AddCommand(reset_password.NewResetPasswordCmd(f))
	cmd.AddCommand(set_password.NewSetPasswordCmd(f))
	cmd.AddCommand(logout.NewLogoutCmd(f))
	cmd.AddCommand(sources.NewCmdSource(f))
	cmd.AddCommand(teams.NewCmdTeam(f))
	cmd.AddCommand(tail.NewTailCmd(f))
	cmd.AddCommand(stream.NewCmdStream(f))
	cmd.AddCommand(views.NewCmdViews(f))
	cmd.AddCommand(alerts.NewCmdAlerts(f))
	cmd.AddCommand(integrations.NewCmdIntegrations(f))
	cmd.AddCommand(sql.NewCmdSql(f))
	cmd.AddCommand(check_endpoint.NewCheckEndpointCmd(f))
	cmd.AddCommand(update_profile.UpdateProfileCmd(f))
	cmd.AddCommand(bootstrap.NewCmdBootstrap(f))
	cmd.AddCommand(roundtrip.NewCmdRoundTrip(f))
	cmd.AddCommand(settings.SettingsCmd(f))
	cmd.AddCommand(delete_profile.DeleteProfileCmd(f))

	// go func() {
	// 	for range cmdCh {
	// 		PromptRootRun(opts)
	// 	}
	// }()

	return cmd, nil
}

func NotLoggedInPromptRun(opts *PromptRootOptions) {
	if opts.Interactive {
		var err error
		opts.NotLoggedInChoice, err = opts.Prompter.Select("What do you want to do?", "", NotLoggedInChoices)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "Failed to read choice\n")
		}
	}
}

func PromptRootRun(opts *PromptRootOptions) {
	cs := opts.IO.ColorScheme()
	if !opts.Interactive {
		return
	}

	if opts.Interactive {
		var err error
		opts.Choice, err = opts.Prompter.Select("What do you want to do?", "", choices)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read choice\n", cs.FailureIcon())
			return
		}
	}
}
