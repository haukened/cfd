package main

import "log"

var debug bool

func Debugf(format string, a ...interface{}) {
	if debug {
		log.Printf(format, a...)
	}
}
