package main

import (
  "bytes"
  "os"
  "strings"
  "io/fs"
  "io/ioutil"
  "path/filepath"
  "github.com/adrg/frontmatter"
  "gopkg.in/yaml.v3"
  "github.com/noirbizarre/gonja"
  "github.com/brothertoad/btu"
)

type pageInfo struct {
  outputPath string
  dataPath string
  perPageDataPath string
}

// Default format for page modification date.
// Should make this configurable.
const TEXTDATE = "January 2, 2006"

// key for page modification date
const DATE_MODIFIED_KEY = "modTime"

var numPageDirParts int

func processPages(pageDir string, outputDir string, globalData map[string]interface{}) {
  btu.DirMustExist(pageDir)
  err := filepath.Walk(pageDir, func(path string, fileInfo fs.FileInfo, err error) error {
    // Check the extension and filename to see if we should ignore this file.
    if !validExtension(fileInfo.Name()) || !validFilename(fileInfo.Name()) {
      return nil
    }

    info := buildPageInfo(path, outputDir, fileInfo)

    // Copy the global data, then read and merge the frontmatter.
    pageData := make(map[string]interface{})
    for k, v := range globalData {
      pageData[k] = v
    }
    b, ferr := ioutil.ReadFile(path)
    btu.CheckError(ferr)
    fm := make(map[string]interface{})
    rest, fmerr := frontmatter.Parse(bytes.NewReader(b), &fm)
    btu.CheckError(fmerr)
    for k, v := range fm {
      pageData[k] = v
    }

    // Read in page-specific data, if any.
    b, ferr = ioutil.ReadFile(info.dataPath)
    if ferr == nil {
      pageMap := make(map[string]interface{})
      yerr := yaml.Unmarshal(b, pageMap)
      btu.CheckError(yerr)
      for k, v := range pageMap {
        pageData[k] = v
      }
    }

    // Get the last-modified date of the file, and add it as a page-specific property.
    pageData[DATE_MODIFIED_KEY] = fileInfo.ModTime().Format(TEXTDATE)

    tpl := gonja.Must(gonja.FromBytes(rest))
    out, err := tpl.Execute(pageData)
    btu.CheckError(err)
    btu.CreateDirForFile(info.outputPath)
    err = os.WriteFile(info.outputPath, []byte(out + "\n"), 0644)
    btu.CheckError(err)

    return nil
  })
  btu.CheckError(err)
}

func validExtension(filename string) bool {
  if strings.HasSuffix(filename, ".html") {
    return true
  }
  // Note that this case also covers .html.j2
  if strings.HasSuffix(filename, ".j2") {
    return true
  }
  return false
}

func validFilename(filename string) bool {
  return !strings.HasPrefix(filename, "__")
}

func buildPageInfo(path string, outputDir string, fileInfo fs.FileInfo) pageInfo {
  var info pageInfo
  relativePath := strings.TrimPrefix(path, config.PageDir)
  outputPath := filepath.Join(outputDir, relativePath)
  dir, base := filepath.Split(outputPath)
  suffix := getSuffix(base)
  // If the filename starts with an underscore, remove the underscore and the suffix and add .html.
  // Otherwise, remove the suffix and add a separator followed by index.html.
  if strings.HasPrefix(base, "_") {
    outputPath = filepath.Join(dir, strings.TrimSuffix(base[1:], suffix) + ".html")
  } else {
    outputPath = strings.TrimSuffix(outputPath, suffix)
    outputPath = filepath.Join(outputPath, "index.html")
  }
  info.outputPath = outputPath
  info.dataPath = strings.TrimSuffix(path, suffix) + ".yaml"
  if config.PerPageDataDir != "" {
    info.perPageDataPath = filepath.Join(config.PerPageDataDir, strings.TrimSuffix(relativePath, suffix) + ".yaml")
  }
  btu.Log(10, "path is %s, relativePath is %s, dataPath is %s, perPageDataPath is %s\n", path, relativePath, info.dataPath, info.perPageDataPath)
  return info
}

func getSuffix(base string) string {
  if strings.HasSuffix(base, ".html") {
    return ".html"
  }
  if strings.HasSuffix(base, ".html.j2") {
    return ".html.j2"
  }
  if strings.HasSuffix(base, ".j2") {
    return ".j2"
  }
  btu.Fatal("Can't get suffix of %s\n", base)
  return ""
}
