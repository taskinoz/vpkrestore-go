package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const hashURL = "https://taskinoz.com/titanfall/pc/" // Replace with your URL

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

	fmt.Printf("%s:\nMD5: %s\nSHA1: %s\nSHA256: %s\n", filePath, hex.EncodeToString(md5Hash.Sum(nil)), hex.EncodeToString(sha1Hash.Sum(nil)), hex.EncodeToString(sha256Hash.Sum(nil)))

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

	for _, file := range vpkFiles {
		md5h, sha1h, sha256h, err := computeHashes(file)
		if err != nil {
			fmt.Printf("Error hashing %s: %v\n", file, err)
			continue
		}
		if rHashes, exists := remoteHashes[file]; exists {
			if md5h != rHashes[0] || sha1h != rHashes[1] || sha256h != rHashes[2] {
				fmt.Printf("Hash mismatch for %s. Downloading correct version...\n", file)
				// Assuming the file is hosted at the same URL with the hash file. Adjust as needed.
				err = downloadFile(hashURL+file, file)
				if err != nil {
					fmt.Printf("Failed to download %s: %v\n", file, err)
				}
			}
		}
	}
}
