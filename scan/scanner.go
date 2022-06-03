package scanner

import (
	"predorock/gitscan/logstreamer"
	"predorock/gitscan/utils"
	"sync"
)

var log = logstreamer.GetInstance().Log

func Scan(folder string, op func(string)) {
	if utils.IsGitFolder(folder) {
		op(folder)
		return
	}

	var folders, err = utils.GetDirs(folder)

	if err != nil {
		log.Fatalf("Error reading folder: %s - %s\n", folder, err.Error())
		return
	}

	for _, f := range folders {
		Scan(folder+"/"+f, op)
	}

}

func ScanParallel(folder string, op func(string)) {
	if utils.IsGitFolder(folder) {
		op(folder)
		return
	}

	var folders, err = utils.GetDirs(folder)

	if err != nil {
		log.Fatalf("Error reading folder: %s - %s\n", folder, err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(folders))

	for _, f := range folders {
		var parent = folder
		var current = f
		var fn = op
		go func() {
			ScanParallel(parent+"/"+current, fn)
			wg.Done()
		}()
	}
	wg.Wait()
}
