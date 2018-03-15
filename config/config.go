package config

import "github.com/jinzhu/configor"

// Configuration is stuff that can be configured externally per env variables or config file (config.yml).
type Configuration struct {
	Server struct {
		Port int `default:"80"`
		SSL struct {
			Enabled         *bool  `default:"false"`
			RedirectToHTTPS *bool  `default:"true"`
			Port            int    `default:"443"`
			CertFile        string `default:""`
			CertKey         string `default:""`
			LetsEncrypt struct {
				Enabled   *bool  `default:"false"`
				AcceptTOS *bool  `default:"false"`
				Cache     string `default:"certs"`
				Hosts     []string
			}
		}
	}
	Database struct {
		Dialect    string `default:"sqlite3"`
		Connection string `default:"gotify.db"`
	}
	DefaultUser struct {
		Name string `default:"admin"`
		Pass string `default:"admin"`
	}
	PassStrength int `default:"10"`
}

// Get returns the configuration extracted from env variables or config file.
func Get() *Configuration {
	conf := new(Configuration)
	configor.New(&configor.Config{ENVPrefix: "GOTIFY"}).Load(conf, "config.yml", "/etc/gotify/config.yml")
	return conf
}