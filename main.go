package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vallerion/rscanner"
	flag "github.com/spf13/pflag"
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
	cache_path := parsePath(cacheFileName)

	list := flag.UintP("list", "l", 0, "Will list out the most recent paths.")
	flag.Lookup("list").NoOptDefVal = "10"
	prevHelp := `Will get the n-th previous path. 
Using -n or --n where n is a number will also work.`
	prev := flag.UintP("previous", "p", 0, prevHelp)

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

	cache, err := os.OpenFile(cache_path, os.O_CREATE|os.O_RDWR, 0644)
	handle(err, "Failed to get or create cache at '%s'", cache_path)
	cacheStat, err := cache.Stat()
	handle(err, "Failed to stat cache")

	getlastLine := func(num uint) string {
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

	getlastLines := func(num uint) []string {
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

	if list != nil && *list > 0 {
		for i, line := range getlastLines(*list) {
			fmt.Printf("%d: %s\n", i, line)
		}
		return
	}

	if flag.NArg() == 0 {
		fmt.Println(getlastLine(*prev))
		return
	} else {
		arg := parsePath(flag.Arg(0))
		_, err = os.Stat(arg)
		handle(err, "Path is not a file or directory: %s", arg)
		cache.Seek(0, io.SeekEnd)
		_, err = cache.WriteString(arg + "\n")
		handle(err, "Failed writing to cache")
		fmt.Println(arg)
	}
}
