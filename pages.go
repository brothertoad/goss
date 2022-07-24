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
      err = os.WriteFile(info.outputPath, []byte(out + "\n"), 0644)
      btu.CheckError(err)
    }

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
  if !strings.HasPrefix(filename, "_") {
    return true
  }
  if filename == "_index.html" || filename == "_index.html.j2" || filename == "_index.j2" {
    return true
  }
  return false
}

func buildPageInfo(path string, outputDir string, fileInfo fs.FileInfo) pageInfo {
  var info pageInfo
  _, base := filepath.Split(path)
  parts := strings.Split(path, string(os.PathSeparator))
  suffix := getSuffix(base)
  // TASK: handle the case where pageDir and/or outputDir has multiple components
  parts[0] = outputDir
  // If the filename is _index.html or one of its variants, change it to index.html.
  // Otherwise, remove the suffix and add a separator followed by index.html.
  if base == "_index.html" || base == "_index.html.j2" || base == "_index.j2" {
    parts[len(parts)-1] = "index.html"
  } else {
    parts[len(parts)-1] = strings.TrimSuffix(base, suffix)
    parts = append(parts, "index.html")
  }
  info.outputPath = filepath.Join(parts...)
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
  log.Fatalf("Can't get suffix of %s\n", base)
  return ""
}
