package main

import (
  "io/fs"
  "os"
  "path/filepath"
  "strings"
  "gopkg.in/yaml.v3"
  "github.com/brothertoad/btu"
  "github.com/brothertoad/goss/gossutil"
)

func loadGlobalData2(dataDir string) map[string]interface{} {
  data := make(map[string]interface{})
  loadDataDir(dataDir, data)
  btu.Info("global data: %+v\n", data)
  return data
}

func loadDataDir(dataDir string, m map[string]interface{}) {
  btu.Debug("loadDataDir %s\n", dataDir)
  pattern := dataDir + string(filepath.Separator) + "*"
  matches, err := filepath.Glob(pattern)
  btu.CheckError(err)
  if matches == nil {
    return
  }
  btu.SetLogLevel(btu.INFO)
  for _, path := range(matches) {
    filename := getRelativePath(path, dataDir)
    if filepath.Ext(path) == ".yaml" {
      btu.Info("Found a yaml file, path is %s, filename is %s\n", path, filename)
      b := btu.ReadFileB(path)
      // Get base name of the file, use that as the key in data, unless file begins with
      // an underscore.
      var yerr error
      if strings.HasPrefix(filename, "_") {
        yerr = yaml.Unmarshal(b, m)
      } else {
        key := strings.TrimSuffix(filename, ".yaml")
        m[key] = make(map[string]interface{})
        yerr = yaml.Unmarshal(b, m[key])
      }
      btu.CheckError(yerr)
    } else if btu.IsDir(path) {
      btu.Info("%s (filename %s) is a directory\n", path, filename)
      if strings.HasPrefix(filename, "_") {
        // Need to concatenate yaml files in directory.
        btu.Info("Need to concatenate files in %s\n", path)
        yamls := gossutil.CatFiles(path, "yaml")
        if yamls == nil {
          btu.Fatal("Got nil when catting files.\n")
        }
        btu.Info("Found %d bytes in %s\n", len(yamls), path)
        // Need to create key (filename minus underscore).
        key := filename[1:]
        m[key] = make(map[string]interface{})
        yerr := yaml.Unmarshal(yamls, m[key])
        btu.CheckError(yerr)
      } else {
        // Recurse into directory.
        btu.Info("Recursing into %s...\n", path)
        m[filename] = make(map[string]interface{})
        loadDataDir(path, m[filename].(map[string]interface{}))
      }
    }
  }
}

func loadGlobalData(dataDir string) map[string]interface{} {
  data := make(map[string]interface{})
  if btu.DirExists(dataDir) {
    err := filepath.Walk(dataDir, func(path string, fileInfo fs.FileInfo, err error) error {
      // Ignore non-yaml files.
      if filepath.Ext(path) != ".yaml" {
        return nil
      }
      relativePath := getRelativePath(path, dataDir)
      btu.Debug("global data file path is %s, relativePath is %s\n", path, relativePath)
      // Create nested data if the file is in a subdirectory.
      parts := strings.Split(relativePath, string(os.PathSeparator))
      m := data
      for j := 0; j < (len(parts) - 1); j++ {
        btu.Debug("Need to create a submap called %s...\n", parts[j])
        m[parts[j]] = make(map[string]interface{})
        m = m[parts[j]].(map[string]interface{})
      }
      b := btu.ReadFileB(path)
      // Get base name of the file, use that as the key in data, unless file begins with
      // an underscore.
      var yerr error
      if strings.HasPrefix(fileInfo.Name(), "_") {
        yerr = yaml.Unmarshal(b, m)
      } else {
        key := strings.TrimSuffix(fileInfo.Name(), ".yaml")
        m[key] = make(map[string]interface{})
        yerr = yaml.Unmarshal(b, m[key])
      }
      btu.CheckError(yerr)
  		return nil
    })
    btu.CheckError(err)
  }
  btu.Debug("global data: %+v\n", data)
  return data
}

func getRelativePath(path, prefix string) string {
  prefixLen := len(prefix)
  s := path[prefixLen:]
  if strings.HasPrefix(s, string(os.PathSeparator)) {
    return s[1:]
  }
  return s
}
