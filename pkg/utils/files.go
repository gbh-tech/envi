package utils

import (
	"errors"
	"io/fs"
	"os"

	"github.com/charmbracelet/log"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, fs.ErrNotExist) {
		log.Warnf("File %s does not exist", path)
		return false
	}

	log.Warnf("File %s does not exist", path)
	return false
}
