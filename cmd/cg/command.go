package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/urfave/cli/v2"
)

var AddCmd = cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Usage:   "add a template dir",
	Action: func(cCtx *cli.Context) error {
		path := cCtx.Args().First()
		if path == "" {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			de, err := os.ReadDir(dir)
			if err != nil {
				return err
			}
			for _, d := range de {
				Add(filepath.Join(dir, d.Name()), "")
			}
			return nil
		}
		return Add(filepath.Join(path), cCtx.Args().Get(1))
	},
}

var RmCmd = cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Usage:   "remove a template dir",
	Action: func(cCtx *cli.Context) error {
		name := cCtx.Args().First()
		if name == "" {
			return fmt.Errorf("must provider a template name")
		}
		return Remove(name)
	},
}

var ListCmd = cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "list all template",
	Action: func(cCtx *cli.Context) error {
		list, err := GetTemplateList()
		if err != nil {
			return err
		}
		for _, v := range list {
			fmt.Println(v)
		}
		return nil
	},
}

func BuildConfCmd() cli.Command {
	fields := reflect.TypeOf(*config)
	var setCommand []*cli.Command
	var getCommand []*cli.Command
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		fieldName := field.Tag.Get("json")
		setCommand = append(setCommand, &cli.Command{
			Name:  fieldName,
			Usage: fmt.Sprintf("set config [%s]", fieldName),
			Action: func(ctx *cli.Context) error {
				return SetConfigItem(fieldName, ctx.Args().First())
			},
		})
		getCommand = append(getCommand, &cli.Command{
			Name:  fieldName,
			Usage: fmt.Sprintf("get config [%s]", fieldName),
			Action: func(ctx *cli.Context) error {
				val, err := GetConfigItem(fieldName)
				if err != nil {
					return err
				}
				fmt.Printf("[%s]:%v", fieldName, val)
				return nil
			},
		})
	}

	return cli.Command{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "operate config",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "list config",
				Action: func(ctx *cli.Context) error {
					config := GetConfig()
					jsonData, err := json.MarshalIndent(&config, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(jsonData))
					return nil
				},
			},
			{
				Name:        "set",
				Usage:       "set config",
				Subcommands: setCommand,
			},
			{
				Name:        "get",
				Usage:       "get config",
				Subcommands: getCommand,
			},
		},
	}
}

func BuildGenerateCmd() cli.Command {
	schemas, err := GetAllSchema()
	if err != nil {
		log.Fatal(err)
	}
	var command []*cli.Command
	for key, value := range schemas {
		fileName := key
		schema := value
		name := schema.Name
		if name == "" {
			name = fileName
		}
		flag := []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "the path for generated code",
			},
		}
		for _, item := range schema.Flags {
			flag = append(flag, &cli.StringFlag{
				Name:    item.Name,
				Value:   item.Default,
				Aliases: item.Aliases,
				Usage:   item.Description,
			})
		}
		command = append(command, &cli.Command{
			Name:        fileName,
			Aliases:     schema.Aliases,
			Description: schema.Description,
			Flags:       flag,
			Action: func(cCtx *cli.Context) error {
				data := make(map[string]interface{})
				for _, v := range schema.Flags {
					name := v.Name
					regex := v.Regex
					val := cCtx.String(name)
					if v.Require && val == "" {
						return fmt.Errorf("require for flag: [%s]", name)
					}
					if regex != "" {
						re := regexp.MustCompile(regex)
						isMatch := re.Match([]byte(val))
						if !isMatch {
							return fmt.Errorf("flag: [%s] must match the regex [%s]", name, regex)
						}
					}
					switch v.Type {
					case "raw":
						var rawData interface{}
						err := json.Unmarshal([]byte(val), &rawData)
						if err != nil {
							return err
						}
						data[name] = rawData
						continue
					case "json":
						file, err := os.Open(val)
						if err != nil {
							return err
						}
						defer file.Close()

						jsonBytes, err := io.ReadAll(file)
						if err != nil {
							return err
						}
						var rawData interface{}
						err = json.Unmarshal(jsonBytes, &rawData)
						if err != nil {
							return err
						}
						data[name] = rawData
						continue
					case "url":
						var rawData interface{}
						err := Request(val, &rawData)
						if err != nil {
							return err
						}
						data[name] = rawData
						continue
					case "openAPI":
						url, err := CheckURL(val)
						if err != nil {
							return err
						}
						loader := openapi3.NewLoader()
						result, err := loader.LoadFromURI(url)
						if err != nil {
							return err
						}
						data[name] = result
						continue
					}
					data[name] = val
				}
				path := cCtx.String("path")
				if path == "" {
					dir, err := os.Getwd()
					if err != nil {
						return err
					}
					path = dir
				}
				return Generate(path, fileName, data, schema)
			},
		})
	}

	return cli.Command{
		Name:        "generate",
		Aliases:     []string{"g"},
		Usage:       "generate by a template",
		Subcommands: command,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "template",
				Usage: "the path of template dir",
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "the path of target dir",
			},
			&cli.StringFlag{
				Name:  "data",
				Usage: "the template variable data",
			},
		},
		Action: func(ctx *cli.Context) error {
			templatePath := ctx.String("template")
			if templatePath == "" {
				return fmt.Errorf("flag: [%s] is required", "template")
			}
			toPath := ctx.String("path")
			if toPath == "" {
				dir, err := os.Getwd()
				if err != nil {
					return err
				}
				toPath = dir
			}
			dataStr := ctx.String("data")
			data := make(map[string]interface{})
			if dataStr != "" {
				json.Unmarshal([]byte(dataStr), &data)
			}
			return Generate(toPath, templatePath, data, nil)
		},
	}
}
