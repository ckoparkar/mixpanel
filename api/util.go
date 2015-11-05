package api

import (
	"fmt"
	"io/ioutil"
	"log"
)

func OverwriteFile(filename string, data []byte) {
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		// if we couldnt open file, try writing to tmp file.
		f, _ := ioutil.TempFile(".", "")
		defer f.Close()
		if err := ioutil.WriteFile(f.Name(), data, 0644); err != nil {
			log.Printf("[ERR] %s, writing to STDOUT.\n")
			fmt.Println(string(data))
		} else {
			log.Printf("[ERR] Couldnt write to %s: %s, instead wrote to: %s", filename, err, f.Name())
		}
	}
}
