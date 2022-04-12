package main

import (
  "os"
  "path/filepath"
  "io/ioutil"
  "io/fs"
  "html/template"
)

type pageInfo struct {
  outputPath string
  dataPath string
}

// These should be put in a structure and read from a config file.
// TASK: Need to add logic to handle the case where one or more of these
// consists for multiple levels (i.e., dir1/dir2).
const pageDir = "pages"
const staticDir = "static"
const sassDir = "sass"
const layoutDir = "layouts"
const dataDir = "data"
const outputDir = "public"
const cleanOutputDir = true
const preCommand = "./pre.sh"
const postCommand = "./post.sh"

var layouts []string
var layoutTemplate *template.Template
var globalData map[string]interface{}

func main() {
  // TASK: read config
  createOutputDir()
  executeCommand(preCommand)
  loadLayouts()
  globalData = loadGlobalData(dataDir)
  processPages(pageDir, globalData)
  executeCommand(postCommand)
}

func loadLayouts() {
  dirMustExist(layoutDir)
  layoutTemplate = template.New("").Funcs(template.FuncMap{
    "include": includeAction,
  })
  err := filepath.Walk(layoutDir, func(path string, fileInfo fs.FileInfo, err error) error {
    // Ignore non-html files.
    if filepath.Ext(path) != ".html" {
      return nil
    }
    b, ferr := ioutil.ReadFile(path)
    checkError(ferr)
    layoutTemplate, err = layoutTemplate.Parse(string(b))
    checkError(err)
		return nil
  })
  checkError(err)
}

func copyFile(src, target string) {
  input, err := ioutil.ReadFile(src)
  checkError(err)
  err = ioutil.WriteFile(target, input, 0644)
  checkError(err)
}

func createOutputDir() {
  if cleanOutputDir {
    err := os.RemoveAll(outputDir)
    checkError(err)
  }
  createDir(outputDir)
}
