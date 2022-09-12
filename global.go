package main

import (
  "io/fs"
  "os"
  "path/filepath"
  "strings"
  "gopkg.in/yaml.v3"
  "github.com/brothertoad/btu"
)

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
