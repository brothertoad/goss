package main

import (
  "bufio"
  "bytes"
  "fmt"
  "log"
  "os"
  "strings"
  "io/fs"
  "io/ioutil"
  "path/filepath"
  "github.com/adrg/frontmatter"
  "gopkg.in/yaml.v3"
  "github.com/noirbizarre/gonja"
  "github.com/Joker/hpp"
  "github.com/brothertoad/btu"
)

type pageInfo struct {
  outputPath string
  dataPath string
}

func processPages(pageDir string, outputDir string, globalData map[string]interface{}) {
  btu.DirMustExist(pageDir)
  err := filepath.Walk(pageDir, func(path string, fileInfo fs.FileInfo, err error) error {
    // Ignore non-html files.
    if filepath.Ext(path) != ".html" {
      return nil
    }
    // Ignore files with a leading underscore, except _index.html.
    if strings.HasPrefix(fileInfo.Name(), "_")  && fileInfo.Name() != "_index.html" {
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
    if ferr != nil {
      pageMap := make(map[string]interface{})
      yerr := yaml.Unmarshal(b, pageMap)
      btu.CheckError(yerr)
      for k, v := range pageMap {
        pageData[k] = v
      }
    }

    if config.TemplateFormat == GOLANG_FORMAT {
      // Clone the layout template, to avoid residue from previous pages.
      t, terr := layoutTemplate.Clone()
      btu.CheckError(terr)

      // Parse what was left of the page after removing the frontmatter.
      t, terr = t.Parse(string(rest))
      btu.CheckError(terr)

      // Determine the layout
      layoutValue, ok := pageData["layout"]
      if !ok {
        log.Fatalf("No layout for page %s\n", path)
      }
      layout := fmt.Sprintf("%v", layoutValue)

      // Now we can execute the template and write the output.
      btu.CreateDirForFile(info.outputPath)
      file := btu.CreateFile(info.outputPath)
      defer file.Close()
      w := bufio.NewWriter(file)
      t.ExecuteTemplate(w, layout, pageData)
      w.Flush()
    } else if config.TemplateFormat == JINJA_FORMAT {
      tpl := gonja.Must(gonja.FromBytes(rest))
      out, err := tpl.Execute(pageData)
      btu.CheckError(err)
      btu.CreateDirForFile(info.outputPath)
      err = os.WriteFile(info.outputPath, []byte(hpp.PrPrint(out)), 0644)
      btu.CheckError(err)
    }

    return nil
  })
  btu.CheckError(err)
}

func buildPageInfo(path string, outputDir string, fileInfo fs.FileInfo) pageInfo {
  var info pageInfo
  _, base := filepath.Split(path)
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
  return info
}
