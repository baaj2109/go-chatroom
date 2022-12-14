package global

import (
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var RootDir string

var once = new(sync.Once)

func Init() {
	once.Do(func() {
		inferRootDir()
		// initConfig()
	})
}

// / 找出項目根目錄
func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	/// recursive call infer 判斷目錄 d 下面使否存在 template
	var infer func(d string) string
	infer = func(d string) string {
		/// 確定根目錄下有 template
		if exists(d + "/template") {
			return d
		}
		return infer(filepath.Dir(d))
	}
	RootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
