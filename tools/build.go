//go:build windows

package tools

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
	hostName = "com.github.yulog.ytmusic_notifier"
)

// Manifest represents the structure of the native messaging host manifest.
type Manifest struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Path           string   `json:"path"`
	Type           string   `json:"type"`
	AllowedOrigins []string `json:"allowed_origins"`
}

func Build(bin string) {
	log.Println("Starting installation for Windows...")

	// 2. Build the Go native application
	log.Println("Building native application (native-app.exe)...")
	buildCmd := exec.Command("go", "build", "-o", bin, ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Failed to build Go native application: %v", err)
	}
	log.Println("Native app build successful.")
}

func Dev(bin string) {
	log.Println("Starting installation for Windows...")

	// 2. Build the Go native application
	log.Println("Building native application (native-app.exe)...")
	buildCmd := exec.Command("go", "build", "-o", bin, ".", "-tags", "debuglog")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Failed to build Go native application: %v", err)
	}
	log.Println("Native app build successful.")
}

// Create the manifest file for native messaging
func CreateManifest(currentDir, bin, extensionID string) (manifestPath string) {
	exePath := filepath.Join(currentDir, bin)
	manifestPath = filepath.Join(currentDir, hostName+".json")

	manifest := Manifest{
		Name:        hostName,
		Description: "YouTube Music Notifier",
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

	return
}

// Register the native messaging host in the Windows Registry
func Register(manifestPath string) {
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
