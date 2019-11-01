// You can edit this code!
// Click here and start typing.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Parse command line flags.
	dbFile := flag.String("db", "", "The database file to store the hashes in.")
	flag.Parse()

	// Output
	fmt.Println(*dbFile)
	fmt.Println("tail:", flag.Args())

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
			fmt.Println(path, info.Size(), info.ModTime())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
