package main

import (
  "os"
  "path/filepath"
  "io/fs"
  "html/template"
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/btu"
)

const configFlag = "config"

// TASK: Need to add logic to handle the case where one or more of the directories
// consists for multiple levels (i.e., dir1/dir2).
// TASK: Use text/template rather than html/template, since we don't process
// user-supplied data (as a non-static site might).

var config gossConfig
var layouts []string
var layoutTemplate *template.Template
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
  loadLayouts(config.LayoutDir)
  globalData = loadGlobalData(config.DataDir)
  processPages(config.PageDir, config.OutputDir, globalData)
  executeCommand(config.Post)
  return nil
}

func loadLayouts(layoutDir string) {
  btu.DirMustExist(layoutDir)
  layoutTemplate = template.New("").Funcs(template.FuncMap{
    "include": btu.ReadFileS,
  })
  err := filepath.Walk(layoutDir, func(path string, fileInfo fs.FileInfo, err error) error {
    // Ignore non-html files.
    if filepath.Ext(path) != ".html" {
      return nil
    }
    layoutTemplate, err = layoutTemplate.Parse(btu.ReadFileS(path))
    btu.CheckError(err)
		return nil
  })
  btu.CheckError(err)
}

func createOutputDir(outputDir string, clean bool) {
  if clean {
    err := os.RemoveAll(outputDir)
    btu.CheckError(err)
  }
  btu.CreateDir(outputDir)
}
