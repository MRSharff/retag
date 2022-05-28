package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/tabwriter"
)

var (
	pattern      = flag.String("o", "(.*)", "The regex pattern to match against the filename.")
	outputFormat = flag.String("n", "${1}", `The format of the new filename. Use $\{1\}, $\{2\}, ..., $\{n\} for captured groups.`)
	separator    = flag.String("s", " ", "Character to optionally replace whitespace with.")
	verbose      = flag.Bool("v", false, "Show more output.")
	testOnly     = flag.Bool("t", false, "Test your regex pattern on a sample input")
	noPrompt     = flag.Bool("y", false, "Do not prompt for confirmation when renaming files.")
)

var replacer *regexp.Regexp
var tw = tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
var directories, files, renamedFiles []string

func main() {
	flag.Parse()
	if !*verbose {
		log.SetOutput(io.Discard)
	}

	removeEscapesFromOutputFormat()

	filecount := flag.NArg()
	if filecount == 0 {
		fmt.Fprintln(os.Stderr, "You must provide a string or a list of files (through globbing) to work on.")
		return
	}

	if *testOnly {
		testPattern()
		return
	}
	replacer = regexp.MustCompile(*pattern)

	// I want to split these now so that I can just work on the list of filenames,
	// but still have the directories by index to reference later.
	directories = make([]string, filecount)
	files = make([]string, filecount)
	for i := 0; i < filecount; i++ {
		directories[i], files[i] = path.Split(flag.Arg(i))
		log.Println(flag.Arg(i))
	}

	renamedFiles = make([]string, filecount)
	for i, old := range files {
		renamedFiles[i] = rename(old)
	}

	if len(renamedFiles) != len(directories) || len(renamedFiles) != len(files) {
		panic("Different amount of files to be renamed than given.\n" +
			"This really shouldn't happen.")
	}

	printRenamings()

	if confirm() {
		renameFiles()
	} else {
		fmt.Println("Rename canceled.")
	}

}

func printRenamings() {
	fmt.Fprintln(tw, "Old\tNew")
	for i := 0; i < len(files); i++ {
		fmt.Fprintln(tw, files[i]+"\t"+renamedFiles[i])
	}
	tw.Flush()
}

func renameFiles() {
	for i := 0; i < len(renamedFiles); i++ {
		oldPath := path.Join(directories[i], files[i])
		newPath := path.Join(directories[i], renamedFiles[i])
		log.Println("Renaming", oldPath, "=>", newPath)
		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

const maxTries = 5

func confirm() bool {
	fmt.Print("Confirm rename?(y/n): ")
	reader := bufio.NewReader(os.Stdin)
	for tries := 0; tries < maxTries; tries++ {
		c, err := reader.ReadString('\n')
		c = strings.TrimSpace(c)
		c = strings.ToLower(c)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Bad confirmation scan: ", err.Error())
		}
		log.Println("Checking confirmation received")
		switch c {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Print("Unrecognized value ", c, "please enter y or n: ")
		}
	}
	return false
}

// removeEscapesFromOutputFormat removes the escapes from the pattern that are
// required during input so that the shell doesn't expand the curly braces
// around the capture groups.
func removeEscapesFromOutputFormat() {
	// TODO: I'm curious if there is a better way to handle this
	*outputFormat = strings.ReplaceAll(*outputFormat, "\\", "")
}

// testPattern outputs useful information when trying to create a pattern to
// replace such as the groups captured.
func testPattern() {
	_, filename := filepath.Split(flag.Arg(0))

	fmt.Println("Pattern: ", *pattern)
	fmt.Println("Filename:", filename)

	var err error
	replacer, err = regexp.Compile(*pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// TODO: Can I try subpatterns and see if groups up until a certain point match?
	// 	Like If I do "testname-(.*) - (.?).txt" I could try to search "testname-(.*)"
	// 	to see if it finds anything, if it does, move up to the next one, "testname-(.*) - (.?)"
	// 	and so on.

	groups := replacer.FindStringSubmatch(filename)
	if len(groups) == 0 {
		fmt.Fprintln(os.Stderr, "No match found.")
		return
	}

	fmt.Println("Groups: ")
	for i, g := range groups[1:] {
		fmt.Printf("${%d}: %s\n", i+1, g)
	}
}

func rename(old string) string {
	// TODO: Revisit this. What if outputFormat contains a $ symbol?
	return strings.ReplaceAll(replacer.ReplaceAllString(old, *outputFormat), " ", *separator)
}
