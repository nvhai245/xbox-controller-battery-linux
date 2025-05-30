package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/systray"
	"github.com/gen2brain/beeep"
)

const (
	powerSupplyPath = "/sys/class/power_supply/"
	iconPath        = "/usr/share/xbox-controller-battery-linux/icons"
	configDir       = ".config/xbox-controller-battery-linux"
	configFileName  = "battery.conf"
)

var theme string

func loadThemeFromConfig() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return "", err
	}
	file, err := os.Open(filepath.Join(homeDir, configDir, configFileName))
	if err != nil {
		err = fmt.Errorf("error opening file: %v", err)
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			if key == "theme" {
				value := strings.TrimSpace(parts[1])
				return value, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		err = fmt.Errorf("error reading config file: %v", err)
		return "", err
	}

	return "", errors.New("theme is not set")
}

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

func loadIcons(theme string) map[string][]byte {
	prefix := filepath.Join(iconPath, theme)
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

func notifyLowBattery(level string) {
	var (
		message string
		appIcon string
	)
	switch level {
	case "low":
		message = fmt.Sprintf("Battery is %s. Please charge soon.", level)
		appIcon = fmt.Sprintf("%s/%s/%s", iconPath, theme, "battery_low.png")
	case "critical":
		message = fmt.Sprintf("Battery is %s! Your controller will turn off soon!", level)
		appIcon = fmt.Sprintf("%s/%s/%s", iconPath, theme, "battery_critical.png")
	}
	err := beeep.Notify("", message, appIcon)
	if err != nil {
		fmt.Println("Notification error:", err)
	}
}

func refreshIndicator(lastNotifiedLevel *string) {
	devicePath, err := findXboxBatteryDevice()
	if err != nil {
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
			// send notification for low battery
			if !charging {
				currentLevel := strings.ToLower(level)
				if lastNotifiedLevel != nil {
					if (currentLevel == "low" || currentLevel == "critical") && currentLevel != *lastNotifiedLevel {
						notifyLowBattery(currentLevel)
						*lastNotifiedLevel = currentLevel
					}
					if currentLevel != "low" && currentLevel != "critical" {
						*lastNotifiedLevel = ""
					}
				}
			}
		}
	}
}

var iconFiles map[string][]byte

func onReady() {
	t, err := loadThemeFromConfig()
	if err != nil {
		theme = "dark"
		fmt.Println(err)
		fmt.Println("using default theme: dark")
	} else {
		theme = t
	}
	iconFiles = loadIcons(theme)
	systray.SetTitle("Xbox Controller Battery")

	// Create menu entries
	mChangeTheme := systray.AddMenuItem("Change theme", "Change the icon theme")
	mExit := systray.AddMenuItem("Exit", "Exit the application")

	// Main loop
	go func() {
		lastNotifiedLevel := new(string)
		for {
			refreshIndicator(lastNotifiedLevel)
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		for range mExit.ClickedCh {
			systray.Quit()
			return
		}
	}()
	go func() {
		for range mChangeTheme.ClickedCh {
			switch theme {
			case "dark":
				theme = "light"
			case "light":
				theme = "dark"
			default:
			}
			iconFiles = loadIcons(theme)
			refreshIndicator(nil)
			setThemeToConfig(theme)
		}
	}()
}

func setThemeToConfig(newTheme string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	err = os.MkdirAll(filepath.Join(homeDir, configDir), 0755) // Create the directory with appropriate permissions
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	file, err := os.OpenFile(filepath.Join(homeDir, configDir, configFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	config := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	config["theme"] = newTheme

	file, err = os.Create(filepath.Join(homeDir, configDir, configFileName))
	if err != nil {
		fmt.Println("Error creating config file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range config {
		_, err := fmt.Fprintf(writer, "%s=%s\n", key, value)
		if err != nil {
			fmt.Println("Error writing to config file:", err)
			return
		}
	}
	writer.Flush()

	fmt.Println("Updated config file successfully.")
}

func main() {
	systray.Run(onReady, func() {})
}
