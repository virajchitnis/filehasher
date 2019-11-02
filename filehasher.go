package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Parse command line flags.
	dbFile := flag.String("db", "", "The database file to store the hashes in.")
	flag.Parse()

	// Output
	fmt.Println("Database file to use: ", *dbFile)

	for _, path := range flag.Args() {
		walkPaths(path)
	}
}

func walkPaths(path string) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// If current path is a directory, skip it.
			if info.IsDir() {
				return nil
			}

			// If file does not exist in database, hash it and add to database.
			// If file exists in database, check if modification time changed. Hash if changed.
			var sha1hash, error = hash_file_sha1(path)
			if error != nil {
				return error
			}
			fmt.Println(path, info.Size(), info.ModTime(), sha1hash)

			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func hash_file_sha1(filePath string) (string, error) {
	//Initialize variable returnSHA1String now in case an error has to be returned
	var returnSHA1String string

	//Open the filepath passed by the argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new SHA1 hash interface to write to
	hash := sha1.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA1String, err
	}

	//Get the 20 bytes hash
	hashInBytes := hash.Sum(nil)[:20]

	//Convert the bytes to a string
	returnSHA1String = hex.EncodeToString(hashInBytes)

	return returnSHA1String, nil
}
