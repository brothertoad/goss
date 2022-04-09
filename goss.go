package main

import (
  "os"
  "fmt"
  "strings"
  "path/filepath"
  "log"
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
const outputDir = "public"
const cleanOutputDir = true

var layouts []string
var layoutTemplate *template.Template

func main() {
  // read config

  // create/clean the output directory, if needed
  createOutputDir()

  loadLayouts()

  // load global data
  // for each page, load it and any page-specific data, and process it
  processPages()
  // copy static files, possibly using an external program, such as rsync or rclone
  // process scss files
}

func loadLayouts() {
  // Need to verify layouts directory exists
  layoutTemplate = template.New("")
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
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(layoutTemplate.DefinedTemplates())
}

func processPages() {
  // Need to verify pages directory exists
  err := filepath.Walk(pageDir, func(path string, fileInfo fs.FileInfo, err error) error {
    // Ignore non-html files.
    if filepath.Ext(path) != ".html" {
      return nil
    }
    // Ignore files with a leading underscore, except _index.html.
    if strings.HasPrefix(fileInfo.Name(), "_")  && fileInfo.Name() != "_index.html" {
      return nil
    }
    // Need to generate the output path, which mirrors the source path.
    info, _ := buildPageInfo(path, fileInfo)
    b, ferr := ioutil.ReadFile(path)
    checkError(ferr)
    // Clone the layout template, so we don't add have residue from previous pages.
    t, terr := layoutTemplate.Clone()
    checkError(terr)
    t, terr = t.Parse(string(b))
    checkError(terr)
    // Now we can execute the template and write the output.
    fmt.Printf("Output will be written to %s (may need to create directory)\n", info.outputPath)
    return nil
  })
  checkError(err)
}

func buildPageInfo(path string, fileInfo fs.FileInfo) (pageInfo, error) {
  var info pageInfo
  dir, base := filepath.Split(path)
  fmt.Printf("path is %s, dir is %s, base is %s\n", path, dir, base)
  parts := strings.Split(path, string(os.PathSeparator))
  // TASK: handle the case where pageDir and/or outputDir has multiple components
  parts[0] = outputDir
  // If the filename is _index.html, change it to index.html.  Otherwise, remove the
  // .html and add a separator followed by index.html.
  if base == "_index.html" {
    parts[len(parts)-1] = "index.html"
  } else {
    lastDir := strings.TrimSuffix(base, ".html")
    parts[len(parts)-1] = lastDir
    parts = append(parts, "index.html")
  }
  info.outputPath = filepath.Join(parts...)
  info.dataPath = strings.TrimSuffix(path, ".html") + ".yaml"
  fmt.Printf("Output path is %s, data path is %s\n", info.outputPath, info.dataPath)
  return info, nil
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
