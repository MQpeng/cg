package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func main() {
	config := GetConfig()
	confCmd := BuildConfCmd()
	generateCmd := BuildGenerateCmd()
	app := &cli.App{
		Commands: []*cli.Command{
			&AddCmd,
			&RmCmd,
			&confCmd,
			&ListCmd,
			&generateCmd,
			{
				Name:  "clone",
				Usage: "git clone repo to template",
				Action: func(cCtx *cli.Context) error {
					clone := exec.Command("git", "clone", cCtx.Args().First(), config.TemplatePath)
					clone.Stdout = os.Stdout
					return clone.Run()
				},
			},
			{
				Name:  "pull",
				Usage: "git pull repo to template",
				Action: func(cCtx *cli.Context) error {
					pull := exec.Command("git", "-C", config.TemplatePath, "pull")
					pull.Stdout = os.Stdout
					return pull.Run()
				},
			},
		},
		Version: Version,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
