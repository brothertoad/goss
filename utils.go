package main

import (
  "log"
  "os"
  "os/exec"
  "path/filepath"
)

func fileExists(path string) bool {
  fileInfo, err := os.Stat(path)
  if err != nil {
    return false
  }
  if !fileInfo.Mode().IsRegular() {
    log.Fatal("%s exists, but is not a file\n", path)
  }
  return true
}

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

func executeCommand(cmd interface{}) {
  // cmd can be either a string or a slice of strings
  var command *exec.Cmd
  switch cmd.(type) {
  case string:
    // command = exec.Command(fmt.Sprintf("%v", cmd))
    command = exec.Command(cmd.(string))
  case []string:
    cmds := cmd.([]string)
    command = exec.Command(cmds[0], cmds[1:]...)
  default:
    log.Fatalf("don't know how to handle command type\n")
  }
  err := command.Run()
  checkError(err)
}

func executeCommands(inputDir string, cmds []string) {
  if !dirExists(inputDir) {
    return
  }
  command := exec.Command(cmds[0], cmds[1:]...)
  err := command.Run()
  checkError(err)
}

func includeAction(path string) string {
  b, err := os.ReadFile(path)
  checkError(err)
  return string(b)
}

func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
