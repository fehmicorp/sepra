package main

import (
	"strings"
)

type Conf struct {
	AppName string
}

var Config = []Conf{
	{AppName: "Auth"},
	{AppName: "Home"},
}

// GetConfig now returns a primitive string instead of a struct pointer
func GetConfig(tag string) string {
	for i := range Config {
		if strings.EqualFold(Config[i].AppName, tag) {
			return Config[i].AppName
		}
	}
	return ""
}
