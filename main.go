package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/systray"
)

const powerSupplyPath = "/sys/class/power_supply/"

// Read the battery level and status from uevent
func readBatteryInfo(devicePath string) (string, bool, error) {
	ueventPath := filepath.Join(devicePath, "uevent")
	file, err := os.Open(ueventPath)
	if err != nil {
		return "", false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var level string
	var charging bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "POWER_SUPPLY_CAPACITY_LEVEL=") {
			level = strings.TrimPrefix(line, "POWER_SUPPLY_CAPACITY_LEVEL=")
		} else if strings.HasPrefix(line, "POWER_SUPPLY_STATUS=") {
			status := strings.TrimPrefix(line, "POWER_SUPPLY_STATUS=")
			charging = strings.ToLower(status) == "charging"
		}
	}

	if level == "" {
		return "", false, fmt.Errorf("no capacity level found")
	}

	return level, charging, nil
}

// Find Xbox battery device
func findXboxBatteryDevice() (string, error) {
	entries, err := os.ReadDir(powerSupplyPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		devicePath := filepath.Join(powerSupplyPath, entry.Name())
		typePath := filepath.Join(devicePath, "type")
		data, err := os.ReadFile(typePath)
		if err != nil {
			continue
		}
		deviceType := strings.TrimSpace(string(data))
		if deviceType != "Battery" {
			continue
		}

		ueventPath := filepath.Join(devicePath, "uevent")
		file, err := os.Open(ueventPath)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "POWER_SUPPLY_MODEL_NAME=") {
				model := strings.TrimPrefix(line, "POWER_SUPPLY_MODEL_NAME=")
				if strings.Contains(strings.ToLower(model), "xbox") && strings.Contains(strings.ToLower(model), "controller") {
					file.Close()
					return devicePath, nil
				}
			}
		}
		file.Close()
	}

	return "", fmt.Errorf("no Xbox controller device found")
}

func detectDarkMode() bool {
	configPaths := []string{
		filepath.Join(os.Getenv("HOME"), ".config", "gtk-4.0", "settings.ini"),
		filepath.Join(os.Getenv("HOME"), ".config", "gtk-3.0", "settings.ini"),
	}

	for _, path := range configPaths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "gtk-application-prefer-dark-theme") && strings.HasSuffix(line, "=1") {
				file.Close()
				return true
			}
		}
		file.Close()
	}
	return false
}

func loadIcons(theme string) map[string][]byte {
	prefix := filepath.Join("icons", theme)
	icons := map[string][]byte{}
	files := []string{
		"battery_high.png",
		"battery_medium.png",
		"battery_low.png",
		"battery_critical.png",
		"battery_charging.png",
		"battery_unknown.png",
		"battery_disconnected.png",
	}
	for _, f := range files {
		path := filepath.Join(prefix, f)
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Failed to load icon %s: %v\n", path, err)
			continue
		}
		key := strings.TrimSuffix(strings.TrimPrefix(f, "battery_"), ".png")
		icons[key] = data
	}
	return icons
}

func updateTrayTooltip(batteryLevel string, charging bool, iconFiles map[string][]byte) {
	// Show battery level (full, high, medium,...)
	status := strings.ToLower(batteryLevel)
	if status == "full" {
		status = "high" // icon full = icon high
	}

	tooltip := fmt.Sprintf("Xbox Controller Battery: %s", batteryLevel)
	if charging {
		tooltip += " (Charging)"
	}
	systray.SetTitle(tooltip)

	iconKey := status
	if charging {
		iconKey = "charging"
	}
	if icon, ok := iconFiles[iconKey]; ok {
		systray.SetIcon(icon)
	} else {
		systray.SetIcon(iconFiles["unknown"])
	}
}

func onReady() {
	theme := "light"
	if detectDarkMode() {
		theme = "dark"
	}
	iconFiles := loadIcons(theme)
	systray.SetTitle("Xbox Controller Battery")

	// Create menu for Exit
	mExit := systray.AddMenuItem("Exit", "Exit the application")

	// Main loop
	go func() {
		for {
			devicePath, err := findXboxBatteryDevice()
			if err != nil {
				fmt.Println("Xbox controller not found.")
				systray.SetTitle("Xbox Controller Battery: Disconnected")
				if icon, ok := iconFiles["disconnected"]; ok {
					systray.SetIcon(icon)
				} else {
					systray.SetIcon(iconFiles["unknown"])
				}
			} else {
				level, charging, err := readBatteryInfo(devicePath)
				if err != nil {
					fmt.Println("Error reading battery info:", err)
					updateTrayTooltip("unknown", false, iconFiles)
				} else {
					updateTrayTooltip(level, charging, iconFiles)
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		for range mExit.ClickedCh {
			systray.Quit()
			return
		}
	}()
}

func main() {
	systray.Run(onReady, func() {})
}
