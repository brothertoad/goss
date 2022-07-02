package main

import (
  "log"
  "os"
  "os/exec"
  "reflect"
  "github.com/brothertoad/btu"
)

func executeCommand(cmd interface{}) {
  // cmd can be either a string or a slice of interfaces, which we can typecast to strings
  // if it is a slice, break out the first so we can verify it is on the path
  var exe string
  var args []string
  switch cmd.(type) {
  case string:
    exe = cmd.(string)
    args = make([]string, 0)
  case []interface{}:
    icmds := cmd.([]interface{})
    cmds := make([]string, len(icmds))
    for n, v := range icmds {
      cmds[n] = v.(string)
    }
    exe = cmds[0]
    args = cmds[1:]
  default:
    log.Fatalf("don't know how to handle command type %s\n", reflect.TypeOf(cmd))
  }
  if exe != "" {
    _, err := exec.LookPath(exe)
    btu.CheckError(err)
    command := exec.Command(exe, args...)
    command.Stdout = os.Stdout
    err = command.Run()
    btu.CheckError(err)
  }
}
