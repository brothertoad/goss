package main

import (
  "strings"
  "path/filepath"
  "io/fs"
  "io/ioutil"
  "gopkg.in/yaml.v3"
)

func loadGlobalData(dataDir string) map[string]interface{} {
  data := make(map[string]interface{})
  err := filepath.Walk(dataDir, func(path string, fileInfo fs.FileInfo, err error) error {
    // Ignore non-yaml files.
    if filepath.Ext(path) != ".yaml" {
      return nil
    }
    b, ferr := ioutil.ReadFile(path)
    checkError(ferr)
    // Get base name of the file, use that as the key in data.
    // TASK: Need to create nested data if files are in subdirectories.
    key := strings.TrimSuffix(fileInfo.Name(), ".yaml")
    data[key] = make(map[string]interface{})
    yerr := yaml.Unmarshal(b, data[key])
    checkError(yerr)
		return nil
  })
  checkError(err)
  return data
}
