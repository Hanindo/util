package util

import (
	"fmt"
	"os"
)

func MkdirP(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, 0755)
		} else {
			return err
		}
	}

	if !fi.IsDir() {
		return fmt.Errorf("%s is not directory", dir)
	}
	return nil
}
