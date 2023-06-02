package main

import (
	"Super-X/conf"
	"flag"
	"fmt"
	peparser "github.com/saferwall/pe"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Options:")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var files string
	var types string

	flag.StringVar(&files, "f", "", "DLL 文件名字")
	flag.StringVar(&types, "t", "go", "输出go源码类型或者cpp类型")
	flag.Usage = usage
	flag.Parse()

	if files == "" || types == "0" {
		usage()
	}
	conf.Dllexport(files, types)

}

// 参数：1 dll文件 ，2 写出类型go or cpp
func Dllexport(filename, mod string) {
	pe, err := peparser.New(filename, &peparser.Options{})
	if err != nil {
		log.Fatalf("Error while opening file: %s, reason: %v", filename, err)
	}

	err = pe.Parse()
	if err != nil {
		log.Fatalf("Error while parsing file: %s, reason: %v", filename, err)
	}

	var sb strings.Builder
	for _, i := range pe.Export.Functions {
		if len(i.Name) == 0 {
			continue
		}
		if i.Name[0] == '_' || i.Name[0] == '?' {
			continue
		}
		if i.Name[0:2] == "??" {
			continue
		}

		//判断写出模版类型
		if mod == "cpp" {
			str := `extern "C" __declspec(dllexport) void ${func}() {}`
			str = strings.ReplaceAll(str, "${func}", i.Name)
			// 写出模版
			sb.WriteString(str + "\n")
		} else {
			str := `//export ${func}
func ${func}(){}`
			str = strings.ReplaceAll(str, "${func}", i.Name)
			// 写出模版
			sb.WriteString(str + "\n")
		}

	}

	err = ioutil.WriteFile("DLL.txt", []byte(sb.String()), 0644)
	if err != nil {
		log.Fatalf("Error while writing file: %v", err)
	}
}
