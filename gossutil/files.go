package gossutil;

import (
  "path/filepath"
  "sort"
  "github.com/brothertoad/btu"
)

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
