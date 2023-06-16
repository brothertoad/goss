package main

import (
  "bytes"
  _ "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "sort"
  "strconv"
  "strings"
  "time"
  "github.com/adrg/frontmatter"
  "gopkg.in/yaml.v3"
  "github.com/brothertoad/btu"
  "github.com/brothertoad/goss/gossutil"
)

const DATA_DIR = "data"
const PAGES_DIR = "pages"
const SRC_DIR = "pages/models"
const J2_SUFFIX = "html.j2"
const YAML_SUFFIX = ".yaml"
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
  // These fields are not needed outside of this app.
  url string
  dataRelativePath string
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
var pageMap map[int]PageDataType

func main() {
  globalData = gossutil.LoadGlobalData(DATA_DIR)
  kitMap = createKitMap(globalData["kits"].(map[string]interface{}))
  pageMap = make(map[int]PageDataType)
  _ = filepath.Walk(SRC_DIR, func(path string, info os.FileInfo, err error) error {
    if strings.HasSuffix(path, J2_SUFFIX) {
      relativePath := path[(len(PAGES_DIR) + 1):]
      base := strings.TrimSuffix(relativePath, "." + J2_SUFFIX)
      dataRelativePath := filepath.Join(PER_PAGE_DIR, base + YAML_SUFFIX)
      kitKey := getKitKey(path)
      kit := kitMap[kitKey]
      // fmt.Printf("Walking %s, relativePath is %s, dataRelativePath is %s, kit is %+v\n", path, relativePath, dataRelativePath, kit)
      pageData := createPageData(kit, relativePath)
      pageData.url = "/" + base + "/"
      pageData.dataRelativePath = dataRelativePath
      pageMap[pageData.Key] = pageData
    }
    return nil
  })
  addPreviousNext()
  for _, pageData := range(pageMap) {
    writePageData(pageData)
  }
}

func addPreviousNext() {
  // Per https://yourbasic.org/golang/how-to-sort-in-go/#bonus-sort-a-map-by-key-or-value
  n := len(pageMap)
  keys := make([]int, 0, n)
  for k := range(pageMap) {
    keys = append(keys, k)
  }
  sort.Ints(keys)
  // Add next to the first page, and previous to the last manually.
  // Have to assign to a copy of the struct, due to Golang.
  // See https://stackoverflow.com/questions/42605337/cannot-assign-to-struct-field-in-a-map
  if entry, ok := pageMap[keys[0]]; ok {
    entry.NextUrl = pageMap[keys[1]].url
    pageMap[keys[0]] = entry
  }
  if entry, ok := pageMap[keys[n-1]]; ok {
    entry.PreviousUrl = pageMap[keys[n-2]].url
    pageMap[keys[n-1]] = entry
  }
  // Now do the ones in between.
  for j := 1; j < (n -1); j++ {
    if entry, ok := pageMap[keys[j]]; ok {
      entry.NextUrl = pageMap[keys[j+1]].url
      entry.PreviousUrl = pageMap[keys[j-1]].url
      pageMap[keys[j]] = entry
    }
  }
}

func createPageData(kit KitType, relativePath string) PageDataType {
  var pageData PageDataType
  pageData.Title = kit.name
  if kit.boxart != "" && kit.boxart != "None" {
    pageData.BoxartUrl = "https://d1dems3vhrlf9r.cloudfront.net/boxart/" + kit.boxart
  }
  if kit.scalematesId != "" {
    pageData.ScalematesUrl = "http://www.scalemates.com/kits/" + kit.scalematesId
  }
  key, err := strconv.Atoi(relativePath[7:11] + relativePath[12:16])
  btu.CheckError(err)
  pageData.Key = key
  month, err := strconv.Atoi(relativePath[12:14])
  btu.CheckError(err)
  pageData.CompletionDate = time.Month(month).String() + ", " + relativePath[7:11]
  return pageData
}

func writePageData(pageData PageDataType) {
  b, err := yaml.Marshal(pageData)
  btu.CheckError(err)
  btu.CreateDirForFile(pageData.dataRelativePath)
  err = os.WriteFile(pageData.dataRelativePath, b, 0644)
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
