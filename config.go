package main

import (
  "io/ioutil"
  "log"
  "gopkg.in/yaml.v3"
)

type gossConfig struct {
  PageDir string `yaml:"pageDir"`
  LayoutDir string `yaml:"layoutDir"`
  DataDir string `yaml:"dataDir"`
  OutputDir string `yaml:"outputDir"`
  Clean bool `yaml:"clean"`
  Pre string `yaml:"pre"`
  Post string `yaml:"post"`
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
  if !fileExists(path) {
    if fileMustExist {
      log.Fatalf("Config file %s does not exist.\n", path)
    }
    return
  }
  b, err := ioutil.ReadFile(path)
  checkError(err)
  err = yaml.Unmarshal(b, config)
}
