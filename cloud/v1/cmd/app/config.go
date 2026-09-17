package main

var (
	Name    = "Sepra Cloud"
	Version = "1.0.0"
	Github  = "https://github.com/fehmicorp/sepra"
	defDir  = "./dist"
	defHost = "0.0.0.0"
	defPort = "8080"
	resDir  = "./pkg/"
)

type Conf struct {
	AppName string
}
