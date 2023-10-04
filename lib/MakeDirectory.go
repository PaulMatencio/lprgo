package lib

import "os"

func CreateDir(directory string) (err error) {
	if _, err = os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, os.ModePerm)
	}
	return
}
