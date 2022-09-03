package main

import (
  "io/ioutil"
  "gopkg.in/yaml.v3"
  "github.com/brothertoad/btu"
)

type gossConfig struct {
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
      btu.Fatal("Config file %s does not exist.\n", path)
    }
    return
  }
  b, err := ioutil.ReadFile(path)
  btu.CheckError(err)
  err = yaml.Unmarshal(b, config)
  btu.CheckError(err)
}
