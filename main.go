package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dirs, err := ioutil.ReadDir(cwd)
	if err != nil {
		panic(err)
	}
	if len(dirs) > 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Fprintf(os.Stderr, "WorkDir does not clean, are you sure want to go on? (Y/n)")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(strings.ToLower(text), " \t\n\r")
		switch text {
		case "y", "yes", "":
			//fmt.Fprintf(os.Stdout, "Continue this command.\n")
		default:
			os.Exit(0)
		}
	}
	cmdline := cli.NewApp()
	cmdline.Name = "quickgen"
	cmdline.Usage = "util for generating codes quickly"
	cmdline.HelpName = "guickgen"
	cmdline.Version = "1.0.0"
	cmdline.Commands = []cli.Command{}

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
			switch flag.Type {
			case "string", "":
				option = &cli.StringFlag{
					Name:  flag.Name,
					Value: flag.Default,
					Usage: flag.Usage,
				}
			case "int":
				var defaultValue int = 0
				if dv, err := strconv.ParseInt(flag.Default, 10, 32); err == nil {
					defaultValue = int(dv)
				}
				option = &cli.IntFlag{
					Name:  flag.Name,
					Value: defaultValue,
					Usage: flag.Usage,
				}
			case "bool", "boolean":
				option = &cli.BoolFlag{
					Name:  flag.Name,
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
