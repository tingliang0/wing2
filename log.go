package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	// log.SetPrefix("wing: ")

	file, err := os.OpenFile("errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Debug = log.New(ioutil.Discard,
		"DBG ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout,
		"INF ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout,
		"WAR ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr),
		"ERR ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
