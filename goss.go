package main

import (
  "fmt"
  "os"
  "gopkg.in/yaml.v3"
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/btu"
)

const configFlag = "config"
const logLevelFlag = "log-level"

var config gossConfig
var globalData map[string]interface{}

func main() {
  app := &cli.App {
    Name: "goss",
    Usage: "a simple static site generator",
    Flags: []cli.Flag {
      &cli.StringFlag {
        Name: configFlag,
        Usage: "configuration file",
      },
      &cli.StringFlag {
        Name: logLevelFlag,
        Usage: "log level",
      },
      &cli.BoolFlag {
        Name: "data2",
        Usage: "use level 2 data loading",
      },
      &cli.BoolFlag {
        Name: "yaml",
        Usage: "print global data in yaml format",
      },
    },
    Action: gossMain,
  }
  app.Run(os.Args)
}

func gossMain(c *cli.Context) error {
  if c.String(logLevelFlag) != "" {
    btu.SetLogLevelByName(c.String(logLevelFlag))
  }
  initConfig(&config)
  if c.String(configFlag) != "" {
    loadConfig(&config, c.String(configFlag), true)
  }
  createOutputDir(config.OutputDir, config.Clean)
  executeCommand(config.Pre)
  if c.Bool("data2") {
    globalData = loadGlobalData2(config.DataDir)
  } else {
    globalData = loadGlobalData(config.DataDir)
  }
  if c.Bool("yaml") {
    bytes, err := yaml.Marshal(globalData)
    btu.CheckError(err)
    fmt.Printf("%s\n", string(bytes[:]))
  }
  processPages(config.PageDir, config.OutputDir, globalData)
  executeCommand(config.Post)
  return nil
}

func createOutputDir(outputDir string, clean bool) {
  if clean {
    err := os.RemoveAll(outputDir)
    btu.CheckError(err)
  }
  btu.CreateDir(outputDir)
}
