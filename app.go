package goappbase

import (
	"errors"
	"fmt"
	"log"
	"reflect"

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

	AppSettingsFilename string // with .yml extension please
	AppSettings         interface{}

	ServiceUnitData *mttools.ServiceData

	rootCmd *cobra.Command
}

func NewAppBase() *AppBase {
	app := AppBase{}

	app.Version = BuildVersion
	app.Commit = BuildCommit
	app.Time = BuildTime

	//set default values
	app.ExecutableName = "UNSET_ExecutableName"
	app.AppName = "UNSET_AppName"

	app.AppSettingsFilename = ".settings.yml"

	app.ServiceUnitData = &mttools.ServiceData{
		Name: "service_name_not_set",
	}

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
	//default settings object should be set
	if app.AppSettings == nil {
		log.Fatalln("AppSettings should not be empty")
	}

	base_settings_type := reflect.TypeOf((*AppSettingsBase)(nil)).Elem()

	if !mttools.IsStructEmbeds(app.AppSettings, base_settings_type) {
		log.Fatalln("AppSettings should embed " + base_settings_type.Name())
	}

	//default basic settings values
	v := reflect.ValueOf(app.AppSettings).Elem()
	base_settings := v.FieldByName(base_settings_type.Name()).Addr().Interface().(*AppSettingsBase)

	base_settings.Production = false
	base_settings.WebserverHostname = "localhost"
	base_settings.WebserverPort = 15115
	base_settings.ServiceName = app.ExecutableName
	base_settings.ServiceUser = "www-data"
	base_settings.ServiceGroup = "www-data"

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
		app.buildUninstallCmd(),
		app.buildInitCmd(),
	)
}

func (app *AppBase) loadSettings(filename string) error {
	if mttools.IsFileExists(filename) {
		if err := mttools.LoadYamlSettingFromFile(filename, app.AppSettings); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("File not found: %s", filename)
	}

	if app.AppSettings.(AppSettingsBase).Production {
		// require some settings in PRODUCTION

		if app.AppSettings.(AppSettingsBase).BaseUrl == "" {
			return errors.New("base_url required in production")
		}
	} else {
	}

	return nil
}
