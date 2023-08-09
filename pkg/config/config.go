package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Host             string `mapstructure:"Host"`
	Port             string `mapstructure:"Port"`
	LogLevel         string `mapstructure:"LOG_LEVEL"`
	LogDirPath       string
	DataFolderPath   string
	ReplayFolderPath string
}

var GlobalConfig *Config

func Setup(configPath string) {
	conf := getDefaultConfig()
	loadConfigFromEnv(&conf)
	loadConfigFromFile(configPath, &conf)
	GlobalConfig = &conf
}

func getDefaultConfig() Config {
	rootPath := getPwdDirPath()
	dataFolderPath := filepath.Join(rootPath, "data")
	replayFolderPath := filepath.Join(dataFolderPath, "replays")
	logDirPath := filepath.Join(dataFolderPath, "logs")

	folders := []string{dataFolderPath, replayFolderPath, logDirPath}
	for _, v := range folders {
		if err := EnsureDirExist(v); err != nil {
			log.Fatalf("Create folder failed: %s", err)
		}
	}
	return Config{
		Host:             "localhost",
		Port:             "8083",
		LogLevel:         "INFO",
		LogDirPath:       logDirPath,
		DataFolderPath:   dataFolderPath,
		ReplayFolderPath: replayFolderPath,
	}

}

func EnsureDirExist(path string) error {
	if !haveDir(path) {
		if err := os.MkdirAll(path, 0700); err != nil {
			return err
		}
	}
	return nil
}

func haveDir(file string) bool {
	fi, err := os.Stat(file)
	return err == nil && fi.IsDir()
}

func have(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getPwdDirPath() string {
	if rootPath, err := os.Getwd(); err == nil {
		return rootPath
	}
	return ""
}

func loadConfigFromEnv(conf *Config) {
	viper.AutomaticEnv()
	if err := viper.Unmarshal(conf); err == nil {
		log.Println("Load config from env")
	}
}

func loadConfigFromFile(path string, conf *Config) {
	var err error
	if have(path) {
		fileViper := viper.New()
		fileViper.SetConfigFile(path)
		if err = fileViper.ReadInConfig(); err == nil {
			if err = fileViper.Unmarshal(conf); err == nil {
				log.Printf("Load config from %s success\n", path)
				return
			}
		}
	}
	if err != nil {
		log.Fatalf("Load config from %s failed: %s\n", path, err)
	}
}
