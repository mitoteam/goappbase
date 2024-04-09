package goappbase

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

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

	AppSettingsFilename string           // with .yml extension please
	AppSettings         interface{}      //pointer to struct embedding AppSettingsBase
	baseSettings        *AppSettingsBase //pointer to *AppSettingsBase, set in internalInit()

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

	v := reflect.ValueOf(app.AppSettings).Elem()
	app.baseSettings = v.FieldByName(base_settings_type.Name()).Addr().Interface().(*AppSettingsBase)

	//default basic settings values
	app.baseSettings.Production = false
	app.baseSettings.WebserverHostname = "localhost"
	app.baseSettings.WebserverPort = 15115
	app.baseSettings.ServiceName = app.ExecutableName
	app.baseSettings.ServiceUser = "www-data"
	app.baseSettings.ServiceGroup = "www-data"
	app.baseSettings.ServiceAutostart = true

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
		app.buildInfoCmd(),
	)
}

func (app *AppBase) loadSettings() error {
	if mttools.IsFileExists(app.AppSettingsFilename) {
		if err := mttools.LoadYamlSettingFromFile(app.AppSettingsFilename, app.AppSettings); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("File not found: %s", app.AppSettingsFilename)
	}

	if app.baseSettings.Production {
		// require some settings in PRODUCTION

		if app.baseSettings.BaseUrl == "" {
			return errors.New("base_url required in production")
		}
	} else {
		// or use pre-defined values in DEV
		if app.baseSettings.BaseUrl == "" {
			app.baseSettings.BaseUrl = "http://" + app.baseSettings.WebserverHostname +
				":" + strconv.Itoa(int(app.baseSettings.WebserverPort))
		}
	}

	return nil
}

func (app *AppBase) saveSettings(comment string) error {
	return mttools.SaveYamlSettingToFile(app.AppSettingsFilename, comment, app.AppSettings)
}

func (app *AppBase) printSettings() {
	mttools.PrintYamlSettings(app.AppSettings)
}
