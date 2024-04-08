package goappbase

import (
	"fmt"
	"log"

	"github.com/mitoteam/mttools"
	"github.com/spf13/cobra"
)

func (app *AppBase) buildRootCmd() {
	app.rootCmd = &cobra.Command{
		Version: BuildVersion,

		//disable default 'completion' subcommand
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},

		Run: func(cmd *cobra.Command, args []string) {
			//show help if no subcommand given
			cmd.Help()
		},

		/*PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			///Load Settings
			if mttools.IsFileExists(app.AppSettingsFilename) {
				if err := app.Global.AppSettings.Load(app.AppSettingsFilename); err != nil {
					return err
				}
			} else {
				if cmd.Name() != "init" && cmd.Name() != "version" {
					log.Fatalln(
						"No " + app.AppSettingsFilename + " file found. Please create one or use `twsbot init` command.",
					)
				}
			}

			return nil
		},*/
	}
}

func (app *AppBase) buildVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the raw version number of " + app.AppName + ".",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(app.Version)
		},
	}
}

func (app *AppBase) buildInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Creates system service to run " + app.AppName,

		Run: func(cmd *cobra.Command, args []string) {
			if mttools.IsSystemdAvailable() {
				if err := app.ServiceUnitData.InstallSystemdService(); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatalf(
					"Directory %s does not exists. Only systemd based services supported for now.\n",
					mttools.SystemdServiceDirPath,
				)
			}
		},
	}

	cmd.PersistentFlags().BoolVar(
		&app.ServiceUnitData.Autostart,
		"autostart",
		false,
		"Set service to be auto started after boot. Please note this does not auto starts service after installation.",
	)

	return cmd
}
