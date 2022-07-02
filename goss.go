package main

import (
  "os"
  "path/filepath"
  "io/ioutil"
  "io/fs"
  "html/template"
  "github.com/brothertoad/btu"
)

// TASK: Need to add logic to handle the case where one or more of the directories
// consists for multiple levels (i.e., dir1/dir2).

var config gossConfig
var layouts []string
var layoutTemplate *template.Template
var globalData map[string]interface{}

func main() {
  initConfig(&config)
  if len(os.Args) > 1 {
    loadConfig(&config, os.Args[1], true)
  }
  createOutputDir(config.OutputDir, config.Clean)
  executeCommand(config.Pre)
  loadLayouts(config.LayoutDir)
  globalData = loadGlobalData(config.DataDir)
  processPages(config.PageDir, config.OutputDir, globalData)
  executeCommand(config.Post)
}

func loadLayouts(layoutDir string) {
  btu.DirMustExist(layoutDir)
  layoutTemplate = template.New("").Funcs(template.FuncMap{
    "include": includeAction,
  })
  err := filepath.Walk(layoutDir, func(path string, fileInfo fs.FileInfo, err error) error {
    // Ignore non-html files.
    if filepath.Ext(path) != ".html" {
      return nil
    }
    b, ferr := ioutil.ReadFile(path)
    btu.CheckError(ferr)
    layoutTemplate, err = layoutTemplate.Parse(string(b))
    btu.CheckError(err)
		return nil
  })
  btu.CheckError(err)
}

func copyFile(src, target string) {
  input, err := ioutil.ReadFile(src)
  btu.CheckError(err)
  err = ioutil.WriteFile(target, input, 0644)
  btu.CheckError(err)
}

func createOutputDir(outputDir string, clean bool) {
  if clean {
    err := os.RemoveAll(outputDir)
    btu.CheckError(err)
  }
  createDir(outputDir)
}
