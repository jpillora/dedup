# Dedup

Dedup is a command-line tool which deduplicates all files in the provided directories by merging them together into the first directory. The merge operation simultaneously removes duplicates and renames files (when a path collision occurs).

### Install

**Binaries**

[![Releases](https://img.shields.io/github/release/jpillora/dedup.svg)](https://github.com/jpillora/dedup/releases) [![Releases](https://img.shields.io/github/downloads/jpillora/dedup/total.svg)](https://github.com/jpillora/dedup/releases)

See [the latest release](https://github.com/jpillora/dedup/releases/latest) or download and install it now with `curl https://i.jpillora.com/dedup! | bash`

**Source**

```sh
$ go get -v github.com/jpillora/dedup
```

### Usage

```
$ dedup --help

  Usage: dedup [options] directories...

  deduplicates all files in the provided directories, while optionally merging
  them into the first directory. The merge operation renames files (when a path
  collision occurs).

  Options:
  --keep, -k       keep duplicates (by default, duplicates are deleted)
  --merge, -m      move unique files into the first directory
  --recursive, -r  searches into nested directories
  --verbose, -v    verbose logs (displays each move and delete)
  --dryrun, -d     runs exactly as configured, except no changes are made
  --workers, -w    number of worker threads (default 4)
  --hash, -h       hashing algorithm to use (default md5, can also choose
                   sha1 and sha256)
  --help
  --version

  Notes:
  * dedup considers two files duplicates if they have matching hash sums.
  * dedup is a destructive operation (unless --keep).
  * dedup on a single directory will only perform deduplication, no moves.
  * dedup renames: when a file is unique, dedup will attempt to move the file.
    if the path already exists the incoming file will be suffixed with the next
    number (for example, if 'foo.txt' exists, the new file will be 'foo-2.txt').
  * enabling 'keep' without 'merge' enabled is a no-op.
  * any error will cause dedup to exit.

  Version:
    0.0.0-src

  Read more:
    github.com/jpillora/dedup
```

### Example

```sh
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

$ dedup --verbose --dryrun --merge foo bar
[DRYRUN] Scanning foo (no changes)
[DRYRUN] Scanning bar (hashed 2)
[DRYRUN] Removing bar/foo.txt dupe-of foo/foo.txt (hashed 3)
[DRYRUN] Moving: bar/bar.txt -> foo/bar.txt (hashed 4)
[DRYRUN] Removing bar/foo-copy.txt dupe-of foo/foo.txt (hashed 5, moved 1, deleted 1)
[DRYRUN] Done (hashed 5, moved 1, deleted 2)

# looks good!

$ dedup --merge foo bar
Done (hashed 5, moved 1, deleted 2)

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
