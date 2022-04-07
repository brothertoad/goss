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
}

// These should be put in a structure and read from a config file.
// Probably need to add logic to handle the case where one or more of these
// consists for multiple levels (i.e., dir1/dir2).
const pageDir = "pages"
const staticDir = "static"
const layoutDir = "layouts"
const outputDir = "public"
const cleanOutputDir = true

const PathSeparatorString = string(os.PathSeparator)

var layouts []string
var layoutTemplate *template.Template

func main() {
  // read config

  // create/clean the output directory, if needed
  createOutputDir()

  // load layouts
  loadLayouts()
  loadLayouts3()

  // load global data
  // for each page, load it and any page-specific data, and process it
  processPages2()
  // copy static files, possibly using an external program, such as rsync or rclone
  // process scss files
}

func loadLayouts3() {
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

func loadLayouts() {
  // Need to verify directory exists
  // Need to load files as byte arrays
  err := filepath.Walk(layoutDir,
    func(path string, fileInfo fs.FileInfo, err error) error {
      if strings.HasSuffix(fileInfo.Name(), ".html") {
        layouts = append(layouts, path)
      }
      return nil
    })
  if err != nil {
    log.Fatal(err)
  }
}

func processPages2() {
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
    fmt.Printf("Processing page file %s, which is named %s\n", path, fileInfo.Name())
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
  fmt.Printf("parts of page file: ")
  fmt.Println(parts)
  // TASK: handle the case where pageDir and/or outputDir has multiple components
  parts[0] = outputDir
  fmt.Printf("parts of output file before checking base: ")
  fmt.Println(parts)
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
  fmt.Printf("Output path is %s\n", info.outputPath)
  return info, nil
}

func processPages() {
  // Need to verify directory exists
  err := filepath.Walk(pageDir,
    func(path string, fileInfo fs.FileInfo, err error) error {
      if strings.HasSuffix(fileInfo.Name(), ".html") && !strings.HasPrefix(fileInfo.Name(), "_") {
        // process the file, and write to output directory, recreating the path
        pageOutputPath := buildOutputDirForFile(path, pageDir) + PathSeparatorString + fileInfo.Name()
        fmt.Printf("Found source page file %s, output page is %s\n", path, pageOutputPath)
        // Read frontmatter and merge with global data
        // Read page-specific data and merge with frontmatter/global data
        // Parse template
        s := []string{path}
        s = append(s, layouts...)
        ts, tsErr := template.ParseFiles(s...)
        if tsErr != nil {
          log.Fatal(tsErr)
        }
        writer, fileError := os.Create(pageOutputPath)
        if fileError != nil {
          log.Fatal(fileError)
        }
        tsErr = ts.Execute(writer, nil)
        if tsErr != nil {
          log.Fatal(tsErr)
        }
        writer.Close()
      }
      return nil
    })
  if err != nil {
    log.Fatal(err)
  }
}

func buildOutputDirForFile(path string, prefix string) string {
  // initial index is just after the prefix - we expect that the first character
  // we keep will be the path separator
  m := len(prefix)
  // last index is the last path separator, less 1 (so we don't keep the separator)
  n := strings.LastIndex(path, PathSeparatorString)
  // now create the outputPath
  outputPath := outputDir + path[m:n]
  err := os.MkdirAll(outputPath, 0755)
  if err != nil {
    log.Fatal(err)
  }
  return outputPath
}

func copyFile(src, target string) {
  input, err := ioutil.ReadFile(src)
  if err != nil {
    log.Fatal(err)
  }
  err = ioutil.WriteFile(target, input, 0644)
  if err != nil {
    log.Fatal(err)
  }
}

func createOutputDir() {
  if cleanOutputDir {
    err := os.RemoveAll(outputDir)
    if err != nil {
      log.Fatal(err)
    }
  }
  err := os.MkdirAll(outputDir, 0755)
  if err != nil {
    log.Fatal(err)
  }
}

func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
