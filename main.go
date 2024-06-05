package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

func main() {

	args := os.Args

	// Arguments validation
	if len(args) > 2 {
		panic(Yellow + "too many arguments" + Reset + Usage)
	}
	if len(args) < 2 {
		panic(Yellow + "missing argument" + Reset + Usage)
	}

	if !CheckFileExists("config.yaml") {
		fmt.Println(Yellow + "config.yaml file does not exists, creating new one..." + Reset)
		CreateDefaultConfig()
		fmt.Println("Please check your new configs and run again the binary")
		os.Exit(0)
	}

	switch args[1] {
	// Establish new baseline
	case "baseline":
		baseline()
	// Begin monitoring established baseline
	case "monitor":
		monitor()
	default:
		panic(Yellow + "invalid argument" + Reset + Usage)
	}

}

// Create a new baseline. Inspects defined directories in config.yaml file, calculates hashes of found files and stores them in a new baseline.csv. This overwrites already existing baseline file
func baseline() {

	if CheckFileExists("baseline.csv") {
		DeleteOldBaseline()
	}

	CreateNewBaseline()
	WriteFileHeaders()

	_, files := GetFilesAndDirectories()

	for _, file := range files {
		result_hash, selected_hash := GetHash(file)
		SaveToBaselineFile(selected_hash, result_hash, file)
	}

	ShowBaseline()

}

func monitor() {

	var read_lines []string
	var mapped_files_from_baseline []map[string]string
	mapped_file_from_right_now := make(map[string]string)
	var all_paths_in_baseline []string

	var stringSlice []string

	if !CheckFileExists("baseline.csv") {
		panic(Red + "Baseline file not found. Please check usage to create one" + Reset + Usage)
	}

	// Get all files, hashes and algos from already defined baseline
	read_lines = ReadBaseLine()
	for _, line := range read_lines {
		stringSlice = strings.Split(line, ",")

		mapped_files_from_baseline = append(
			mapped_files_from_baseline,
			map[string]string{
				"hashing_algorithm": stringSlice[0],
				"value":             stringSlice[1],
				"path":              stringSlice[2],
			},
		)

		all_paths_in_baseline = append(all_paths_in_baseline, stringSlice[2])
	}

	// Infinite loop to begin monitoring
	for {

		// Get actual files and directories
		_, files := GetFilesAndDirectories()
		for _, file := range files {
			result_hash, selected_hash := GetHash(file)

			mapped_file_from_right_now["hashing_algorithm"] = selected_hash
			mapped_file_from_right_now["value"] = result_hash
			mapped_file_from_right_now["path"] = file

			if slices.Contains(all_paths_in_baseline, file) {
				for _, mapped_file_from_baseline := range mapped_files_from_baseline {

					if mapped_file_from_right_now["path"] == mapped_file_from_baseline["path"] {

						if mapped_file_from_right_now["value"] == mapped_file_from_baseline["value"] {
							break
						} else {
							fmt.Println(Yellow + "ALERT! File modified: " + Reset + file)
							AddToLog(file, "FILE MODIFIED")
							break
						}
					}

				}
			} else {
				fmt.Println(Red + "ALERT! New file was created: " + Reset + file)
				AddToLog(file, "NEW FILE CREATED")
			}

		}

		time.Sleep(config.Monitoring.CheckingFrequency * time.Second)

	}

}
