package main

import (
  "os"
  "github.com/brothertoad/btu"
)

func includeAction(path string) string {
  b, err := os.ReadFile(path)
  btu.CheckError(err)
  return string(b)
}
