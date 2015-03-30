package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var VERSION string = "0.0.0-src" //set via ldflags

var help = `
	Usage: dedup [options] <dir> [dir] [dir]
	
	Version: ` + VERSION + `

	Dedup de-deduplicates all provided directories by merging
	all into the first. The merge operation simultaneously
	removes duplicates and renames files due to path collisions.
	
	Options:
	  --keep, keep duplicates (by default, duplicates are deleted)
	  -v, verbose logs (display each move and delete)
	  --version, display version
	  -h --help, this help text

	Notes:
	  * dedup considers two files duplicates if they have
	    matching sha1 sums
	  * dedup is not recursive (only works on files)
	  * dedup is a destructive operation (unless --keep)
	  * dedup on a single directory will only perform
	    deduplication, no moves
	  * dedup renames: when a file is unique, dedup will
	    attempt to move the file. if the path already
	    exists the incoming file will be suffixed with
	    the next number (for example, if 'foo.txt' exists,
	    the new file will be 'foo-2.txt')
	  * any error will cause dedup to exit

	Read more: https://github.com/jpillora/dedup

`

var (
	mut = &sync.Mutex{}
	//mut guards these 3 vars
	deletes = 0
	moves   = 0
	hashes  = map[string]string{}
)

func report() string {
	return fmt.Sprintf("moves %d deletes %d", moves, deletes)
}

func main() {
	h1 := flag.Bool("h", false, "")
	h2 := flag.Bool("help", false, "")
	verbose := flag.Bool("v", false, "")
	version := flag.Bool("version", false, "")
	keep := flag.Bool("keep", false, "")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, help)
		os.Exit(1)
	}
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *h1 || *h2 {
		flag.Usage()
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
	}

	dst := args[0]
	info, err := os.Stat(dst)
	check(err)

	if !info.IsDir() {
		check(errors.New("Must be directory"))
	}

	//========================

	//use all da cpus
	runtime.GOMAXPROCS(runtime.NumCPU())

	files, err := ioutil.ReadDir(dst)
	check(err)

	//initialize index
	fmt.Printf("Indexing '%s' (#%d files)\n", dst, len(files))

	each(files, func(f os.FileInfo) {
		n := f.Name()
		if f.IsDir() || strings.HasPrefix(filepath.Base(n), ".") {
			return //skip hidden
		}
		dstpath := filepath.Join(dst, n)
		hex := hash(dstpath)
		//guard hashes index
		mut.Lock()
		defer mut.Unlock()
		//look at index
		if other, exists := hashes[hex]; exists {
			if !*keep {
				if *verbose {
					fmt.Printf("  Remove duplicate '%s' (of '%s')\n", n, other)
				}
				err := os.Remove(dstpath)
				check(err)
				deletes++
			}
			return
		}
		hashes[hex] = n
		if *verbose {
			fmt.Printf("  Indexed '%s' (%s)\n", n, hex)
		}
	})

	srcs := args[1:]
	//compare others against index
	for _, src := range srcs {
		files, err := ioutil.ReadDir(src)
		check(err)

		fmt.Printf("Merging in '%s' (#%d files)\n", src, len(files))
		each(files, func(f os.FileInfo) {
			n := f.Name()
			if f.IsDir() || strings.HasPrefix(filepath.Base(n), ".") {
				return
			}
			srcpath := filepath.Join(src, n)
			hex := hash(srcpath)
			//guard hashes index
			mut.Lock()
			defer mut.Unlock()
			//look at index
			_, exists := hashes[hex]
			if exists {
				if !*keep {
					if *verbose {
						fmt.Printf("  Remove duplicate: %s\n", n)
					}
					err := os.Remove(srcpath)
					check(err)
					deletes++
				}
				return
			}
			//mark as exists
			hashes[hex] = n
			//find next availbe file name
			dstpath := filepath.Join(dst, n)
			i := 1
			for {
				if _, err := os.Stat(dstpath); os.IsNotExist(err) {
					break
				}
				i++
				ext := filepath.Ext(n)
				name := strings.TrimSuffix(n, ext)
				newname := fmt.Sprintf("%s-%d%s", name, i, ext)
				dstpath = filepath.Join(dst, newname)
			}
			//missing, move in
			if *verbose {
				fmt.Printf("  Moving: %s -> %s\n", n, dstpath)
			}
			err = os.Rename(srcpath, dstpath)
			check(err)
			moves++
		})
	}

	fmt.Printf("Done (%s)\n", report())
}

//fan-out for-each over fileinfos
func each(files []os.FileInfo, fn func(os.FileInfo)) {
	queue := make(chan os.FileInfo)
	wg := &sync.WaitGroup{}
	// the work operation
	work := func() {
		defer wg.Done()
		for f := range queue {
			fn(f)
		}
	}
	// add #CPU many workers and start all
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go work()
	}
	// queue all work
	for _, f := range files {
		queue <- f
	}
	//all items queued, mark end
	close(queue)
	//will only yield once all workers
	//have completed their task
	wg.Wait()
}

//hash the file at path to an sha1 hex string
func hash(p string) string {
	f, err := os.Open(p)
	check(err)
	h := sha1.New()
	io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil))
}

//halt program if error is encountered
func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s (%s)\n", err, report())
		os.Exit(1)
	}
}
