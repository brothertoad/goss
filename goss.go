package main

import (
  "os"
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/btu"
)

const configFlag = "config"

// TASK: Need to add logic to handle the case where one or more of the directories
// consists for multiple levels (i.e., dir1/dir2).

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
    },
    Action: gossMain,
  }
  app.Run(os.Args)
}

func gossMain(c *cli.Context) error {
  initConfig(&config)
  if c.String(configFlag) != "" {
    loadConfig(&config, c.String(configFlag), true)
  }
  createOutputDir(config.OutputDir, config.Clean)
  executeCommand(config.Pre)
  globalData = loadGlobalData(config.DataDir)
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
