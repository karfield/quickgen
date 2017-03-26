package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmdline := cli.NewApp()
	cmdline.Name = "quickgen"
	cmdline.Usage = "util for generating codes quickly"
	cmdline.HelpName = "guickgen"
	cmdline.Version = "1.0.0"
	cmdline.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "sync,s",
			Usage: "set '--sync' to sync templates before generate layouts",
		},
	}
	cmdline.Commands = []cli.Command{
		cli.Command{
			Name:   "sync",
			Usage:  "sync templates",
			Action: syncCommand,
		},
	}
	cmdline.Before = beforeRun

	configs := ScanConfigs()
	for _, config := range configs {
		subcmd := cli.Command{
			Name:   config.Name,
			Usage:  config.Description,
			Action: GenerateAction(config),
		}
		subcmd.Flags = []cli.Flag{
			cli.StringFlag{
				Name:  "workdir,w",
				Usage: "Give me a work dir",
				Value: cwd,
			},
		}
		for _, flag := range config.Flags {
			var option cli.Flag
			flagName := flag.Name
			if flag.Shortname != "" {
				flagName = flagName + "," + flag.Shortname
			}
			switch flag.Type {
			case "string", "":
				option = &cli.StringFlag{
					Name:  flagName,
					Value: flag.Default,
					Usage: flag.Usage,
				}
			case "int":
				var defaultValue int = 0
				if dv, err := strconv.ParseInt(flag.Default, 10, 32); err == nil {
					defaultValue = int(dv)
				}
				option = &cli.IntFlag{
					Name:  flagName,
					Value: defaultValue,
					Usage: flag.Usage,
				}
			case "bool", "boolean":
				option = &cli.BoolFlag{
					Name:  flagName,
					Usage: flag.Usage,
				}
			}
			if option != nil {
				subcmd.Flags = append(subcmd.Flags, option)
			}
		}
		cmdline.Commands = append(cmdline.Commands, subcmd)
	}

	cmdline.Run(os.Args)
}

func syncCommand(context *cli.Context) error {
	return SyncTemplates()
}

func beforeRun(context *cli.Context) error {
	if context.Bool("sync") {
		fmt.Fprintf(os.Stdout, "Sync templates ...")
		if err := SyncTemplates(); err != nil {
			fmt.Fprintf(os.Stdout, " failed.\n")
			fmt.Fprintf(os.Stderr, "sync failure reason: %v\n", err)
		} else {
			fmt.Fprintf(os.Stdout, " synced!\n")
		}
	}
	return nil
}
