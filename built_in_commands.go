package goappbase

import (
	"errors"
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

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			//Load Settings
			if mttools.IsFileExists(app.AppSettingsFilename) {
				if err := app.loadSettings(); err != nil {
					return err
				}
			} else {
				if cmd.Name() != "init" && cmd.Name() != "version" {
					log.Fatalln(
						"No "+app.AppSettingsFilename+" file found. Please create one or use `%s init` command.", app.ExecutableName,
					)
				}
			}

			return nil
		},
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
				unitData := &mttools.ServiceData{
					Name:      app.baseSettings.ServiceName,
					User:      app.baseSettings.ServiceUser,
					Group:     app.baseSettings.ServiceGroup,
					Autostart: app.baseSettings.ServiceAutostart,
				}

				if err := unitData.InstallSystemdService(); err != nil {
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

	return cmd
}

func (app *AppBase) buildUninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove installed system service " + app.AppName,

		Run: func(cmd *cobra.Command, args []string) {
			if mttools.IsSystemdAvailable() {
				unitData := &mttools.ServiceData{
					Name: app.baseSettings.ServiceName,
				}

				if err := unitData.UninstallSystemdService(); err != nil {
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

	return cmd
}

func (app *AppBase) buildInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates settings file with defaults in working directory.",

		RunE: func(cmd *cobra.Command, args []string) error {
			if mttools.IsFileExists(app.AppSettingsFilename) {
				return errors.New("Can not initialize existing file: " + app.AppSettingsFilename)
			}

			comment := `File was created automatically by '` + app.AppName + ` init' command. There are all
available options listed here with its default values. Recommendation is to edit options you
want to change and remove all others with default values to keep this as simple as possible.
`

			if err := app.saveSettings(comment); err != nil {
				return err
			}

			fmt.Println("Default app settings written to " + app.AppSettingsFilename)

			return nil
		},
	}

	return cmd
}

func (app *AppBase) buildInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Prints info about app, settings, status etc.",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", app.AppName)
			fmt.Print("================================\n")
			fmt.Printf("Version: %s\n", app.Version)

			// Settings
			fmt.Print("\n================================\n")
			fmt.Print("SETTINGS\n")
			fmt.Print("================================\n")
			app.printSettings()
		},
	}

	return cmd
}
