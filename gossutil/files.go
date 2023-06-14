package gossutil;

import (
  "os"
  "path/filepath"
  "sort"
  "strings"
  "gopkg.in/yaml.v3"
  "github.com/brothertoad/btu"
)

func LoadGlobalData(dataDir string) map[string]interface{} {
  data := make(map[string]interface{})
  loadDataDir(dataDir, data)
  return data
}

func loadDataDir(dataDir string, m map[string]interface{}) {
  pattern := dataDir + string(filepath.Separator) + "*"
  matches, err := filepath.Glob(pattern)
  btu.CheckError(err)
  if matches == nil {
    return
  }
  for _, path := range(matches) {
    filename := getRelativePath(path, dataDir)
    if filepath.Ext(path) == ".yaml" {
      btu.Debug("Found a yaml file, path is %s, filename is %s\n", path, filename)
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
      btu.Debug("%s (filename %s) is a directory\n", path, filename)
      if strings.HasPrefix(filename, "_") {
        // Need to concatenate yaml files in directory.
        btu.Debug("Need to concatenate files in %s\n", path)
        yamls := CatFiles(path, "yaml")
        if yamls == nil {
          btu.Fatal("Got nil when catting files.\n")
        }
        btu.Debug("Found %d bytes in %s\n", len(yamls), path)
        // Need to create key (filename minus underscore).
        key := filename[1:]
        m[key] = make(map[string]interface{})
        yerr := yaml.Unmarshal(yamls, m[key])
        btu.CheckError(yerr)
      } else {
        // Recurse into directory.
        btu.Debug("Recursing into %s...\n", path)
        m[filename] = make(map[string]interface{})
        loadDataDir(path, m[filename].(map[string]interface{}))
      }
    }
  }
}

func getRelativePath(path, prefix string) string {
  prefixLen := len(prefix)
  s := path[prefixLen:]
  if strings.HasPrefix(s, string(os.PathSeparator)) {
    return s[1:]
  }
  return s
}

// Concatenates all the files in the directory that have the
// specified extension.
func CatFiles(dir, ext string) []byte {
  pattern := dir + string(filepath.Separator) + "*." + ext
  matches, err := filepath.Glob(pattern)
  btu.CheckError(err)
  if matches == nil {
    btu.Error("No matches foundfor %s.\n", pattern)
    return make([]byte, 0)
  }
  sort.Strings(matches) // don't know if this is necessary
  b := make([]byte, 0)
  for _, match := range(matches) {
    content := btu.ReadFileB(match)
    b = append(b, content...)
  }
  return b
}
