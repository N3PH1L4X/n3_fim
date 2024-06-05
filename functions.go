package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Crawls recursively every single directory under each specified directory in the confing.yaml file and returns a list of directories found and files found
func GetFilesAndDirectories() ([]string, []string) {

	var baseline_filename = config.Monitoring.BaselineFilename

	var recursive_directories []string
	var files_to_baseline []string

	for _, root_directory := range config.Monitoring.DirectoriesToMonitor {

		err := filepath.WalkDir(root_directory, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if d.IsDir() {
				recursive_directories = append(recursive_directories, path)
			} else {
				if filepath.Base(path) != baseline_filename {
					files_to_baseline = append(files_to_baseline, path)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error walking the path:", err)
		}
	}

	return recursive_directories, files_to_baseline
}

// Calculates specified hash of a given file and returns it with the hashing algorithm used like: return contents_hash, selected_hash
func GetHash(filename string) (string, string) {

	var contents_hash string
	var selected_hash = config.Monitoring.HashingAlgorithm

	body, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	switch selected_hash {
	case "md5":
		hash := md5.Sum([]byte(body))
		contents_hash = hex.EncodeToString(hash[:])
	case "sha256":
		hash := sha256.Sum256([]byte(body))
		contents_hash = hex.EncodeToString(hash[:])
	case "crc32":
		hash := crc32.ChecksumIEEE([]byte(body))
		contents_hash = hex.EncodeToString([]byte{
			byte(hash >> 24),
			byte(hash >> 16),
			byte(hash >> 8),
			byte(hash),
		})
	default:
		panic("error selecting hashing algorithm, check your config.yaml file")
	}

	return contents_hash, selected_hash
}

// Saves calculated hash to new baseline.csv file
func SaveToBaselineFile(hashing_algorithm string, resulting_hash string, hashed_file string) {

	var baseline_filename = config.Monitoring.BaselineFilename
	var text_to_append string = hashing_algorithm + "," + resulting_hash + "," + hashed_file + "\n"

	f, err := os.OpenFile(baseline_filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(text_to_append); err != nil {
		panic(err)
	}
}

// Deletes old baseline file
func DeleteOldBaseline() {
	var baseline_filename = config.Monitoring.BaselineFilename

	err := os.Remove(baseline_filename)
	if err != nil {
		log.Fatal(err)
	}
}

// Creates new baseline file
func CreateNewBaseline() {
	var baseline_filename = config.Monitoring.BaselineFilename

	_, err := os.Create(baseline_filename)
	if err != nil {
		log.Fatal(err)
	}
}

// Writes CSV file headers as: algorithm,hash,file
func WriteFileHeaders() {
	var baseline_filename = config.Monitoring.BaselineFilename

	f, err := os.OpenFile(baseline_filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString("algorithm,hash,file\n")
	if err != nil {
		panic(err)
	}
}

// Checks if a file exists in the same directory as the binary by using a filename, returns boolean value
func CheckFileExists(filename_to_check string) bool {

	_, error := os.Stat(filename_to_check)
	return !errors.Is(error, os.ErrNotExist)
}

// Shows the resulting baseline after creating it when running: ./n3_fim baseline
func ShowBaseline() {
	var baseline_filename = config.Monitoring.BaselineFilename

	dat, err := os.ReadFile(baseline_filename)
	if err != nil {
		panic(err)
	}

	fmt.Print(Green + "New baseline defined successfully:\n\n" + Reset)
	fmt.Print(string(dat))
}

// If there is no config file in the same directory as the binary then this will create a new default config.yaml file. You then have to edit it as you need
func CreateDefaultConfig() {
	f, err := os.OpenFile("config.yaml", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(DefaultConfig); err != nil {
		panic(err)
	}
}

// Read baseline file and return array of lines read
func ReadBaseLine() []string {

	filePath := config.Monitoring.BaselineFilename
	readFile, err := os.Open(filePath)

	var read_lines []string

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	for _, line := range fileLines {

		if line != "algorithm,hash,file" {
			read_lines = append(read_lines, line)
		}
	}

	return read_lines
}

// Helpful function to debug code and print arrays in a more human readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
	return string(s)
}

func AddToLog(file_flagged string, flag_reason string) {
	time_when_flagged := time.Now()
	var log_filename string = config.Monitoring.LogFilename

	var formatted_time string = time_when_flagged.Format("2006-01-02 15:04:05")
	var absolute_path string = ResolveAbsolutePath(file_flagged)

	f, err := os.OpenFile(log_filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(formatted_time + " : " + flag_reason + " : " + absolute_path + "\n"); err != nil {
		panic(err)
	}
}

func ResolveAbsolutePath(relative_path string) string {
	abs, err := filepath.Abs(relative_path)
	if err != nil {
		log.Fatal(err)
	}
	return abs
}
