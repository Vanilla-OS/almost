package core

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	Config  = "/etc/almost.ini"
	Section = "Almost"
	Defaults = map[string]interface{}{
		"Almost::CurrentMode": 0,
		"Almost:.DefaultMode": 0,
		"Almost::PersistModeStatus": false,
		"Almost::PkgManager::EntryPoint": "/usr/bin/apt",
	}
)

func init() {	
	if !RootCheck(false) {
		return
	}

	if _, err := os.Stat(Config); os.IsNotExist(err) {
		f, err := os.Create(Config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		f.Close()
	}
	viper.SetConfigFile(Config)
	viper.SetConfigType("ini")
	viper.SetDefault(Section, Defaults)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.WriteConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Load() error {
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func Set(key, value string) error {
	if err := Load(); err != nil {
		return err
	}
	viper.Set(Section+"."+key, value)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	if err := Load(); err != nil {
		return "", err
	}
	return viper.GetString(Section + "." + key), nil
}

func Show() error {
	if err := Load(); err != nil {
		return err
	}
	for k, v := range viper.GetStringMap(Section) {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
