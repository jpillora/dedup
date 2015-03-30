package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const help = `
	Usage: dedup <destination> [source] [source]

	Dedup indexes and removes all duplicates from <destination> and
	then merges in all files from each [source] while deleting duplicates.
	If files are different, though they have the same name, a number will
	be prefixed to the filename (if 'foo.txt' exists, the new file will
	be 'foo-2.txt'). If you only specify a <destination> then only de-
	deplication will be performed.
	
	Read more: https://github.com/jpillora/dedup

`

func main() {
	h1 := flag.Bool("h", false, "")
	h2 := flag.Bool("help", false, "")
	flag.Parse()

	if *h1 || *h2 {
		fmt.Fprint(os.Stderr, help)
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 1 {
		check(errors.New("Missing destination"))
	}

	dst := args[0]
	srcs := args[1:]

	info, err := os.Stat(dst)
	check(err)

	if !info.IsDir() {
		check(errors.New("Must be directory"))
	}

	files, err := ioutil.ReadDir(dst)
	check(err)

	hashes := map[string]string{}

	//initialize index
	fmt.Printf("Indexing '%s' (#%d files)\n", dst, len(files))
	for _, f := range files {
		n := f.Name()
		if f.IsDir() || strings.HasPrefix(filepath.Base(n), ".") {
			continue //skip hidden
		}
		dstpath := filepath.Join(dst, n)
		hex := hash(dstpath)
		if other, exists := hashes[hex]; exists {
			fmt.Printf("  Duplicate '%s' of '%s'\n", n, other)
			err := os.Remove(dstpath)
			check(err)
			continue
		}
		hashes[hex] = n
		fmt.Printf("  Indexed '%s' (%s)\n", n, hex)
	}

	//compare others against index
	for _, src := range srcs {
		files, err := ioutil.ReadDir(src)
		check(err)

		fmt.Printf("Merging %s (#%d files)\n", src, len(files))
		for _, f := range files {
			n := f.Name()
			if f.IsDir() || strings.HasPrefix(filepath.Base(n), ".") {
				continue
			}
			srcpath := filepath.Join(src, n)
			_, exists := hashes[hash(srcpath)]
			if exists {
				fmt.Printf("  Remove duplicate: %s\n", n)
				err := os.Remove(srcpath)
				check(err)
				//delete?
				continue
			}
			//find next availbe file name
			dstpath := filepath.Join(dst, n)
			i := 1
			for {
				if _, err := os.Stat(dstpath); os.IsNotExist(err) {
					break
				}
				i++
				ext := filepath.Ext(n)
				base := strings.TrimSuffix(n, ext)
				newname := fmt.Sprintf("%s-%d%s", base, i, ext)
				dstpath = filepath.Join(dst, newname)
			}
			//missing, move in
			fmt.Printf("  Moving: %s -> %s\n", n, dstpath)
			err = os.Rename(srcpath, dstpath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Printf("Done\n")
}

func hash(p string) string {
	f, err := os.Open(p)
	check(err)
	h := sha1.New()
	io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil))
}

func check(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
