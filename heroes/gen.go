package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func main() {
	var srcFiles, dstFiles []string
	err := filepath.Walk("../images/heroes/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		srcFiles = append(srcFiles, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range srcFiles {
		// separate path/filename
		idx := strings.LastIndex(file, "/")
		name := file[idx+1:]

		// track errata/CK
		errata := false
		if strings.LastIndex(name, "-errata.png") > 0 {
			name = name[:strings.LastIndex(name, "-")] + ".png"
			errata = true
		}
		ck := false
		if strings.LastIndex(name, "-ck.png") > 0 {
			ck = true
		}

		// chop off expansion
		name = name[:strings.LastIndex(name, "-")]
		if ck {
			name += "CK"
		}
		if errata {
			name += "Errata"
		}
		name += ".png"

		// title case
		tmp := []byte(name)
		for j := range tmp {
			if j == 0 || tmp[j-1] == '-' {
				tmp[j] -= 0x20 // capitalize
			}
		}
		name = string(tmp)

		// save
		dstFiles = append(dstFiles, "./"+strings.ReplaceAll(name, "-", ""))
	}

	for i := range dstFiles {
		Copy(srcFiles[i], dstFiles[i])
	}
}
