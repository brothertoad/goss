package main

import (
  "log"
  "os/exec"
  "reflect"
)

func executeCommand(cmd interface{}) {
  // cmd can be either a string or a slice of strings
  // if it is a slice, break out the first so we can check
  var exe string
  var args []string
  switch cmd.(type) {
  case string:
    // command = exec.Command(fmt.Sprintf("%v", cmd))
    exe = cmd.(string)
    args = make([]string, 0)
  case []string:
    cmds := cmd.([]string)
    exe = cmds[0]
    args = cmds[1:]
  default:
    log.Fatalf("don't know how to handle command type %s\n", reflect.TypeOf(cmd))
  }
  if exe != "" {
    command := exec.Command(exe, args...)
    err := command.Run()
    checkError(err)
  }
}
