package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/vallerion/rscanner"
)

func handle(err error, context string, args ...any) {
	if err == nil {
		return
	}
	fmt.Printf("%s\nError: %v", fmt.Sprintf(context, args...), err)
	os.Exit(1)
}

func getHomeDir() string {
	user, err := user.Current()
	handle(err, "Failed to get the current user")
	return user.HomeDir
}

const cacheFileName = "~/.active.txt"

var homeDir = getHomeDir()

func parsePath(path string) string {
	path = strings.ReplaceAll(path, "~", homeDir)
	res, err := filepath.Abs(path)
	handle(err, "Failed to get absolute path: %v", err)
	return res
}

func main() {
	cachePath := parsePath(cacheFileName)

	list := flag.UintP("list", "l", 0, "Will list out the most recent paths.")
	flag.Lookup("list").NoOptDefVal = "10"
	prevHelp := `Will get the n-th previous path. 
Using -n or --n where n is a number will also work.`
	prev := flag.UintP("previous", "p", 0, prevHelp)
	printCache := flag.BoolP("cache", "c", false, "Will print path to the cache")

	if flag.NArg() > 2 {
		handle(fmt.Errorf("Too many arguments: %v", flag.NArg()), "")
		return
	}

	// Allow for -1 to be -p=1
	for i, arg := range os.Args {
		if !strings.HasPrefix(arg, "-") {
			continue
		}
		num, err := strconv.ParseInt(arg[1:], 10, 32)
		if err == nil {
			if num < 0 {
				num = -num
			}
			os.Args[i] = "-p=" + strconv.FormatInt(num, 10)
		}
	}

	flag.Parse()

	if printCache != nil && *printCache {
		fmt.Print(cachePath)
		os.Exit(4)
		return
	}

	cache, err := os.OpenFile(cachePath, os.O_CREATE|os.O_RDWR, 0644)
	handle(err, "Failed to get or create cache at '%s'", cachePath)
	cacheStat, err := cache.Stat()
	handle(err, "Failed to stat cache")

	getLastLine := func(num uint) string {
		bscanner := rscanner.NewScanner(cache, cacheStat.Size())
		var lastLine []byte
		i := uint(0)
		for bscanner.Scan() {
			lastLine = bscanner.Bytes()
			if len(lastLine) == 0 {
				continue
			}
			if i == num {
				return strings.TrimSpace(string(lastLine))
			}
			i++
		}
		handle(bscanner.Err(), "Error scanning cache")
		return strings.TrimSpace(string(lastLine))
	}

	getLastLines := func(num uint) []string {
		bscanner := rscanner.NewScanner(cache, cacheStat.Size())
		lastLines := make([]string, 0, num)
		if num == 0 {
			return lastLines
		}
		i := uint(0)
		for bscanner.Scan() {
			bytes := bscanner.Bytes()
			if len(bytes) == 0 {
				continue
			}

			lastLines = append(lastLines, strings.TrimSpace(string(bytes)))
			i++
			if i == num {
				break
			}
		}
		handle(bscanner.Err(), "Error scanning cache")
		return lastLines
	}

	getDigits := func(num uint) int {
		if num == 0 {
			return 1 // Special case for 0, which has one digit
		}
		count := 0
		for num != 0 {
			num /= 10
			count++
		}
		return count
	}

	getLineFormat := func(maxNum uint) string {
		digits := getDigits(maxNum)
		return "%" + strconv.Itoa(digits) + "d: %s\n"
	}

	// Handle list flag
	if list != nil && *list > 0 {
		lineFormat := getLineFormat(*list)
		for i, line := range getLastLines(*list) {
			fmt.Printf(lineFormat, i, line)
		}
		return
	}

	doFuzzyList := func() {
		defer os.Exit(3)

		lines := getLastLines(100)
		lineCount := uint(len(lines))
		if lineCount <= 0 {
			return
		}

		seen := map[string]struct{}{}
		for _, line := range lines {
			_, wasSeen := seen[line]
			if wasSeen {
				continue
			}
			fmt.Println(line)
			seen[line] = struct{}{}
		}
	}

	if flag.NArg() == 0 { // Handle no arg
		fmt.Println(getLastLine(*prev))
		return
	} else { // Handle 1 arg
		arg := flag.Arg(0)
		arg = parsePath(arg)
		_, err = os.Stat(arg)
		if err != nil {
			doFuzzyList()
			return
		}

		handle(err, "Path is not a file or directory: %s", arg)
		cache.Seek(0, io.SeekEnd)
		_, err = cache.WriteString(arg + "\n")
		handle(err, "Failed writing to cache")
		fmt.Println(arg)
	}
}
