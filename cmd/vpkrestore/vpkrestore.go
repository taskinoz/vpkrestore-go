package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var debug bool
var auto bool

// Constants for loading bar
const totalBars = 40 // Total number of bars in the loading bar
const barChar = "█"  // Character to use for the loading bar
var lastBars int     // To track the number of bars from the last update

const hashURL = "https://taskinoz.com/titanfall/pc/" // Replace with your URL

func init() {
	flag.BoolVar(&debug, "d", false, "Enable debug mode")
	flag.BoolVar(&auto, "a", false, "Auto download mismatched files without prompt")
}

func debugPrint(v ...interface{}) {
	if debug {
		fmt.Println(v...)
	}
}

func displayLoadingBar(total, current int) {
	percentage := float64(current) / float64(total)
	bars := int(percentage * totalBars)

	if bars != lastBars { // Only update if the bars changed
		fmt.Printf("\rChecking VPK's [%-40s] ", strings.Repeat(barChar, bars))
		lastBars = bars
	}
}

func getFilesWithExtension(extension string) ([]string, error) {
	var files []string
	entries, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), extension) {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func computeHashes(filePath string) (string, string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()

	if _, err := io.Copy(io.MultiWriter(md5Hash, sha1Hash, sha256Hash), file); err != nil {
		return "", "", "", err
	}

	debugPrint("%s:\nMD5: %s\nSHA1: %s\nSHA256: %s\n", filePath, hex.EncodeToString(md5Hash.Sum(nil)), hex.EncodeToString(sha1Hash.Sum(nil)), hex.EncodeToString(sha256Hash.Sum(nil)))

	return hex.EncodeToString(md5Hash.Sum(nil)), hex.EncodeToString(sha1Hash.Sum(nil)), hex.EncodeToString(sha256Hash.Sum(nil)), nil
}

func downloadHashes() (map[string][3]string, error) {
	resp, err := http.Get(hashURL + "hash.txt")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	hashes := make(map[string][3]string)
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 4 {
			hashes[parts[0]] = [3]string{parts[1], parts[2], parts[3]}
		}
	}
	return hashes, nil
}

func shouldDownload(filename string) bool {
	if auto {
		return true
	}

	fmt.Printf("Found mismatched hashes for: %s\nDo you want to restore this file? (y/n) ", filename)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

func downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	flag.Parse()

	debugPrint("Debug mode enabled")

	vpkFiles, err := getFilesWithExtension(".vpk")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d VPK files\n", len(vpkFiles))

	remoteHashes, err := downloadHashes()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d remote hashes\n", len(remoteHashes))

	for index, file := range vpkFiles {
		displayLoadingBar(len(vpkFiles), index)

		md5h, sha1h, sha256h, err := computeHashes(file)
		if err != nil {
			fmt.Printf("\nError hashing %s: %v\n", file, err)
			continue
		}
		if rHashes, exists := remoteHashes[file]; exists {
			if md5h != rHashes[0] || sha1h != rHashes[1] || sha256h != rHashes[2] {
				if shouldDownload(file) {
					fmt.Printf("\nDownloading correct version of %s...\n", file)
					err = downloadFile(hashURL+file, file)
					if err != nil {
						fmt.Printf("Failed to download %s: %v\n", file, err)
					}
				}
			}
		}
	}

	// To ensure the loading bar is complete at the end
	displayLoadingBar(len(vpkFiles), len(vpkFiles))
	fmt.Println() // Newline after the loading bar
}
