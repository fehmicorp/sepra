package main

import (
	"fmt"
	"plugin"
)

func main() {
	// RunWebServer()
	RunPackage("config")
	select {}
}

func RunPackage(pkgName string) {
	path := fmt.Sprintf("%s%s.so", resDir, pkgName)
	p, err := plugin.Open(path)
	if err != nil {
		fmt.Printf("Failed to open plugin: %v\n", err)
		return
	}

	symbol, err := p.Lookup("GetConfig")
	if err != nil {
		fmt.Printf("Failed to find GetConfig symbol: %v\n", err)
		return
	}

	// Match the plugin's exact signature using primitive strings
	getConfigFunc, ok := symbol.(func(string) string)
	if !ok {
		fmt.Println("Plugin function signature mismatch")
		return
	}

	result := getConfigFunc("Auth")
	if result != "" {
		fmt.Printf("Successfully received data from dynamic .so library! Found App: %s\n", result)
	} else {
		fmt.Println("App configuration not found.")
	}
}
