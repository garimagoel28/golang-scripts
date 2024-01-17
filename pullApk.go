package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Read package_name from user input
	fmt.Print("Enter package_name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	packageName := strings.TrimSpace(scanner.Text())

	// Ask the user for a folder name
	fmt.Print("Enter folder name for saving APKs: ")
	scanner.Scan()
	folderName := strings.TrimSpace(scanner.Text())

	// Create the destination directory based on the provided folder name
	destDirectory := filepath.Join(".", folderName)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(destDirectory, 0755); err != nil {
		fmt.Println("Error creating destination directory:", err)
		return
	}

	// Get the paths of the APKs
	paths, err := getApkPaths(packageName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Pull all the APKs to the specified destination directory
	for _, path := range paths {
		if err := pullApk(path, destDirectory); err != nil {
			fmt.Printf("Error pulling APK from path %s: %v\n", path, err)
		}
	}
}

func getApkPaths(packageName string) ([]string, error) {
	// Run the adb shell pm path command
	cmd := exec.Command("adb", "shell", "pm", "path", packageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute adb shell pm path command: %v", err)
	}

	// Extract paths from the output
	var paths []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "package:") {
			// Extract path from the line
			path := strings.TrimPrefix(line, "package:")
			paths = append(paths, path)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning output: %v", err)
	}

	return paths, nil
}

func pullApk(path, destDirectory string) error {
	// Specify the destination file path
	destFilePath := filepath.Join(destDirectory, filepath.Base(path))

	// Run the adb pull command
	cmd := exec.Command("adb", "pull", path, destFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute adb pull command: %v\n%s", err, output)
	}

	fmt.Printf("APK pulled successfully from path %s to %s\n", path, destFilePath)
	return nil
}
