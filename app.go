package goappbase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	Version     string //Version (auto set by compiler)
	BuildCommit string //Git commit hash
	BuildTime   string //Build time
	BuildWith   string //Build information

	AppSettingsFilename string           // with .yml extension please
	AppSettings         interface{}      //pointer to struct embedding AppSettingsBase
	baseSettings        *AppSettingsBase //pointer to *AppSettingsBase, set in internalInit()

	serviceAutostart bool

	rootCmd *cobra.Command

	//contexts and timeout settings
	BaseContext     context.Context
	ShutdownTimeout time.Duration

	//web router
	webRouter           *gin.Engine
	WebRouterLogQueries bool                // true = extended query logging (--query-log option of `run`)
	BuildWebRouterF     func(r *gin.Engine) // function to build web router for `run` command

	//callbacks (aka event handlers)
	PreRunF  func() error // called before starting `run` command. Stops executions inf error returned.
	PostRunF func() error // called after starting `run` command. Stops executions inf error returned.
}

// Initializes new application.
// settings - application settings default values. Pointer to struct that embeds AppSettingsBase.
func NewAppBase(settings interface{}) *AppBase {
	app := AppBase{}

	//default settings values
	app.AppSettingsFilename = ".settings.yml"
	if settings == nil {
		log.Fatalln("settings should not be empty")
	}

	base_settings_type := reflect.TypeOf((*AppSettingsBase)(nil)).Elem()

	if !mttools.IsStructEmbeds(settings, base_settings_type) {
		log.Fatalln("settings structure should embed " + base_settings_type.Name())
	}

	app.AppSettings = settings

	v := reflect.ValueOf(app.AppSettings).Elem()
	app.baseSettings = v.FieldByName(base_settings_type.Name()).Addr().Interface().(*AppSettingsBase)

	app.baseSettings.checkDefaultValues(&AppSettingsBase{
		WebserverHostname: "localhost",
		WebserverPort:     15115,
		ServiceName:       app.ExecutableName,
		ServiceUser:       "www-data",
		ServiceGroup:      "www-data",
	})

	//global application base context
	app.BaseContext = context.Background()

	//compilation data
	app.Version = BuildVersion
	app.BuildCommit = BuildCommit
	app.BuildTime = BuildTime
	app.BuildWith = runtime.Version()

	//set default values
	app.ExecutableName = "UNSET_ExecutableName"
	app.AppName = "UNSET_AppName"

	app.ShutdownTimeout = 10 * time.Second

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
		app.buildUninstallCmd(),
		app.buildInitCmd(),
		app.buildInfoCmd(),
		app.buildRunCmd(),
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

	// Settings post-processing

	if app.baseSettings.Production {
		// require some settings in PRODUCTION
		if app.baseSettings.BaseUrl == "" {
			return errors.New("base_url required in production")
		}

		if app.baseSettings.WebserverCookieSecret == "" {
			return errors.New("webserver_cookie_secret required in production")
		}
	} else {
		// or use pre-defined values in DEV
		if app.baseSettings.BaseUrl == "" {
			app.baseSettings.BaseUrl = "http://" + app.baseSettings.WebserverHostname +
				":" + strconv.Itoa(int(app.baseSettings.WebserverPort))
		}

		if app.baseSettings.WebserverCookieSecret == "" {
			app.baseSettings.WebserverCookieSecret = "DEFAULT_DEV_SECRET"
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
