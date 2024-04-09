package goappbase

type AppSettingsBase struct {
	Production bool `yaml:"production" yaml_comment:"Production mode"`

	BaseUrl string `yaml:"base_url" yaml_comment:"Base external site URL (with protocol and port, no trailing slash)"`

	WebserverHostname     string `yaml:"webserver_hostname" yaml_comment:"Webserver hostname"`
	WebserverPort         uint16 `yaml:"webserver_port" yaml_comment:"Webserver port number"`
	WebserverCookieSecret string `yaml:"webserver_cookie_secret" yaml_comment:"Secret string to encrypt cookies. Required in Production mode."`

	ServiceName  string `yaml:"service_name" yaml_comment:"Service name for 'install' command"`
	ServiceUser  string `yaml:"service_user" yaml_comment:"User for 'install' command"`
	ServiceGroup string `yaml:"service_group" yaml_comment:"Group for 'install' command"`
}

func (s *AppSettingsBase) checkDefaultValues(d *AppSettingsBase) {
	if s.WebserverHostname == "" {
		s.WebserverHostname = d.WebserverHostname
	}

	if s.WebserverPort == 0 {
		s.WebserverPort = d.WebserverPort
	}

	if s.ServiceName == "" {
		s.ServiceName = d.ServiceName
	}

	if s.ServiceUser == "" {
		s.ServiceUser = d.ServiceUser
	}

	if s.ServiceGroup == "" {
		s.ServiceGroup = d.ServiceGroup
	}
}
