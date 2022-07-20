package main

import (
  "io/ioutil"
  "log"
  "gopkg.in/yaml.v3"
  "github.com/brothertoad/btu"
)

const JINJA_FORMAT  = "jinja"
const GOLANG_FORMAT = "go"

type gossConfig struct {
  TemplateFormat string `yaml:"templateFormat"`
  PageDir string `yaml:"pageDir"`
  LayoutDir string `yaml:"layoutDir"`
  DataDir string `yaml:"dataDir"`
  OutputDir string `yaml:"outputDir"`
  Clean bool `yaml:"clean"`
  Pre interface{} `yaml:"pre"`
  Post interface{} `yaml:"post"`
}

const DEFAULT_CONFIG_FILE = "goss.yaml"

func initConfig(config *gossConfig) {
  config.TemplateFormat = JINJA_FORMAT
  config.PageDir = "pages"
  config.LayoutDir = "layouts"
  config.DataDir = "data"
  config.OutputDir = "public"
  config.Clean = true
  config.Pre = ""
  config.Post = ""
  loadConfig(config, DEFAULT_CONFIG_FILE, false)
}

func loadConfig(config *gossConfig, path string, fileMustExist bool) {
  if !btu.FileExists(path) {
    if fileMustExist {
      log.Fatalf("Config file %s does not exist.\n", path)
    }
    return
  }
  b, err := ioutil.ReadFile(path)
  btu.CheckError(err)
  err = yaml.Unmarshal(b, config)
  btu.CheckError(err)
}
