package g

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	Basedirs   []string
	CurrentDir string
	TempDir    string = "tmp"
	Product    bool
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	Basedirs = append(Basedirs, filepath.Dir(filepath.Dir(file)))
	dir, _ := os.Getwd()
	Basedirs = append(Basedirs, dir)

	_, Product = os.LookupEnv("product")
	if Product {
		file, _ = filepath.Abs(os.Args[0])
		CurrentDir = filepath.Dir(file)
	}

	os.Chdir(CurrentDir)

}
