package main

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
  "github.com/adrg/frontmatter"
  "gopkg.in/yaml.v3"
  "github.com/brothertoad/btu"
  "github.com/brothertoad/goss/gossutil"
)

const DATA_DIR = "data"
const PAGES_DIR = "pages"
const SRC_DIR = "pages/models"
const J2_SUFFIX = "html.j2"
const YAML_SUFFIX = "yaml"
const PER_PAGE_DIR = "perPage"

type KitType map[string]string

var globalData map[string]interface{}
var kits map[string]KitType

func main() {
  globalData = gossutil.LoadGlobalData(DATA_DIR)
  kits = globalData["kits"].(map[string]KitType)
  _ = filepath.Walk(SRC_DIR, func(path string, info os.FileInfo, err error) error {
    if strings.HasSuffix(path, J2_SUFFIX) {
      relativePath := path[(len(PAGES_DIR) + 1):]
      dataRelativePath := filepath.Join(PER_PAGE_DIR, (strings.TrimSuffix(relativePath, J2_SUFFIX) + YAML_SUFFIX))
      kitKey := getKitKey(path)
      kit := kits[kitKey]
      fmt.Printf("Walking %s, relativePath is %s, dataRelativePath is %s, kitKey is %s\n", path, relativePath, dataRelativePath, kitKey)
      pageData := createPageData(kit)
      writePageData(dataRelativePath, pageData)
    }
    return nil
  })
}

func createPageData(kit KitType) map[string]interface{} {
  data := make(map[string]interface{})
  data["title"] = kit["name"]
  return data
}

func writePageData(path string, data map[string]interface{}) {
  b, err := yaml.Marshal(data)
  btu.CheckError(err)
  btu.CreateDirForFile(path)
  err = os.WriteFile(path, b, 0644)
  btu.CheckError(err)
}

func getKitKey(path string) string {
  b, err := ioutil.ReadFile(path)
  btu.CheckError(err)
  fm := make(map[string]interface{})
  _, err = frontmatter.Parse(bytes.NewReader(b), &fm)
  btu.CheckError(err)
  return fm["kit"].(string)
}
