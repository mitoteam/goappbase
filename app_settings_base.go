package goappbase

type AppSettingsBase struct {
	Production bool `yaml:"production" yaml_comment:"Production mode"`

	BaseUrl string `yaml:"base_url" yaml_comment:"Base external site URL (with protocol and port, no trailing slash)"`

	WebserverHostname string `yaml:"webserver_hostname" yaml_comment:"Webserver hostname"`
	WebserverPort     uint16 `yaml:"webserver_port" yaml_comment:"Webserver port number"`

	ServiceName      string `yaml:"service_name" yaml_comment:"Service name for 'install' command"`
	ServiceUser      string `yaml:"service_user" yaml_comment:"User for 'install' command"`
	ServiceGroup     string `yaml:"service_group" yaml_comment:"Group for 'install' command"`
	ServiceAutostart bool   `yaml:"service_autostart" yaml_comment:"Set autostart while installing service"`
}
