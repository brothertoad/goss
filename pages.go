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
}

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
  _, base := filepath.Split(path)
  parts := strings.Split(path, string(os.PathSeparator))
  btu.Debug("pageDir is %s\n", config.PageDir)
  btu.Debug("path is %s, parts is %+v\n", path, parts)
  suffix := getSuffix(base)
  // TASK: handle the case where pageDir has multiple components
  parts[0] = outputDir
  // If the filename starts with an underscore, remove the underscore and the suffix and add .html.
  // Otherwise, remove the suffix and add a separator followed by index.html.
  if strings.HasPrefix(base, "_") {
    parts[len(parts)-1] = strings.TrimSuffix(base[1:], suffix) + ".html"
  } else {
    parts[len(parts)-1] = strings.TrimSuffix(base, suffix)
    parts = append(parts, "index.html")
  }
  info.outputPath = filepath.Join(parts...)
  btu.Debug("outputPath is %s\n", info.outputPath)
  info.dataPath = strings.TrimSuffix(path, suffix) + ".yaml"
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
