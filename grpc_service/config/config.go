package config

import (
	"flag"
	"path/filepath"

	"github.com/airborne12/golang-samples/grpc_service/log"
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
	ClientCA         string
	CredStoreAddress string
	CredStoreCA      string
	SN               string
	Insecure         bool
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
	flag.String("serverCert", "certificates/server.crt", "server TLS cert")
	flag.String("serverKey", "certificates/server.key", "server TLS key")
	flag.String("clientCA", "", "client CA")
	flag.String("allowedCN", "", "allowed CommonName for client authentication")
	flag.String("credstoreAddress", "", "credstore grpc address")
	flag.String("credstoreCA", "", "credstore server ca")
	flag.String("sn", "example.com", "server name for ssl")

	flag.Bool("insecure", false, "https on or off")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetDefault("listenAddr", "localhost:2503")
	viper.SetDefault("serverCert", "")
	viper.SetDefault("serverKey", "")
	viper.SetDefault("clientCA", "")
	viper.SetDefault("allowedCN", "")
	viper.SetDefault("credStoreAddress", "")
	viper.SetDefault("credStoreCA", "example.com")
	viper.SetDefault("sn", "")
	viper.SetDefault("insecure", false)

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
	c.ServerCert = viper.GetString("serverCert")
	c.ServerKey = viper.GetString("serverKey")
	c.AllowedCN = viper.GetString("allowedCN")
	c.ClientCA = viper.GetString("clientCA")
	c.CredStoreAddress = viper.GetString("credStoreAddress")
	c.CredStoreCA = viper.GetString("credStoreCA")
	c.SN = viper.GetString("sn")
	c.Insecure = viper.GetBool("insecure")
	return
}
