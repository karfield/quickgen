package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
			data[flag.Name] = c.Cmd.String(flag.Name)
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
