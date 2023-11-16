package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/urfave/cli/v2"
)

func main() {
	InitConfig()
	schemas, err := GetAllSchema()
	if err != nil {
		log.Fatal(err)
	}
	var command []*cli.Command
	for fileName, schema := range schemas {
		name := schema.Name
		if name == "" {
			name = fileName
		}
		flag := []cli.Flag {
			&cli.StringFlag{
				Name:  "path",
				Usage: "the path for generated code",
			},
		}
		for _, item := range schema.Flags {
			flag = append(flag, &cli.StringFlag{
				Name:  item.Name,
				Value: item.Default,
				Aliases: item.Aliases,
				Usage: item.Description,
			})
		}
		command = append(command, &cli.Command{
			Name: fileName,
			Aliases: schema.Aliases,
			Description: schema.Description,
			Flags: flag,
			Action: func(cCtx *cli.Context) error {
				data := make(map[string]interface{})
				for _, v := range schema.Flags {
					val := cCtx.String(v.Name)
					if v.Require && val == "" {
						return fmt.Errorf("require for flag: [%s]", v.Name)
					}
					if v.Regex != "" {
						re := regexp.MustCompile(v.Regex)
						isMatch := re.Match([]byte(val))
						if !isMatch {
							return fmt.Errorf("flag: [%s] must match the regex [%s]", v.Name, v.Regex)
						}
					}
					data[v.Name] = val
				}
				path := cCtx.String("path")
				if path == "" {
					dir, err := os.Getwd()
					if err != nil {
						return err
					}
					path = dir
				}
				return Generate(path, fileName, data)
			},
		})
	}
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a template dir",
				Action: func(cCtx *cli.Context) error {
					return Add(filepath.Join(cCtx.Args().First()), cCtx.Args().Get(1))
				},
			},
			{
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
			},
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "generate by a template",
				Subcommands: command,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
