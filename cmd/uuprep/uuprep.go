package main

import (
  "fmt"
  "os"
  "github.com/brothertoad/goss/gossutil"
)

func main() {
  fmt.Printf("main function in uuprep\n")
  b := gossutil.CatFiles(os.Args[1], "yaml")
  fmt.Printf("Found %d bytes.\n", len(b))
}
