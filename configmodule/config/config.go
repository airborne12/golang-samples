package config

import (
	"flag"
	"path/filepath"

	"github.com/airborne12/golang-samples/logmodule/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var logger = log.Get()

//Config struct contains config items for whole project
type Config struct {
	ListenAddr       string
	ServerCert       string
	ServerKey        string
	AllowedCN        string
	CredStoreAddress string
	CredStoreCA      string
}

var gConfig *Config

func init() {
	gConfig = newConfig()
	gConfig.parse()
}

func newConfig() *Config {
	return &Config{}
}

//Get returns a config object
func Get() *Config {
	return gConfig
}

func (c *Config) parse() {
	//command line args
	flag.String("listenAddr", "localhost:2503", "Address to listen")
	flag.String("configFile", "/etc/example/example.yaml", "Config file to load")
	flag.String("server-cert", "", "server TLS cert")
	flag.String("server-key", "", "server TLS key")
	flag.String("client-ca", "", "client CA")
	flag.String("allowed-cn", "", "allowed CommonName for client authentication")
	flag.String("credstore-address", "", "credstore grpc address")
	flag.String("credstore-ca", "", "credstore server ca")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetDefault("listenAddr", "localhost:2503")
	viper.SetDefault("serverCert", "")
	viper.SetDefault("serverKey", "")
	viper.SetDefault("clientCA", "")
	viper.SetDefault("allowedCN", "")
	viper.SetDefault("credStoreAddress", "")
	viper.SetDefault("credStoreCA", "")

	viper.SetConfigType("yaml")
	name := filepath.Base(viper.GetString("configFile"))
	path := filepath.Dir(viper.GetString("configFile"))
	viper.SetConfigName(name)    // name of config file (without extension)
	viper.AddConfigPath(path)    // path to look for the config file in
	viper.AddConfigPath("conf/") // optionally look for config in the working directory
	err := viper.ReadInConfig()  // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		logger.Errorf("Loading config file error:%s", err)
	}
	//if command line args and yaml config file duplicate defined, command first
	c.ListenAddr = viper.GetString("listenAddr")
	return
}
