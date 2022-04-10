package main

import (
  "os/exec"
)

func copyStaticFiles(staticDir string, outputDir string, cmds []string) {
  if !dirExists(staticDir) {
    return
  }
  command := exec.Command(cmds[0], cmds[1:]...)
  err := command.Run()
  checkError(err)
}
