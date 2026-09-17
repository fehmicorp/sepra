package main

import (
	"runtime"
)

var Target = struct {
	OS       string
	Arch     string
	Hostname string
}{
	OS:       runtime.GOOS,
	Arch:     runtime.GOARCH,
	Hostname: "",
}

type AppConfig struct {
	AppName         string
	Description     string
	Icon            string
	Tagline         string
	Title           string
	Width           int
	Height          int
	Min             bool
	Max             bool
	Quit            bool
	Version         string
	Domain          string
	InstallationDir string `json:"installDir"`
}

var Conf = AppConfig{
	AppName:         "SEPRA",
	Description:     "Fehmi Agent Installer is a web-based application that allows users to easily install and manage Fehmi agents on their systems. It provides a user-friendly interface for configuring and deploying agents, making it simple for users to monitor and control their cloud infrastructure.",
	Icon:            "assets/logo.png",
	Tagline:         "Endpoint Security Installer",
	Title:           "Fehmi Corporation",
	Width:           850,
	Height:          600,
	Min:             true,
	Max:             false,
	Quit:            false,
	Version:         "v1.0.1",
	Domain:          "fehmicorp.in",
	InstallationDir: `%ProgramFiles%\fehmi\agent`,
}

type SystemPayload struct {
	Hostname    string `json:"hostname"`
	TargetOS    string `json:"targetOs"`
	TargetArch  string `json:"targetArch"`
	TargetBuild string `json:"targetBuild"`
	AgentVer    string `json:"agentVersion"`
}
