package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config
var Config struct {
	Version string `mapstructure:"version"`
	Common  struct {
		NodeName         string `mapstructure:"node_name"`
		RunMode          string `mapstructure:"run_mode"`
		DataDir          string `mapstructure:"data_dir"`
		MasterAddr       string `mapstructure:"master_addr"`
		AuthServerAddr   string `mapstructure:"auth_server_addr"`
		TunnelListenPort string `mapstructure:"tunnel_listen_port"`
		WebListenPort    string `mapstructure:"web_listen_port"`
		Storage          struct {
			DBType string `mapstructure:"db_type"`
		} `mapstructure:"storage"`
		Cache struct {
			CacheType string `mapstructure:"cache_type"`
		} `mapstructure:"cache"`
	} `mapstructure:"common"`
	SSHService struct {
		Enabled bool              `mapstructure:"enabled"`
		Labels  map[string]string `mapstructure:"labels"`
	} `mapstructure:"ssh_service"`
	AuthService struct {
		Enabled              bool   `mapstructure:"enabled"`
		AuthServerListenPort string `mapstructure:"auth_server_listen_port"`
		WebServerListenPort  string `mapstructure:"web_server_listen_port"`
	} `mapstructure:"auth_service"`
}

// ParseConfig 配置解析
func ParseConfig() {
	vp := viper.New()
	vp.AddConfigPath("/etc")
	vp.AddConfigPath("./etc")
	vp.AddConfigPath("../etc")
	vp.AddConfigPath("../../etc")
	vp.SetConfigName("secops")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = vp.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
}
