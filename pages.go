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
)

func processPages(globalData map[string]interface{}) {
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
    info := buildPageInfo(path, fileInfo)
    // Clone the layout template, so we don't add have residue from previous pages.
    t, terr := layoutTemplate.Clone()
    checkError(terr)

    // Copy the global data, then read and merge the frontmatter.
    pageData := make(map[string]interface{})
    for k, v := range globalData {
      pageData[k] = v
    }
    b, ferr := ioutil.ReadFile(path)
    checkError(ferr)
    fm := make(map[string]interface{})
    rest, fmerr := frontmatter.Parse(bytes.NewReader(b), &fm)
    checkError(fmerr)
    for k, v := range fm {
      pageData[k] = v
    }

    // TASK: Need to read in page-specific data, if any.

    // Parse what was left of the page after removing the frontmatter.
    t, terr = t.Parse(string(rest))
    checkError(terr)

    // Determine the layout
    layoutValue, ok := pageData["layout"]
    if !ok {
      log.Fatalf("No layout for page %s\n", path)
    }
    layout := fmt.Sprintf("%v", layoutValue)

    // Now we can execute the template and write the output.
    fmt.Printf("Output will be written to %s\n", info.outputPath)
    createDirForFile(info.outputPath)
    file, ferr := os.Create(info.outputPath)
    checkError(nil)
    defer file.Close()
    w := bufio.NewWriter(file)
    t.ExecuteTemplate(w, layout, pageData)
    w.Flush()
    return nil
  })
  checkError(err)
}

func buildPageInfo(path string, fileInfo fs.FileInfo) pageInfo {
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
  return info
}
