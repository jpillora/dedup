# dedeup

Deduplicate files

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

<tmpl,code: dedup --help>
```

	Usage: dedup [options] <dir> [dir] [dir]

	Dedup de-deduplicates all directories, combines all dirs into
	the first <dir>, while removing duplicates.

	If two files having matching paths, the second file will be
	prefixed to the filename (if 'foo.txt' exists, the new file
	will be 'foo-2.txt').

	If you only specify one <dir> then only dedeplication will be
	performed.
	
	Options:
	  --keep, keep duplicates (by default, duplicates are deleted)
	  -v, verbose logs (display moves and deletes)

	Warnings:
	  * dedup is not recursive
	  * dedup is a destructive operation (unless --keep)
	  * any error will cause dedup to exit

	Read more: https://github.com/jpillora/dedup

```
</tmpl>

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