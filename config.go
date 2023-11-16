package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppName cg
const AppName string = ".cg"
// TemplateDir is a dir name for template
const TemplateDir string = "templates"
// ConfigName is the name of config.json
const ConfigName string = AppName + "-config.json"
// SchemaName is the template schema file name
const SchemaName string = "schema.json"

// Config config.json
type Config struct {
	TemplatePath string `json:"templatePath"`
	FileNameTag  string `json:"fileNameTag"`
	FileStartTag string `json:"fileStartTag"`
	FileEndTag   string `json:"fileEndTag"`
}

// Schema schema.json
type Schema struct {
	Name        string   `json:"name"`
	Aliases     []string   `json:"aliases"`
	Description string   `json:"description"`
	Flags       []Flag    `json:"flags"`
}

// Flag flag
type Flag struct {
	Name        string   `json:"name"`
	Default     string   `json:"default"`
	Aliases     []string   `json:"aliases"`
	Regex       string   `json:"regex"`
	Description string   `json:"description"`
	Options     []string `json:"options"`
	Require     bool     `json:"require"`
}

var config *Config

func GetConfig() *Config {
	if config != nil {
		return config
	}
	config = InitConfig()
	return config
}

func GetSchemaPath(templateName string) string {
	config := GetConfig()
	return filepath.Join(config.TemplatePath, templateName, SchemaName)
}

func InitConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	baseDir := filepath.Join(home, AppName)
	err = MakeDirIfNotExist(baseDir)
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(baseDir, ConfigName)
	isConfigExist := CheckPathExists(configPath)
	if !isConfigExist {
		err = os.Mkdir(filepath.Join(baseDir, TemplateDir), 0755)
		if err != nil {
			panic(err)
		}
		config := Config{
			TemplatePath: filepath.Join(baseDir, TemplateDir),
			FileNameTag:  "__",
			FileStartTag: "{{",
			FileEndTag:   "}}",
		}
		WriteConfig(&config, configPath)
		return &config
	}
	config, err := ReadConfig(configPath)
	if err != nil {
		panic(err)
	}
	return config
}

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func ReadSchema(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var schemas Schema
	err = json.Unmarshal(data, &schemas)
	if err != nil {
		return nil, err
	}
	return &schemas, nil
}

func WriteConfig(config *Config, path string) error {
	data, err := json.Marshal(&config)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
