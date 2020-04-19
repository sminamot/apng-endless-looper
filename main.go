package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/kettek/apng"
)

type option struct {
	overwrite bool
	suffix    string
}

// flag
var (
	w bool   // overwrite
	s string // filename's suffix
)

func main() {
	log.SetOutput(os.Stderr)
	flag.BoolVar(&w, "w", false, "overwrite")
	flag.StringVar(&s, "s", "_loop", "new filename's suffix")
	flag.Parse()
	o := &option{
		overwrite: w,
		suffix:    s,
	}
	args := flag.Args()
	os.Exit(run(args, o))
}

func run(args []string, o *option) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <target_file1> [<target_file2> ...]\n", os.Args[0])
		return 1
	}
	for _, file := range args {
		func() {
			f, err := os.Open(file)
			if err != nil {
				log.Printf("failed to open %s, ignored...\n", file)
				return
			}
			defer f.Close()

			a, err := apng.DecodeAll(f)
			if err != nil {
				log.Printf("failed to decode %s, ignored...\n", file)
				return
			}

			if len(a.Frames) == 1 {
				return
			}

			var newFile *os.File
			if o.overwrite {
				tmp, err := ioutil.TempFile("", "tmp")
				if err != nil {
					log.Println("failed to create new file, ignored...")
					return
				}
				defer os.Remove(tmp.Name())
				newFile = tmp
			} else {
				dir := filepath.Dir(filepath.Clean(f.Name()))
				base := filepath.Base(f.Name())
				ext := filepath.Ext(f.Name())
				fileName := filepath.Join(dir, base[:len(base)-len(ext)]+o.suffix+ext)
				newFile, err = os.Create(fileName)
				if err != nil {
					log.Println("failed to create new file, ignored...")
					return
				}
			}

			a.LoopCount = 0
			if err := apng.Encode(newFile, a); err != nil {
				log.Printf("failed to encode %s, ignored...\n", file)
				return
			}

			if o.overwrite {
				if err := os.Rename(newFile.Name(), f.Name()); err != nil {
					log.Println("failed to create new file, ignored...")
					return
				}
			}
		}()
	}

	return 0
}
