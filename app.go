package goappbase

import (
	"log"

	"github.com/mitoteam/mttools"
	"github.com/spf13/cobra"
)

const DEV_MODE_LABEL = "DEV"

// Variables to be set by compiler
var (
	BuildVersion = DEV_MODE_LABEL
	BuildCommit  = DEV_MODE_LABEL
	BuildTime    = DEV_MODE_LABEL
)

type AppBase struct {
	ExecutableName  string //executable command name
	AppName         string //Long name
	LongDescription string //Long description

	Version string //Version (auto set by compiler)
	Commit  string //Git commit hash
	Time    string //Build time

	ServiceUnitData *mttools.ServiceData

	rootCmd *cobra.Command
}

func NewAppBase() *AppBase {
	app := AppBase{}

	app.Version = BuildVersion
	app.Commit = BuildCommit
	app.Time = BuildTime

	app.ExecutableName = "UNSET_ExecutableName"
	app.AppName = "UNSET_AppName"

	app.ServiceUnitData = &mttools.ServiceData{}

	app.buildRootCmd()

	return &app
}

func (app *AppBase) Run() {
	app.internalInit()

	//cli application - we just let cobra to do its job
	if err := app.rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func (app *AppBase) internalInit() {
	//setup root cmd
	app.rootCmd.Use = app.ExecutableName
	app.rootCmd.Long = app.AppName

	if app.LongDescription != "" {
		app.rootCmd.Long += " - " + app.LongDescription

	}

	//add built-in commands
	app.rootCmd.AddCommand(
		app.buildVersionCmd(),
		app.buildInstallCmd(),
	)
}

func (app *AppBase) GetRootCmd() *cobra.Command {
	return app.rootCmd
}
