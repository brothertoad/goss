package main

import (
  "path/filepath"
  "log"
  "os"
)

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
