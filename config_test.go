package main

import (
	"os"
	"path/filepath"
	"testing"
)

func getPath() (string, string, string){
    homeDir, _ := os.UserHomeDir()
    baseDir := filepath.Join(homeDir, AppName)
    defaultTemplatePath := filepath.Join(baseDir, TemplateDir)
    configPath := filepath.Join(baseDir, ConfigName)
    return baseDir, defaultTemplatePath, configPath
}

func TestReadConfig(t *testing.T) {
    _, defaultTemplatePath, configPath := getPath()
    config, _ := ReadConfig(configPath)

    if config.TemplatePath != defaultTemplatePath {
        t.Fatalf("`templatePath` field in config.json must be %s", defaultTemplatePath)
    }
}

func TestInitConfig(t *testing.T) {
    baseDir, defaultTemplatePath, _ := getPath()
	config := InitConfig()
	if config.TemplatePath != filepath.Join(baseDir, TemplateDir) {
        t.Fatalf("default template path must is [%s]", defaultTemplatePath)
    }
    configPath := filepath.Join(baseDir, ConfigName)
    isExist := CheckPathExists(configPath)
    if !isExist {
        t.Fatalf("config path [%s] must exist", configPath)
    }
}
