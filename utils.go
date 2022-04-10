package main

import (
  "path/filepath"
  "log"
  "os"
)

func dirExists(dir string) bool {
  fileInfo, err := os.Stat(dir)
  if err != nil {
    return false
  }
  if !fileInfo.IsDir() {
    log.Fatal("%s exists, but is not a directory\n", dir)
  }
  return true
}

func dirMustExist(dir string) {
  if !dirExists(dir) {
    log.Fatal("%s does not exist\n", dir)
  }
}

func createDir(dir string) {
  err := os.MkdirAll(dir, 0755)
  checkError(err)
}

func createDirForFile(path string) {
  dir, _ := filepath.Split(path)
  createDir(dir)
}

func includeCommand(path string) string {
  b, err := os.ReadFile(path)
  checkError(err)
  return string(b)
}

func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
