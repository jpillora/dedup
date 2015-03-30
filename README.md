# Dedup

Dedup is a command-line tool which deduplicates all files in the provided directories by merging them together into the first directory. The merge operation simultaneously removes duplicates and renames files (when a path collision occurs).

### Install

**Binaries**

See [Releases](https://github.com/jpillora/dedup/releases/latest)

**Source**

``` sh
$ go get -v github.com/jpillora/dedup
```

### Usage

```
$ dedup --help
```

<tmpl,code: go run main.go --help>
```

	Usage: dedup [options] <dir> [dir] [dir]
	
	Version: 0.0.0-src

	deduplicates all files in the provided directories
	by merging them together into the first directory.
	The merge operation simultaneously removes duplicates
	and renames files (when a path collision occurs).
	
	Options:
	  --keep, keep duplicates (by default, duplicates
	    are deleted)
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

```
</tmpl>

### Example

``` sh
$ cd example/
$ tree .
.
├── bar
│   ├── bar.txt
│   ├── foo-copy.txt
│   └── foo.txt
└── foo
    ├── bazz.txt
    └── foo.txt

2 directories, 5 files

$ dedup foo bar
Indexing 'foo' (#2 files)
Merging in 'bar' (#3 files)
Done (moves 1 deletes 2)

$ tree .
.
├── bar
└── foo
    ├── bar.txt
    ├── bazz.txt
    └── foo.txt

2 directories, 3 files
```

#### MIT License

Copyright © 2015 Jaime Pillora &lt;dev@jpillora.com&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.