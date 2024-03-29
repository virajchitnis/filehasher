package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var filename string
var verbose bool
var update bool
var files int64

func main() {
	// Parse command line flags.
	dbFile := flag.String("db", "filehasher.db3", "The database file to store the hashes in")
	displayAll := flag.Bool("verbose", false, "Verbose mode")
	doUpdate := flag.Bool("update", false, "Update the database if files have been added or updated")
	flag.Parse()

	// Output
	filename = *dbFile
	verbose = *displayAll
	update = *doUpdate

	// Database setup
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create files table in database if it does not exist.
	sqlStmt := `
	create table if not exists files (path string UNIQUE, size integer, mod_time date, hash string);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// Go through all files and directories that have been passed to the program.
	files = 0
	for _, path := range flag.Args() {
		walkPaths(path)
	}

	fmt.Println()
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

			files++
			fmt.Printf("\r[%d] ", files)

			// Query database for this particular file.
			db, err := sql.Open("sqlite3", filename)
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()

			stmt, err := db.Prepare("select size, mod_time, hash from files where path = ?")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			var size int64
			var mod_time time.Time
			var hash string
			err = stmt.QueryRow(path).Scan(&size, &mod_time, &hash)
			switch err {
			case sql.ErrNoRows:
				if update {
					// If file does not exist in database, hash it and add to database.
					var sha1hash, error = hash_file_sha1(path)
					if error != nil {
						return error
					}

					tx, err := db.Begin()
					if err != nil {
						log.Fatal(err)
					}
					stmt, err := tx.Prepare("insert into files values(?, ?, ?, ?)")
					if err != nil {
						log.Fatal(err)
					}
					defer stmt.Close()
					_, err = stmt.Exec(path, info.Size(), info.ModTime().UTC(), sha1hash)
					if err != nil {
						log.Fatal(err)
					}
					tx.Commit()

					fmt.Println("\u001b[34mAA --\u001b[0m", path, info.Size(), info.ModTime().UTC(), sha1hash)
				}
			case nil:
				// Hash the file to check data consistency.
				var sha1hash, error = hash_file_sha1(path)
				if error != nil {
					return error
				}
				// If file exists in database, check if modification time changed.
				var changed = false
				if info.Size() != size {
					changed = true
				}
				if !info.ModTime().Equal(mod_time) {
					changed = true
				}

				if changed {
					// File has changed.
					fmt.Println("\u001b[35mUO --\u001b[0m", path, size, mod_time, hash)
					fmt.Println("\u001b[35mUN --\u001b[0m", path, info.Size(), info.ModTime().UTC(), sha1hash)

					if update {
						tx, err := db.Begin()
						if err != nil {
							log.Fatal(err)
						}
						stmt, err := tx.Prepare("update files set size = ?, mod_time = ?, hash = ? where path = ?")
						if err != nil {
							log.Fatal(err)
						}
						defer stmt.Close()
						_, err = stmt.Exec(info.Size(), info.ModTime().UTC(), sha1hash, path)
						if err != nil {
							log.Fatal(err)
						}
						tx.Commit()
					}
				} else {
					// File has not changed.
					if sha1hash != hash {
						// File is damaged.
						fmt.Println("\u001b[31mDO --\u001b[0m", path, size, mod_time, hash)
						fmt.Println("\u001b[31mDN --\u001b[0m", path, info.Size(), info.ModTime().UTC(), sha1hash)
					} else if verbose {
						fmt.Println("\u001b[32m-- --\u001b[0m", path, size, mod_time, hash)
					}
				}
			}

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
