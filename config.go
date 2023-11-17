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
// Version tag app version. [be replaced by CI]
var Version string = "unknown"

// Config config.json
type Config struct {
	TemplatePath string `json:"templatePath"`
	FileNameTag  string `json:"fileNameTag"`
	FileStartTag string `json:"fileStartTag"`
	FileEndTag   string `json:"fileEndTag"`
	Version string `json:"version"`
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
var configMap map[string]interface{}

// GetConfig get config
func GetConfig() *Config {
	if config != nil {
		return config
	}
	config = InitConfig()
	return config
}

// GetConfigMap get config map
func GetConfigMap() (map[string]interface{}, error){
	if configMap != nil {
		return configMap, nil
	}
	config := GetConfig()
	jsonData, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, configMap)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

// ConfigMapInstance build config by configMap
func ConfigMapInstance(configMap map[string]interface{}) (*Config, error) {
	jsonData, err := json.Marshal(configMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetSchemaPath get schema path
func GetSchemaPath(templateName string) string {
	config := GetConfig()
	return filepath.Join(config.TemplatePath, templateName, SchemaName)
}

// GetPath get common path
func GetPath() (string, string) {
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
	return baseDir, configPath
}

// InitConfig init config
func InitConfig() *Config {
	baseDir, configPath := GetPath()
	isConfigExist := CheckPathExists(configPath)
	if !isConfigExist {
		err := os.Mkdir(filepath.Join(baseDir, TemplateDir), 0755)
		if err != nil {
			panic(err)
		}
		config := Config{
			TemplatePath: filepath.Join(baseDir, TemplateDir),
			FileNameTag:  "__",
			FileStartTag: "{{",
			FileEndTag:   "}}",
			Version: Version,
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

// ReadConfig read config by path
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

// ReadSchema read schema by path
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

// WriteConfigDefault writes default config
func WriteConfigDefault(config *Config){
	_, configPath := GetPath()
	WriteConfig(config, configPath)
}

// WriteConfig writes config by path
func WriteConfig(config *Config, path string) error {
	data, err := json.Marshal(&config)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// SetConfigItem set config
func SetConfigItem(key string, value interface{}) error {
	configMap, err := GetConfigMap()
	if err != nil {
		return err
	}
	configMap[key] = value
	config, err := ConfigMapInstance(configMap)
	if err != nil {
		return err
	}
	_, configPath := GetPath()
	return WriteConfig(config, configPath)
}

// GetConfigItem get config by key
func GetConfigItem(key string) (interface{}, error) {
	configMap, err := GetConfigMap()
	if err != nil {
		return nil, err
	}
	return configMap[key], nil
}
