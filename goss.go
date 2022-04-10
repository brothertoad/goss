package main

import (
  "os"
  "fmt"
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
// Probably need to add logic to handle the case where one or more of these
// consists for multiple levels (i.e., dir1/dir2).
const pageDir = "pages"
const staticDir = "static"
const layoutDir = "layouts"
const dataDir = "data"
const outputDir = "public"
const cleanOutputDir = true

var layouts []string
var layoutTemplate *template.Template
var globalData map[string]interface{}

func main() {
  // read config

  createOutputDir()
  loadLayouts()
  globalData = loadGlobalData(dataDir)
  processPages(globalData)
  // copy static files, possibly using an external program, such as rsync or rclone
  // process scss files
}

func loadLayouts() {
  dirMustExist(layoutDir)
  layoutTemplate = template.New("").Funcs(template.FuncMap{
    "include": includeCommand,
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
  fmt.Println(layoutTemplate.DefinedTemplates())
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
