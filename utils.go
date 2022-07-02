package main

import (
  "os"
  "path/filepath"
  "github.com/brothertoad/btu"
)

/*
func _fileExists(path string) bool {
  fileInfo, err := os.Stat(path)
  if err != nil {
    return false
  }
  if !fileInfo.Mode().IsRegular() {
    log.Fatal("%s exists, but is not a file\n", path)
  }
  return true
}

func _dirExists(dir string) bool {
  fileInfo, err := os.Stat(dir)
  if err != nil {
    return false
  }
  if !fileInfo.IsDir() {
    log.Fatal("%s exists, but is not a directory\n", dir)
  }
  return true
}

func _dirMustExist(dir string) {
  if !_dirExists(dir) {
    log.Fatal("%s does not exist\n", dir)
  }
}
*/

func createDir(dir string) {
  err := os.MkdirAll(dir, 0755)
  btu.CheckError(err)
}

func createDirForFile(path string) {
  dir, _ := filepath.Split(path)
  createDir(dir)
}

func includeAction(path string) string {
  b, err := os.ReadFile(path)
  btu.CheckError(err)
  return string(b)
}

/*
func _checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
*/
