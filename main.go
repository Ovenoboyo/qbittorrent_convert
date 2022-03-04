package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeebo/bencode"
)

type torrentMap = map[string]interface{}

func main() {
	directory := flag.String("dir", "", "Directory to scan for .fastresume files")
	originalPath := flag.String("ogPath", "", "Original path which is to be replaced")
	newPath := flag.String("newPath", "", "New path which to which original path will be changed")
	toLinux := flag.Bool("linux", true, "Should be true if path is to be converted to linux format, False for windows")

	flag.Parse()

	fmt.Println(*directory, *originalPath, *newPath, *toLinux)

	if *directory != "" && *originalPath != "" && *newPath != "" {
		readDir(*directory, func(file string) {
			writeBencode(file, replacePaths(readBencode(file), *originalPath, *newPath, *toLinux))
		})
	}
}

func readDir(directory string, callback func(file string)) {
	absPath, err := filepath.Abs(directory)
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		callback(filepath.Join(absPath, file.Name()))
	}
}

func readBencode(file string) torrentMap {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var torrent torrentMap

	bencode.DecodeBytes(data, &torrent)

	return torrent
}

func writeBencode(file string, torrent torrentMap) {
	data, err := bencode.EncodeBytes(torrent)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(file, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func isString(data interface{}) bool {
	switch data.(type) {
	case string:
		return true
	}
	return false
}

func convertToLinux(data string) string {
	return strings.ReplaceAll(data, "\\", "/")
}

func convertToWindows(data string) string {
	return strings.ReplaceAll(data, "/", "\\")
}

func fixPath(data string, oldPath string, newPath string, toLinux bool) string {
	tmp := strings.ReplaceAll(data, oldPath, newPath)

	if toLinux {
		return convertToLinux(tmp)
	}

	return convertToWindows(tmp)
}

func replacePaths(torrent torrentMap, originalPath string, newPath string, toLinux bool) torrentMap {

	qBtSavePath := torrent["qBt-savePath"]
	if isString(qBtSavePath) {
		torrent["qBt-savePath"] = fixPath(qBtSavePath.(string), originalPath, newPath, toLinux)
	}

	savePath := torrent["save_path"]
	if isString(savePath) {
		torrent["save_path"] = fixPath(savePath.(string), originalPath, newPath, toLinux)
	}

	return torrent
}
