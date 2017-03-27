package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

type GeneratorContext struct {
	Cmd         *cli.Context
	Config      *Config
	Step        *Step
	WorkDir     string
	TemplateDir string
}

func GenerateAction(config *Config) cli.ActionFunc {
	return func(context *cli.Context) error {
		if isWorkDirEmpty() {
			reader := bufio.NewReader(os.Stdin)
			fmt.Fprintf(os.Stderr, "WorkDir does not clean, are you sure to continue? (Y/n)")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(strings.ToLower(text), " \t\n\r")
			switch text {
			case "y", "yes", "":
				//fmt.Fprintf(os.Stdout, "Continue this command.\n")
			default:
				return errors.New("stopped by user")
			}
		}
		for idx, step := range config.Steps {
			gc := &GeneratorContext{
				Config:      config,
				Cmd:         context,
				Step:        &step,
				WorkDir:     context.String("workdir"),
				TemplateDir: config.Path,
			}
			var err error
			fmt.Printf("Step %d: ", idx)
			switch step.Type {
			case "copy":
				fmt.Printf("copy %s", step.Res)
				err = gc.copyAndApplyTemplate()
			case "mkdir", "mkdirs":
				fmt.Printf("mkdir %v", step.Dirs)
				err = gc.makeDirs()
			}
			fmt.Print("\n")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Step.%d failed: %v, start rollback\n", idx, err)
				// FIXME: rollback here!
				break
			}
		}

		fmt.Println("Finished.")
		if config.FinishNotice != "" {
			fmt.Println(config.FinishNotice)
		}
		return nil
	}
}

func (c *GeneratorContext) copyAndApplyTemplate() error {
	content, err := ioutil.ReadFile(filepath.Join(c.TemplateDir, c.Step.Res))
	if err != nil {
		return err
	}
	tpl, err := template.New("generator").Parse(string(content))
	if err != nil {
		return err
	}

	dstPath := filepath.Join(c.WorkDir, c.Step.Target)
	makeSureCanCreateFile(dstPath)

	file, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	data := map[string]interface{}{}
	for _, flag := range c.Config.Flags {
		switch flag.Type {
		case "string", "":
			value := c.Cmd.String(flag.Name)
			if len(flag.Options) > 0 {
				hasOne := false
				for _, opt := range flag.Options {
					if opt == value {
						hasOne = true
						break
					}
				}
				if !hasOne {
					return errors.New(fmt.Sprintf("option '%s'(value: %s) not matched in %v", flag.Name, value, flag.Options))
				}
			}
			data[flag.Name] = value
		case "int":
			data[flag.Name] = c.Cmd.Int(flag.Name)
		case "bool", "boolean":
			data[flag.Name] = c.Cmd.Bool(flag.Name)
		}
	}

	return tpl.Execute(file, data)
}

func (c *GeneratorContext) makeDirs() (err error) {
	defer func() {
		if err != nil {
			for _, dir := range c.Step.Dirs {
				// FIXME: remove dirs
				dir = dir
			}
		}
	}()
	for _, dir := range c.Step.Dirs {
		absDir := filepath.Join(c.WorkDir, dir)
		err := os.MkdirAll(absDir, os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	return nil
}

func makeSureCanCreateFile(filename string) {
	parentDir := filepath.Dir(filename)
	if fileExist(parentDir) {
		return
	}
	if parentDir != "/" && parentDir != "" {
		os.MkdirAll(parentDir, os.FileMode(0755))
	}
}
