package main

import (
  "fmt"
  "os"
  "path/filepath"
  "strings"
  "github.com/brothertoad/goss/gossutil"
)

const DATA_DIR = "data"
const PAGES_DIR = "pages"
const SRC_DIR = "pages/models"
const J2_SUFFIX = "html.j2"
const YAML_SUFFIX = "yaml"
const PER_PAGE_DIR = "perPage"

var globalData map[string]interface{}

func main() {
  globalData = gossutil.LoadGlobalData(DATA_DIR)
  _ = filepath.Walk(SRC_DIR, func(path string, info os.FileInfo, err error) error {
    if strings.HasSuffix(path, J2_SUFFIX) {
      relativePath := path[(len(PAGES_DIR) + 1):]
      dataRelativePath := filepath.Join(PER_PAGE_DIR, (strings.TrimSuffix(relativePath, J2_SUFFIX) + YAML_SUFFIX))
      fmt.Printf("Walking %s, relativePath is %s, dataRelativePath is %s\n", path, relativePath, dataRelativePath)
    }
    return nil
  })
}
