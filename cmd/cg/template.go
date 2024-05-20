package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/valyala/fasttemplate"
)

// SchemaFileName is code template name
const SchemaFileName string = "schema.json"

// Add would copy template to TemplatePath
func Add(fromPath, toName string) error {
	err := Test(fromPath)
	if err != nil {
		return err
	}
	config := GetConfig()
	if toName == "" {
		toName = filepath.Base(fromPath)
	}
	fmt.Printf("add template [%s] as %s\n", fromPath, toName)
	return CopyDir(fromPath, filepath.Join(config.TemplatePath, toName))
}

// Remove remove a template
func Remove(name string) error {
	config := GetConfig()
	return os.RemoveAll(filepath.Join(config.TemplatePath, name))
}

// GetTemplateList get all template list
func GetTemplateList() ([]string, error) {
	config := GetConfig()
	dir, err := os.ReadDir(config.TemplatePath)
	if err != nil {
		return nil, err
	}
	var list []string
	for _, file := range dir {
		if file.IsDir() {
			schemaExist := CheckPathExists(filepath.Join(config.TemplatePath, file.Name(), SchemaFileName))
			if !schemaExist {
				continue
			}
			list = append(list, file.Name())
		}
	}
	return list, nil
}

// GetAllSchema get all schema
func GetAllSchema() (map[string]*Schema, error) {
	list, err := GetTemplateList()
	if err != nil {
		return nil, err
	}
	var result = make(map[string]*Schema)
	for _, name := range list {
		schemas, err := GetSchemas(name)
		if err != nil {
			return nil, err
		}
		fileName := name
		result[fileName] = schemas
	}
	return result, nil
}

// GetSchemas is generate code by template
func GetSchemas(name string) (*Schema, error) {
	schemaPath := GetSchemaPath(name)
	exist := CheckPathExists(schemaPath)
	if !exist {
		return nil, fmt.Errorf("schema is not found for [%s] in [%s]", name, schemaPath)
	}
	schemas, err := ReadSchema(schemaPath)
	if err != nil {
		return nil, err
	}
	return schemas, nil
}

func GetSchemaByPath(schemaDir string) (*Schema, error) {
	schemaPath := filepath.Join(schemaDir, SchemaFileName)
	exist := CheckPathExists(schemaPath)
	if !exist {
		return nil, fmt.Errorf("schema is not found in [%s]", schemaPath)
	}
	schemas, err := ReadSchema(schemaPath)
	if err != nil {
		return nil, err
	}
	return schemas, nil
}

// Generate generates code by template name
func Generate(toPath, name string, data map[string]interface{}, driver *Schema) error {
	config := GetConfig()
	fromPath := filepath.Join(config.TemplatePath, name)
	return GenerateByPath(toPath, fromPath, data, driver)
}

// Generate generates code by template path
func GenerateByPath(toPath, fromPath string, data map[string]interface{}, driver *Schema) error {
	var config *Config
	if driver == nil {
		schema, err := GetSchemaByPath(fromPath)
		if err != nil {
			return err
		}
		driver = schema
	}
	config = driver.Config
	baseConfig := GetConfig()
	if config == nil {
		config = baseConfig
	} else {
		config.Merge(baseConfig)
	}
	if !CheckPathExists(fromPath) {
		return fmt.Errorf("template is not exist in [%s]", fromPath)
	}
	err := CopyDirWithFunc(fromPath, toPath, func(s string) string {
		t := fasttemplate.New(s, config.FileNameTag, config.FileNameTag)
		return t.ExecuteString(data)
	}, func(dst io.Writer, src io.Reader) (written int64, err error) {
		reader := bufio.NewReader(src)
		var content string
		for {
			line, err := reader.ReadString('\n')
			content = content + line
			if err != nil {
				break
			}
		}

		switch driver.Driver {
		case "text/template":
			return io.Copy(dst, TextTemplate(content, data, config))
		case "liquid":
			return io.Copy(dst, LiquidTemplate(content, data, config))
		case "fasttemplate":
			return io.Copy(dst, FastTemplate(content, data, config))
		default:
			return io.Copy(dst, FastTemplate(content, data, config))
		}

	}, func(path string) bool {
		return path == SchemaFileName
	})
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("✅ template \"%v\" generated in %v", driver.Name, toPath))
	return nil
}

// Test would check path is correct template
func Test(fromPath string) error {
	schemaExist := CheckPathExists(filepath.Join(fromPath, SchemaFileName))
	if !schemaExist {
		return errors.New("schema file not found")
	}
	files, err := os.ReadDir(fromPath)
	if err != nil {
		return err
	}
	if len(files) < 2 {
		return errors.New("template dir should contain at least 2 files")
	}
	return nil
}
