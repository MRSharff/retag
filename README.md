# Retag
Easier bulk file renaming and retagging on the command line.

I know there is a 3 liner with the sed tool for bulk renaming but you better get it
right the first time. There is big room for error with it and I wanted something a
a little more forgiving.

Uses Golangs regex syntax. See https://golang.org/s/re2syntax

This is super early development and the flags will probably change.

TODO:
- [ ] Interactive?
- [ ] GUI (mp3tag style is goal)
- [ ] Tag support
- [ ] Filename from Tag
- [ ] Tag from Filename
- [ ] Support output formatting for number (leading zeroes, etc)

## Install
If you have go installed:
`go install github.com/MRSharff/retag`

## Run

### Example
`retag -o "(.*) - (.*).mp3" -n "$\{2\} $\{1\} (2019).mp3" -s "_" MyFavoriteSong\ -\ 01.mp3`

Output:

```
Old                     New
MyFavoriteSong - 01.txt 01_MyFavoriteSong_(2019).txt
Confirm rename?(y/n): y
```

### Testing Your Command
Use the -t flag to test your pattern and see if you are able to match a sample filename and see which groups are matched.
`retag -t -o "(.*) - (.*).txt" -n "$\{2\} $\{1\} (2019).txt" -s "_" MyFavoriteSong\ -\ 01.txt`

Output:

```
Pattern:  (.*) - (.*).txt
Filename: MyFavoriteSong - 01.txt
Groups: 
${1}: MyFavoriteSong
${2}: 01
```
