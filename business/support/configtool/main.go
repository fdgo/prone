package main

import (
	"business/support/config/decrypt"
	"business/support/libraries/loggers"
	"flag"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//./configtool.exe -type=encrypt -key="1234567890abcdef" -file="d:\\go\\src\\business\\support\\config\\config.yml"

func main() {
	t := flag.String("type", "", "encrypt config.yml")
	key := flag.String("key", "", "key str")
	file := flag.String("file", "", "input file path")
	flag.Parse()

	if t == nil || key == nil || file == nil {
		loggers.Error.Println("invalid params")
		return
	}

	r, e := PathExists(*file)
	if !r {
		loggers.Error.Println("file not exist ", e)
	}

	if *t == "encrypt" {
		r, e := decrypt.EncryptConfig(*file, *key)
		if e != nil || !r {
			loggers.Error.Println("encrypt ", e)
		}
		loggers.Debug.Printf("encrypt success")
	} else if *t == "decrypt" {
		r, e := decrypt.DecryptConfig(*file, *key)
		if e != nil || !r {
			loggers.Error.Println("decrypt ", e)
		}

		loggers.Debug.Printf("decrypt success")
	}
}
