package main

import (
  "log"
  "os"
)

func createDir(dir string) {
  err := os.MkdirAll(dir, 0755)
  checkError(err)
}

func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
