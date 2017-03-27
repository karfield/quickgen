package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/BurntSushi/toml"
)

type Flag struct {
	Name      string
	Shortname string
	Type      string
	Default   string
	Options   []string
	Usage     string `toml:"desc" yaml:"desc"`
}

type Step struct {
	Type   string
	Dirs   []string
	Res    string
	Target string
}

type GenConfig struct {
	Name         string
	Description  string `toml:"desc" yaml:"desc"`
	Flags        []Flag
	Steps        []Step
	FinishNotice string
}

type Config struct {
	GenConfig
	Path string
}

func ParseFromFile(filename string) (*Config, error) {
	var config Config
	if strings.HasSuffix(filename, ".toml") {
		if _, err := toml.DecodeFile(filename, &config.GenConfig); err != nil {
			return nil, err
		} else {
			//fmt.Printf("%v\n", metadata)
		}
	} else {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(data, &config.GenConfig)
		if err != nil {
			return nil, err
		}
	}

	for _, flag := range config.Flags {
		if flag.Name == "" {
			return nil, errors.New(fmt.Sprintf("missing name for flag %v", flag))
		}
		switch flag.Type {
		case "string", "", "bool", "boolean":
		default:
			return nil, errors.New(fmt.Sprintf("illegal flag type %v", flag))
		}
	}

	config.Path = filepath.Dir(filename)

	for _, step := range config.Steps {
		switch step.Type {
		case "copy":
			if step.Res == "" || step.Target == "" {
				return nil, errors.New(fmt.Sprintf("missing copy res or target for step %v", step))
			}
			if !fileExist(filepath.Join(config.Path, step.Res)) {
				return nil, errors.New(fmt.Sprintf("missing resource file for %v", step))
			}
		}
	}

	return &config, nil
}

func ScanConfigs() []*Config {
	configs := []*Config{}
	templateDir := templateDir()
	if dirs, err := ioutil.ReadDir(templateDir); err != nil {
		return configs
	} else {
		for _, dir := range dirs {
			if dir.IsDir() {
				configDir := filepath.Join(templateDir, dir.Name())
				configPath := filepath.Join(configDir, "GEN.toml")
				if !fileExist(configPath) {
					configPath = filepath.Join(configDir, "GEN.yaml")
				}
				if fileExist(configPath) {
					if config, err := ParseFromFile(configPath); err != nil {
						fmt.Fprintf(os.Stderr, "parse config file(%s) error %v\n", configPath, err)
						continue
					} else {
						configs = append(configs, config)
					}

				}
			}
		}
	}
	return configs
}
