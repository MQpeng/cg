package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/magiconair/properties/assert"
)

func getPath() (string, string, string) {
	InitConfig()
	homeDir, _ := os.UserHomeDir()
	baseDir := filepath.Join(homeDir, AppName)
	defaultTemplatePath := filepath.Join(baseDir, TemplateDir)
	configPath := filepath.Join(baseDir, ConfigName)
	return baseDir, defaultTemplatePath, configPath
}

func TestReadConfig(t *testing.T) {
	_, defaultTemplatePath, configPath := getPath()
	config, _ := ReadConfig(configPath)

	assert.Equal(t, config.TemplatePath, defaultTemplatePath)
}

func TestInitConfig(t *testing.T) {
	baseDir, _, _ := getPath()
	config := InitConfig()
	assert.Equal(t, config.TemplatePath, filepath.Join(baseDir, TemplateDir))

	configPath := filepath.Join(baseDir, ConfigName)
	isExist := CheckPathExists(configPath)

	assert.Equal(t, isExist, true)
}
