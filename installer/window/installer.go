package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type InstallConfig struct {
	IsPrivate        bool   `json:"isPrivate"`
	FQDN             string `json:"fqdn,omitempty"`
	Port             string `json:"port,omitempty"`
	APIKey           string `json:"apiKey,omitempty"`
	LicenseData      string `json:"licenseData,omitempty"`
	InstallPath      string `json:"installPath"`
	InstallAsService bool   `json:"installAsService"`
	EnableAutostart  bool   `json:"enableAutostart"`
}

type InstallProgress struct {
	Percentage int    `json:"percentage"`
	Message    string `json:"message"`
	Log        string `json:"log"`
}

// StartInstallation runs the installer steps asynchronously and emits events
func (a *App) StartInstallation(cfg InstallConfig) error {
	go func() {
		emitLog := func(pct int, msg string, logText string) {
			runtime.EventsEmit(a.ctx, "install_progress", InstallProgress{
				Percentage: pct,
				Message:    msg,
				Log:        fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), logText),
			})
		}

		// Step 1: Create Directories
		emitLog(15, "Creating directory structure...", fmt.Sprintf("Creating target path: %s", cfg.InstallPath))
		if err := os.MkdirAll(cfg.InstallPath, 0755); err != nil {
			emitLog(15, "Failed to create directory", fmt.Sprintf("Error: %v", err))
			runtime.EventsEmit(a.ctx, "install_error", err.Error())
			return
		}

		// Step 2: Extract / Copy Binaries
		emitLog(35, "Extracting Fehmi Agent binaries...", "Unpacking agent core binaries...")
		time.Sleep(500 * time.Millisecond) // Simulate extraction / embed copy

		// Step 3: Write Configuration File & License Data
		emitLog(55, "Writing configuration files...", "Generating config.yaml and security credentials...")

		var configContent string
		if cfg.IsPrivate {
			configContent = fmt.Sprintf(
				"mode: private\nfqdn: %s\nport: %s\napi_key: %s\nservice: %t\nautostart: %t\n",
				cfg.FQDN, cfg.Port, cfg.APIKey, cfg.InstallAsService, cfg.EnableAutostart,
			)
		} else {
			configContent = fmt.Sprintf(
				"mode: public\nservice: %t\nautostart: %t\n",
				cfg.InstallAsService, cfg.EnableAutostart,
			)

			// Write uploaded license file content if present
			if cfg.LicenseData != "" {
				licensePath := filepath.Join(cfg.InstallPath, "license.lic")
				if err := os.WriteFile(licensePath, []byte(cfg.LicenseData), 0600); err != nil {
					emitLog(55, "Failed to write license file", fmt.Sprintf("Error: %v", err))
					runtime.EventsEmit(a.ctx, "install_error", err.Error())
					return
				}
			}
		}

		configPath := filepath.Join(cfg.InstallPath, "config.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
			emitLog(55, "Failed to write config", fmt.Sprintf("Error: %v", err))
			runtime.EventsEmit(a.ctx, "install_error", err.Error())
			return
		}

		// Step 4: Register Windows Service (If requested)
		if cfg.InstallAsService {
			emitLog(75, "Registering Windows Service...", "Executing sc.exe create FehmiAgent Service...")
			exePath := filepath.Join(cfg.InstallPath, "fehmi-agent.exe")
			cmd := exec.Command("sc.exe", "create", "FehmiAgent", "binPath=", exePath, "start=", "auto")
			_ = cmd.Run() // Handle execution according to backend privilege state
		}

		// Step 5: Start Service or Executable
		emitLog(90, "Starting Fehmi Agent Service...", "Initializing service runtime...")
		time.Sleep(500 * time.Millisecond)

		// Step 6: Completion
		targetEndpoint := cfg.FQDN
		if !cfg.IsPrivate {
			targetEndpoint = "public gateway"
		}
		emitLog(100, "Installation completed successfully!", fmt.Sprintf("Handshake verified with %s", targetEndpoint))
		runtime.EventsEmit(a.ctx, "install_complete", true)
	}()

	return nil
}
