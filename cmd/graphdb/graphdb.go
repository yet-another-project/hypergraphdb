package main

import (
    "os"
    "fmt"
    "flag"

    "github.com/GeertJohan/go.linenoise"
    "github.com/golang/glog"

    "github.com/yet-another-project/hypergraphdb"
)

func main() {
    flag.Parse()
    exitCode := 0

    director := graphdb.NewCommandsDirector()
    for {
        str, err := linenoise.Line("graphdb> ")
        if err != nil {
            if err == linenoise.KillSignalError {
                break
            }
            glog.Errorf("Unexpected error: %s\n", err)
            exitCode = 1
            break
        }
        fields := director.Prepare(str)
        if len(fields) == 0 {
            continue
        }
        status := director.Execute(fields[0], fields[1:])
        if status {
            fmt.Println("OK")
        }
    }

    glog.Flush()

    if exitCode != 0 {
        os.Exit(exitCode)
    }

}
