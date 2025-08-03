package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
	"github.com/yulog/ytmusic-extension/tools"
)

type Build mg.Namespace

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

var (
	BIN                 string = "native-app"
	VERSION             string = getVersion()
	CURRENT_REVISION, _        = sh.Output("git", "rev-parse", "--short", "HEAD")
	BUILD_LDFLAGS       string = "-s -w -X main.revision=" + CURRENT_REVISION
	BUILD_TARGET        string = "."

	DEST = filepath.Join("dist", "contents")
)

// func init() {
// 	VERSION = getVersion()
// 	CURRENT_REVISION, _ = sh.Output("git", "rev-parse", "--short", "HEAD")
// }

func getVersion() string {
	_, err := exec.LookPath("gobump")
	if err != nil {
		fmt.Println("installing gobump")
		sh.Run("go", "install", "github.com/x-motemen/gobump/cmd/gobump@latest")
	}
	v, _ := sh.Output("gobump", "show", "-r", BUILD_TARGET)
	return v
}

func bin() string {
	bin := BIN
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	return bin
}

// A build step that requires additional params, or platform specific steps for example
func (Build) App() error {
	// mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags="+BUILD_LDFLAGS, "-o", bin(), BUILD_TARGET)
	return cmd.Run()
}

func (Build) Release() error {
	defer os.RemoveAll(DEST)
	mg.Deps(Collect)
	fmt.Println("Packaging...")
	bin := "ytmusic_extension_setup"
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags="+BUILD_LDFLAGS+" -H=windowsgui", "-o", bin, "./installer")
	return cmd.Run()
}

// Build the Go native application, and setup extension for develop
func (Build) Dev() error {
	log.Println("Starting installation for Windows...")

	// 1. Get current directory (where install.go is located)
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
		return err
	}
	log.Printf("Working directory: %s", currentDir)

	log.Println("Building native application...")
	bin := bin()
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags="+BUILD_LDFLAGS, "-o", bin, BUILD_TARGET, "-tags", "debuglog")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to build Go native application: %v", err)
		return err
	}
	log.Println("Native app build successful.")

	extensionID := "pckomeadlamhidfmcmhkfdehlfblaeah" // Replace with your actual extension ID later
	manifestPath := tools.CreateManifest(currentDir, bin, extensionID)
	tools.Register(manifestPath)
	return nil
}

func Collect() error {
	mg.Deps(Build.App)
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	info, err := os.Stat(currentDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(DEST, info.Mode().Perm())
	if err != nil {
		return err
	}

	err = os.CopyFS(filepath.Join(DEST, "chrome-extension"), os.DirFS("chrome-extension"))
	if err != nil {
		return err
	}

	w, err := os.Create(filepath.Join(DEST, bin()))
	if err != nil {
		return err
	}
	defer w.Close()

	r, err := os.Open(bin())
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build.App)
	fmt.Println("Installing...")
	cmd := exec.Command("go", "install", "-ldflags="+BUILD_LDFLAGS, BUILD_TARGET)
	return cmd.Run()
}

// Manage your deps, or running package managers.
// func InstallDeps() error {
// 	fmt.Println("Installing Deps...")
// 	cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
// 	return cmd.Run()
// }

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("goxz")
	bin := BIN
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	os.RemoveAll(bin)
}

func ShowVersion() {
	fmt.Println(getVersion())
}

func Credits() {
	_, err := exec.LookPath("gocredits")
	if err != nil {
		fmt.Println("installing gocredits")
		sh.Run("go", "install", "github.com/Songmu/gocredits")
	}
	s, _ := sh.Output("gocredits", ".")
	f, _ := os.Create("CREDITS")
	f.WriteString(s)
	defer f.Close()
}

func Cross() {
	_, err := exec.LookPath("goxz")
	if err != nil {
		fmt.Println("installing goxz")
		sh.Run("go", "install", "github.com/Songmu/goxz/cmd/goxz@latest")
	}
	sh.Run("goxz", "-n", BIN, "-pv=v"+VERSION, BUILD_TARGET)
	// sh.Run("goxz", "-n", BIN+"_lite", "-pv=v"+VERSION, "-build-tags=lite", BUILD_TARGET)
}

func Bump() {
	_, err := exec.LookPath("gobump")
	if err != nil {
		fmt.Println("installing gobump")
		sh.Run("go", "install", "github.com/x-motemen/gobump/cmd/gobump@latest")
	}
	sh.Run("gobump", "up", "-w", BUILD_TARGET)
}

func Upload() {
	_, err := exec.LookPath("ghr")
	if err != nil {
		fmt.Println("installing ghr")
		sh.Run("go", "install", "github.com/tcnksm/ghr@latest")
	}
	sh.Run("ghr", "-draft", "v"+VERSION, "goxz")
}
