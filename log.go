package main

import (
	logPkg "log"
	"os"
)

var log = logPkg.New(os.Stdout, "", 0)
