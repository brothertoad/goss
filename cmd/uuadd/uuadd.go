package main

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "strconv"
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

type KitType struct {
  name string
  boxart string
  scalematesId string
  brand string
  scale string
  number string
}

type PageDataType struct {
  Title string `yaml:"title"`
  BoxartUrl string `yaml:"boxartUrl"`
  ScalematesUrl string `yaml:"scalematesUrl"`
  CompletionDate string `yaml:"completionDate"`
  PreviousUrl string `yaml:"previousUrl"`
  NextUrl string `yaml:"nextUrl"`
  Key int `yaml:"key"`
}

// Keys for KitType
const KIT_KEY_NAME = "name"
const KIT_KEY_BOXART = "boxart"
const KIT_KEY_SCALEMATES = "scalematesId"
const KIT_KEY_BRAND = "brand"
const KIT_KEY_SCALE = "scale"
const KIT_KEY_NUMBER = "number"

var globalData map[string]interface{}
var kitMap map[string]KitType

func main() {
  globalData = gossutil.LoadGlobalData(DATA_DIR)
  kitMap = createKitMap(globalData["kits"].(map[string]interface{}))
  _ = filepath.Walk(SRC_DIR, func(path string, info os.FileInfo, err error) error {
    if strings.HasSuffix(path, J2_SUFFIX) {
      relativePath := path[(len(PAGES_DIR) + 1):]
      dataRelativePath := filepath.Join(PER_PAGE_DIR, (strings.TrimSuffix(relativePath, J2_SUFFIX) + YAML_SUFFIX))
      kitKey := getKitKey(path)
      kit := kitMap[kitKey]
      fmt.Printf("Walking %s, relativePath is %s, dataRelativePath is %s, kit is %+v\n", path, relativePath, dataRelativePath, kit)
      pageData := createPageData(kit)
      pageData = addBuildInfo(pageData, relativePath)
      fmt.Printf("pageData is %+v\n", pageData)
      writePageData(dataRelativePath, pageData)
    }
    return nil
  })
}

func createPageData(kit KitType) PageDataType {
  var pageData PageDataType
  pageData.Title = kit.name
  if kit.boxart != "" && kit.boxart != "None" {
    pageData.BoxartUrl = "https://d1dems3vhrlf9r.cloudfront.net/boxart/" + kit.boxart
  }
  if kit.scalematesId != "" {
    pageData.ScalematesUrl = "http://www.scalemates.com/kits/" + kit.scalematesId
  }
  return pageData
}

func addBuildInfo(pageData PageDataType, relativePath string) PageDataType {
  key, err := strconv.Atoi(relativePath[7:11] + relativePath[12:16])
  btu.CheckError(err)
  pageData.Key = key
  pageData.CompletionDate = relativePath[7:11]
  return pageData
}

func writePageData(path string, pageData PageDataType) {
  b, err := yaml.Marshal(pageData)
  btu.CheckError(err)
  btu.CreateDirForFile(path)
  err = os.WriteFile(path, b, 0644)
  btu.CheckError(err)
}

func createKitMap(m map[string]interface{}) map[string]KitType {
  kitMap := make(map[string]KitType)
  for k, v := range(m) {
    var kit KitType
    vmap := v.(map[string]interface{})
    kit.name = vmap[KIT_KEY_NAME].(string)
    kit.boxart = vmap[KIT_KEY_BOXART].(string)
    kit.scalematesId = vmap[KIT_KEY_SCALEMATES].(string)
    kit.brand = vmap[KIT_KEY_BRAND].(string)
    kit.scale = vmap[KIT_KEY_SCALE].(string)
    kit.number = vmap[KIT_KEY_NUMBER].(string)
    kitMap[k] = kit
  }
  return kitMap
}

func getKitKey(path string) string {
  b, err := ioutil.ReadFile(path)
  btu.CheckError(err)
  fm := make(map[string]interface{})
  _, err = frontmatter.Parse(bytes.NewReader(b), &fm)
  btu.CheckError(err)
  return fm["kit"].(string)
}
