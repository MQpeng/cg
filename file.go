package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// As return as is
func As(s string) string {
	return s
}

// CheckPathExists used for checking the path exist
func CheckPathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return true
}

// MakeDirIfNotExist used for making dir if the dir is not exist
func MakeDirIfNotExist(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// CopyFile is copy file to target dir
func CopyFile(from, to string) error {
	return CopyFileWithFunc(from, to, As, io.Copy)
}

// CopyFileWithFunc is copy file to target dir
func CopyFileWithFunc(from, to string, as func(string) string, copy func(dst io.Writer, src io.Reader) (written int64, err error)) error {
	srcFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(as(to))
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

// CopyDir used for copy dir
func CopyDir(from, to string) error {
	return CopyDirWithFunc(from, to, As, io.Copy, nil)
}

// CopyDirWithFunc used for copy dir with name func
func CopyDirWithFunc(from, to string, 
	as func(string) string, 
	copy func(dst io.Writer, src io.Reader) (written int64, err error),
	exclude func(path string) bool,
	) error {
	if !CheckPathExists(from) {
		return fmt.Errorf("from path must be exist:[%s]", from)
	}
	srcDir, err := os.Open(from)
	if err != nil {
		return err
	}
	defer srcDir.Close()
	err = MakeDirIfNotExist(to)
	if err != nil {
		return err
	}
	fileInfos, err := srcDir.Readdir(-1)
	if err != nil {
		return err
	}
	excludeExist := exclude != nil
	for _, fileInfo := range fileInfos {
		fileName := fileInfo.Name()
		if excludeExist {
			if exclude(fileName) {
				continue
			}
		}
		srcPath := filepath.Join(from, fileName)
		dstPath := filepath.Join(to, as(fileName))
		if fileInfo.IsDir() {
			err = CopyDirWithFunc(srcPath, dstPath, as, copy, exclude)
		} else {
			err = CopyFileWithFunc(srcPath, dstPath, as, copy)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
