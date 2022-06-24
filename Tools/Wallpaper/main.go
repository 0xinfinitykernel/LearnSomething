package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	log.Println("Starting...")
	flag.Parse()

	args := strings.Join(flag.Args()[:1], " ")
	println("文件名称及路径", args)

	fin, err := os.Open(args)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	argsNew := args + ".mp4"
	fout, err := os.Create(argsNew)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	// Offset is the number of bytes you want to exclude
	_, err = fin.Seek(2, io.SeekStart)
	if err != nil {
		panic(err)
	}

	n, err := io.Copy(fout, fin)
	// find . -mindepth 1 -maxdepth 1 -type f -print -exec go run /Users/xxx/Dropbox/Code/Tools/Wallpaper/main.go {} \;
	// 批量处理所有文件
	fmt.Printf("Copied %d bytes, err: %v \n", n, err)
	fmt.Println("文件输出路径: ", argsNew)
}
