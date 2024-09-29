package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Directories []string `yaml:"directories"`
	FileTypes   []string `yaml:"file_types"`
	DeleteAfter string   `yaml:"delete_after"`
}

func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func parseCustomDuration(duration string) (time.Duration, error) {
	if strings.HasSuffix(duration, "d") {
		daysStr := strings.TrimSuffix(duration, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, fmt.Errorf("invalid number of days: %v", err)
		}
		// Convert days to hours
		return time.Duration(days) * 24 * time.Hour, nil
	}
	return time.ParseDuration(duration) // Handle standard Go duration formats
}

func checkAndDeleteFiles(config *Config, dryRun bool) {
	deleteAfterDuration, err := parseCustomDuration(config.DeleteAfter)
	if err != nil {
		log.Fatalf("Invalid delete_after duration: %v", err)
	}
	cutoff := time.Now().Add(-deleteAfterDuration)

	for _, dir := range config.Directories {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			ext := filepath.Ext(info.Name())
			if !contains(config.FileTypes, ext) {
				return nil
			}

			if info.ModTime().Before(cutoff) {
				if dryRun {
					fmt.Printf("[DRY RUN] Would delete %s (modified: %v)\n", path, info.ModTime())
				} else {
					fmt.Printf("Deleting %s (modified: %v)\n", path, info.ModTime())
					err := os.Remove(path)
					if err != nil {
						log.Printf("Error deleting file: %v\n", err)
					}
				}
			}
			return nil
		})

		if err != nil {
			log.Printf("Error walking the path %v: %v\n", dir, err)
		}
	}
}

// Utility function to check if a string slice contains a value
func contains(slice []string, item string) bool {
	for _, value := range slice {
		if value == item {
			return true
		}
	}
	return false
}

func main() {
	defaultConfig := "filepruner-config.yml"

	configFile := flag.String("config", "", "Configuration file name")
	dryRunFlag := flag.Bool("dry-run", false, "Enable dry run mode (no files will be deleted)")
	flag.Parse()

	if *configFile == "" {
		*configFile = os.Getenv("FILEPRUNER_CONFIG")
	}

	dryRunEnv := os.Getenv("FILEPRUNER_DRY_RUN")
	dryRun := *dryRunFlag || (dryRunEnv == "true")

	if *configFile == "" {
		*configFile = defaultConfig
	}

	fmt.Printf("Using config file: %s\n", *configFile)
	if dryRun {
		fmt.Println("Dry run mode is enabled. No files will be deleted.")
	}

	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	checkAndDeleteFiles(config, dryRun)
}
