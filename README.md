# topfew
A program that finds records in which a 
certain field or combination of fields occurs  
most frequently

## Usage

```shell
tf 
  -n, --n [number of lines]
  -f, --fields [fieldlist]
  -h, -help, --help
  -g, --grep [regexp]
  -v, --vgrep [regexp]
  -s, --sed [regexp] [replacement]
  -w, --width [number of file segments]
  -sample
  [filename]
```
## Options
`-n integer`, `--number integer` How many of the highest‐occurrence‐count lines to print out. The
default value is 10.

`-f fieldlist, --fields fieldlist` Specifies which fields should be extracted from incoming records
and used in computing occurrence counts. The fieldlist must be a
comma‐separated  list  of  integers  identifying  field numbers,
which start at one, for example 3 and 2,5,6.  The fields
must be provided in order, so 3,1,7 is an error.

If no fieldlist is provided, **tf** treats the whole input record as a single field.

`-g, regexp`, `--grep regexp`

The  initial **g** suggests `grep`. These options apply the provided
regular expression to, respectively, each record as it  is  read
and  each  field‐set  as it is extracted, and if the regexp does
not match the record or field, cause tf to bypass the record.

These options can be provided multiple times; the provided regu‐
lar  expressions will be applied in the order they appear on the
command line.

`-v regexp`, `--vgrep regegxp`

The initial **v** suggests "grep ‐v". These operations are  the  in‐
verse  of  `‐grecord` and `‐gfield`, rejecting records and extracted
fields that match the  provided  regular  expression.   As  with
those operations, these can be provided multiple times.

`-s regexp replacement`, `--sed regexp replacement`

As its name suggests, applies sed‐style editing by replacing any
text that matches the provided regexp with the provided replace‐
ment.   It  works on the fields in the fieldlist after they have
been extracted from the record.

If ()‐enclosed capturing groups appear in the regexp,  they  may
be referred to as **$1**, **$2**, and so on in, the replacement.

This  option can be provided many times, and the replacement op‐
erations are performed in the order they appear on  the  command
line.

`--sample`

It can be tricky to get the regular expressions in the `−g`,
`−v`, and `−s` options  right.  Specifying
`-−sample`  causes  **tf**  to  print lines to the standard output that
display the filtering and field‐editing logic.  It can  only  be
used when processing standard input, not a file.

`-w integer`, `--width integer`

If a file name is specified then **tf**, rather than reading it from
end to end, will divide it into segements and process it in multiple 
parallel threads. The optimal number of threads depends in a 
complicated way on how many cores your CPU has what kind of cores
they are, and the storage architecture.

The default is the result of the Go `runtime.NumCPU()` calls and
often produces good results.

`-h`, `-help`, `--help`

Describes the function and options of tf.

## Examples

To find the IP address that most commonly hits your
web site, given an Apache logfile named `access_log`

`tf -fields 1 access_log`

The same effect could be achieved with

`awk '{print $1}' access_log | sort | uniq -c | sort -rn | head`

But tf is usualy much faster.

Do the same, but exclude high-traffic bots (omiting `access_log`)

`tf -fields 1 -vrecord googlebot -vrecord bingbot` 

Most popular IP addresses from May 2020.

`tf -fields 1 -grecord '\[../May/2020' `

Most popular hour/minute of the day for retrievals

`tf -fields 4 -sed "\\[" ""  -sed '^[^:]*:' ''  -sed ':..$' '' `

