package main

import (
  "strings"
  "path/filepath"
  "io/fs"
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
      // TASK: Need to create nested data if files are in subdirectories.  Perhaps not bother
      // to support that, in which case we should skip those files with a warning.
      b := btu.ReadFileB(path)
      // Get base name of the file, use that as the key in data, unless file begins with
      // an underscore.
      var yerr error
      if strings.HasPrefix(fileInfo.Name(), "_") {
        yerr = yaml.Unmarshal(b, data)
      } else {
        key := strings.TrimSuffix(fileInfo.Name(), ".yaml")
        data[key] = make(map[string]interface{})
        yerr = yaml.Unmarshal(b, data[key])
      }
      btu.CheckError(yerr)
  		return nil
    })
    btu.CheckError(err)
  }
  return data
}
