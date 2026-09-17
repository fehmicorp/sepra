package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fehmicorp/pkg/v1/utils/os/win"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) GetConfig() AppConfig {
	return Conf
}
func (a *App) WindowMinimize() {
	runtime.WindowMinimise(a.ctx)
}
func (a *App) WindowToggleMaximize() {
	runtime.WindowToggleMaximise(a.ctx)
}
func (a *App) WindowClose() {
	runtime.Quit(a.ctx)
}

func getDetailedOSVersion() string {
	if Target.OS != "windows" {
		return Target.OS
	}

	// Open the CurrentVersion registry key
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return "Win"
	}
	defer k.Close()

	productName, _, err := k.GetStringValue("ProductName")
	if err != nil || productName == "" {
		productName = "Win"
	}

	// 1. Convert "Windows 10 Pro" -> "Win 10 Pro"
	if len(productName) >= 7 && productName[:7] == "Windows" {
		productName = "Win" + productName[7:]
	}

	// 2. Fix Windows 11 being misreported as Windows 10 in registry
	buildNumber, _, err := k.GetStringValue("CurrentBuildNumber")
	if err == nil {
		var build int
		fmt.Sscanf(buildNumber, "%d", &build)
		if build >= 22000 {
			// Windows 11 builds start at 22000
			if len(productName) >= 6 && productName[:6] == "Win 10" {
				productName = "Win 11" + productName[6:]
			}
		}
	}

	return productName
}

func getArch() string {
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.QUERY_VALUE,
	)
	var ver string
	if err == nil {
		defer k.Close()
		ver, _, err = k.GetStringValue("DisplayVersion")
		if err != nil || ver == "" {
			ver, _, _ = k.GetStringValue("ReleaseId")
		}
	}
	var arch string
	switch Target.Arch {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "ARM64"
	case "386":
		arch = "x86"
	default:
		arch = Target.Arch
	}
	if ver == "" {
		return arch
	}
	return fmt.Sprintf("%s | %s", ver, arch)
}

func getTargetBuild() string {
	if Target.OS != "windows" {
		return ""
	}

	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return ""
	}
	defer k.Close()

	// Reads main build number (e.g. "22631" or "19045")
	buildNumber, _, err := k.GetStringValue("CurrentBuildNumber")
	if err != nil || buildNumber == "" {
		buildNumber, _, _ = k.GetStringValue("CurrentBuild")
	}

	// Reads minor revision / UBR (e.g. 4890)
	ubr, _, err := k.GetIntegerValue("UBR")
	if err == nil && ubr > 0 {
		return fmt.Sprintf("%s.%d", buildNumber, ubr)
	}

	return buildNumber
}

func (a *App) GetSystemInfo() SystemPayload {
	hostname, _ := os.Hostname()

	return SystemPayload{
		Hostname:    hostname,
		TargetOS:    getDetailedOSVersion(),
		TargetArch:  getArch(),
		TargetBuild: getTargetBuild(),
		AgentVer:    Conf.Version,
	}
}

type SysPayload struct {
	IsAdmin    bool   `json:"isAdmin"`
	TargetArch string `json:"targetArch"`
}

func (a *App) GetSysInfo() SysPayload {
	return SysPayload{
		IsAdmin:    win.IsAdmin(),
		TargetArch: Target.Arch,
	}
}

func (a *App) RelaunchAsAdmin() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	verbPtr, err := windows.UTF16PtrFromString("runas")
	if err != nil {
		return err
	}

	exePtr, err := windows.UTF16PtrFromString(exe)
	if err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	cwdPtr, _ := windows.UTF16PtrFromString(cwd)

	var showCmd int32 = windows.SW_NORMAL

	err = windows.ShellExecute(0, verbPtr, exePtr, nil, cwdPtr, showCmd)
	if err != nil {
		return err
	}

	// Terminate the current non-admin instance
	os.Exit(0)
	return nil
}

type ServerConfig struct {
	IsPrivate bool         `json:"isPrivate"`
	FQDN      string       `json:"fqdn,omitempty"`
	Port      int          `json:"port,omitempty"`
	APIKey    string       `json:"apiKey,omitempty"`
	License   *LicenseInfo `json:"license,omitempty"`
}

type LicenseInfo struct {
	Key       string `json:"key"`
	Issuer    string `json:"issuer"`
	ExpiresAt int64  `json:"expiresAt"`
}
