//go:build windows

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const (
	hostName    = "com.google.chrome.example.echo"
	exeName     = "native-app.exe"
	extensionID = "pckomeadlamhidfmcmhkfdehlfblaeah" // Replace with your actual extension ID later
)

// Manifest represents the structure of the native messaging host manifest.
type Manifest struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Path           string   `json:"path"`
	Type           string   `json:"type"`
	AllowedOrigins []string `json:"allowed_origins"`
}

func main() {
	log.Println("Starting installation for Windows...")

	// 1. Get current directory (where install.go is located)
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	log.Printf("Working directory: %s", currentDir)

	// 2. Build the Go native application
	log.Println("Building native application (native-app.exe)...")
	buildCmd := exec.Command("go", "build", "-o", exeName, "main.go", "log_release.go")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Failed to build Go native application: %v", err)
	}
	log.Println("Native app build successful.")

	// 3. Create the manifest file for native messaging
	exePath := filepath.Join(currentDir, exeName)
	manifestPath := filepath.Join(currentDir, hostName+".json")

	manifest := Manifest{
		Name:        hostName,
		Description: "Chrome Native Messaging Example Host",
		Path:        exePath,
		Type:        "stdio",
		AllowedOrigins: []string{
			"chrome-extension://" + extensionID + "/",
		},
	}

	manifestBytes, err := json.MarshalIndent(manifest, "", "    ")
	if err != nil {
		log.Fatalf("Failed to create manifest JSON: %v", err)
	}

	if err := os.WriteFile(manifestPath, manifestBytes, 0644); err != nil {
		log.Fatalf("Failed to write manifest file: %v", err)
	}
	log.Printf("Created manifest file at: %s", manifestPath)

	// 4. Register the native messaging host in the Windows Registry
	log.Println("Registering native messaging host in Windows Registry...")
	keyPath := `Software\Google\Chrome\NativeMessagingHosts\` + hostName
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		log.Fatalf("Failed to create or open registry key: %v", err)
	}
	defer key.Close()

	if err := key.SetStringValue("", manifestPath); err != nil {
		log.Fatalf("Failed to set registry key value: %v", err)
	}
	log.Println("Registry key set successfully.")

	fmt.Println("\nInstallation complete.")
	fmt.Println("You can now load the Chrome extension in developer mode.")
}
