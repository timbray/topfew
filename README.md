# topfew

[![Tests](https://github.com/timbray/topfew/actions/workflows/tests.yaml/badge.svg)](https://github.com/timbray/topfew/actions/workflows/tests.yaml)
[![codecov](https://codecov.io/gh/timbray/topfew/branch/main/graph/badge.svg)](https://codecov.io/gh/timbray/topfew)
[![Go Report Card](https://goreportcard.com/badge/github.com/timbray/topfew)](https://goreportcard.com/report/github.com/timbray/topfew)
[![timbray/topfew](https://img.shields.io/github/go-mod/go-version/timbray/topfew)](https://github.com/timbray/topfew)


A program that finds and prints out the top few records in which a certain field or combination of fields occurs most frequently.

This is release 2.0 of Topfew.

## Examples

To find the IP address that most commonly hits your web site, given an Apache logfile named `access_log`.

`topfew --fields 1 access_log`

The same effect could be achieved with

`awk '{print $1}' access_log | sort | uniq -c | sort -rn | head`

But **topfew** is usually much faster.

Do the same, but exclude high-traffic bots (omitting the filename).

`topfew --fields 1 --vgrep googlebot --vgrep bingbot`

Most popular IP addresses from May 2020.

`topfew --fields 1 --grep '\[../May/2020'`

Most popular hour/minute of the day for retrievals.

`topfew --fields 4 --sed "\\[" ""  --sed '^[^:]*:' ''  --sed ':..$' ''`

## Usage

```shell
topfew
	-n, --number (output line count) [default is 10]
	-f, --fields (field list) [default is the whole record]
	-q, --quotedfields [respect "-delimited space-separated fields]
	-p, --fieldseparator (regexp) [use provided regexp to separate fields]
	-g, --grep (regexp) [may repeat, default is accept all]
	-v, --vgrep (regexp) [may repeat, default is reject none]
	-s, --sed (regexp) (replacement) [may repeat, default is no changes]
	-w, --width (segment count) [default is result of runtime.numCPU()]
	--sample
	-h, -help, --help
	filename [default is stdin]

All the arguments are optional; if none are provided, topfew will read records 
from the standard input and list the 10 which occur most often.
```
## Options
`-n integer`, `--number integer` 

How many of the highest‐occurrence‐count lines to print out. 
The default value is 10.

`-f fieldlist, --fields fieldlist`

Specifies which fields should be extracted from incoming records and used in computing occurrence counts.
The fieldlist must be a comma‐separated  list  of  integers  identifying  field numbers, which start at one, for example 3 and 2,5,6.
The fields must be provided in order, so 3,1,7 is an error.

If no fieldlist is provided, **topfew** treats the whole input record as a single field.

`-p separator, --fieldseparator separator` 

Provides a regular expression that is used as a field separator instead of the default white space.
This is likely to incur a significant performance cost.

`-q, --quotedfields`

Some files, for example Apache httpd logs, use space-separation but also
allow spaces within fields which are delimited by `"`. The -q/--quotedfields
argument allows **topfew** to process these correctly. It is an error to specify both
-p and -q.

`-g regexp`, `--grep regexp`

The  initial **g** suggests `grep`.
This option applies the provided regular expression to each record as it is read and if the regexp does not match the record, **topfew** bypasses it.

This option can be provided multiple times; the provided regular expressions will be applied in the order they appear on the command line.

`-v regexp`, `--vgrep regegxp`

The initial **v** suggests `grep ‐v`. This operation is the  inverse  of `-g` and `-‐grep`, rejecting records that match the  provided regular  expression.  
As  with `grep`, it can be provided multiple times.

`-s regexp replacement`, `--sed regexp replacement`

As its name suggests, applies sed‐style editing by replacing any text that matches the provided regexp with the provided replacement.
It  works on the fields in the fieldlist after they have been extracted from the record.

If ()‐enclosed capturing groups appear in the regexp,  they  may be referred to as **$1**, **$2**, and so on in, the replacement.

This  option can be provided many times, and the replacement operations are performed in the order they appear on  the  command line.

`--sample`

It can be tricky to get the regular expressions in the `−g`, `−v`, and `−s` options  right.
Specifying `-−sample`  causes  **topfew**  to  print lines to the standard output that display the filtering and field‐editing logic.
It can  only  be used when processing standard input, not a file.

`-w integer`, `--width integer`

If a file name is specified then **topfew**, rather than reading it from end to end, will divide it into segments and process it in multiple parallel threads.
The optimal number of threads depends in a complicated way on how many cores your CPU has what kind of cores they are, and the storage architecture.

The default is the result of the Go `runtime.NumCPU()` calls and often produces good results.

`-h`, `-help`, `--help`

Describes the function and options of **topfew**.

## Records and fields

Records are separated by newlines, fields within records by white space, defined as one or more space or tab characters.

The field separator can be overridden with the --fieldseparator option.

## Case study: Apache access_log

Here is a line from an Apache httpd `access_log` file. For readability, the fields are 
separated by line-breaks and numbered. Note that the fields are mostly space-separated, but that field 6,
summarizing the request and its result, is delimited by quote characters `"`.

```
1. 202.113.19.244 
2. - 
3. - 
4. [12/Mar/2007:08:04:39 
5. -0800] 
6. "GET /ongoing/picInfo.xml?o=http://www.tbray.org/ongoing/When/200x/2007/03/10/Beautiful-Code HTTP/1.1" 
7. 200 
8. 137 
9. "http://www.tbray.org/ongoing/When/200x/2007/03/10/Beautiful-Code" 
10. "Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.2) Gecko/20070219 Firefox/2.0.0.2"
```

The fetch of `picInfo.xml` signals that this is an actual browser request, likely signifying that 
a human was involved; the URL following the `o=` is the resource the human looked at. Here is a 
**topfew** invocation that yields a list of the top 5 URLs that were fetched by a human:

```shell
topfew -g picInfo.xml -f 6 -q -s '\?utm.*' '' -s " HTTP/..." "" -s "GET .*\/ongoing" ""
```

Note the `-g` to select only lines with `picInfo.xml`, the `-q` to request correct processing
of quote-delimited fields, and the sequence of `-s` patterns to clean up the results.

## Performance issues

Since the effect of topfew can be exactly duplicated with a combination of `awk`, `grep`, `sed` and `sort`, you wouldn’t be using it if you didn’t care about performance. 
Topfew is quite highly tuned and pushes your computer’s I/O subsystem and Go runtime hard.
Therefore, the observed effects of combinations of options can vary dramatically from system to system.

For example, if I want to list the top records containing the string `example` from a file named `big-file` I could do either of the following:

```shell
topfew -g example big-file 
grep example big-file |topfew
```

When I benchmark topfew on a modern Apple-Silicon Mac and an elderly spinning-rust Linux VPS, I observe that the first option is faster on Mac, the second on Linux.

Only one performance issue is uncomplicated: Topfew will **always** run faster on a named file than a standard-input stream.

## Credits

Tim Bray created version 0.1 of Topfew, and the path toward 1.0 was based chiefly on ideas stolen from Dirkjan Ochtman and contributed by Simon Fell.
The GitHub CI was based on Michael Gasch’s implementation from my Quamina repository, and he helped with Topfew’s as well.
