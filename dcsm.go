package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var composefile string
var servicedir string

const defaultComposeContent = `version: '3'
services:
`

func main() {
	flag.StringVar(&composefile, "f", "docker-compose.yml", "path to the unified docker-compose.yml file")
	flag.StringVar(&servicedir, "sdir", "services", "directory where service.yml files are held")
	flag.Parse()

	cf := parse(composefile)
	if cf.Services == nil {
		cf.Services = map[string]service{}
	}
	serviceFiles := parseServices(servicedir)

	for _, s := range serviceFiles {
		cf = merge(cf, s)
	}

	o, err := yaml.Marshal(cf)
	check(err)
	fmt.Printf("writing to %s\n", composefile)
	err = os.WriteFile(composefile, o, 0644)
	check(err)
	fmt.Println("done")
}

func check(err error) {
	if err != nil {
		panic(fmt.Sprintf("failed: %+v", err))
	}
}

func parse(filepath string) compose {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		os.WriteFile(filepath, []byte(defaultComposeContent), 0644)
	}
	c, err := os.ReadFile(filepath)
	check(err)
	var o compose
	err = yaml.Unmarshal(c, &o)
	check(err)
	return o
}

func parseServices(dir string) []compose {
	o := []compose{}

	files, err := os.ReadDir(dir)
	check(err)
	for _, f := range files {

		ext := filepath.Ext(f.Name())

		if !f.IsDir() && (ext == ".yaml" || ext == ".yml") {
			o = append(o, parse(filepath.Join(dir, f.Name())))
		}
	}
	return o
}

func merge(o compose, n compose) compose {
	for k, v := range n.Services {
		if _, ok := o.Services[k]; ok {
			fmt.Printf("[%s]already present, will replace\n", k)
		} else {
			fmt.Printf("[%s] will be added\n", k)
		}
		if v.Restart == "" {
			v.Restart = "unless-stopped"
		}
		if v.Hostname == "" {
			v.Hostname = k
		}
		o.Services[k] = v
	}
	return o
}

type compose struct {
	Version  string             `yaml:"version"`
	Services map[string]service `yaml:"services"`
}

type service struct {
	Image   string   `yaml:"image"`
	Volumes []string `yaml:"volumes,omitempty"`
	Ports   []string `yaml:"ports,omitempty"`

	Restart  string `yaml:"restart"`  // default to unless-stopped
	Hostname string `yaml:"hostname"` // default to service name
}
