package scan

import (
	"io/ioutil"
	"synchronizer/internal/filestruct"
	"synchronizer/pkg/logger"
)

//amount of bytes in gigabyte
const gigabyte int64 = 1076741824

//delta compensates a time error. secPerGb - average time for coping one GB. Every parameter in seconds
var (
	delta int64 = 60
	secPerGb int64 = 20
)

//PathScan scans embedded folders and files in them for creating a file tree. Returns an error
func PathScan(path string, tree *filestruct.FoldersInfo, logs *logger.Logger) error {
	elems, err := ioutil.ReadDir(path)
	if err != nil {
		logs.Errs.Printf("current path %s was not read correctly", path)
	}

	for _, elem := range elems {
		if elem.IsDir() {
			tree.Folders = append(tree.Folders, &filestruct.FoldersInfo{
				Name: elem.Name(),
				Mode: elem.Mode(),
			})

			err = PathScan(path + "/" + elem.Name(), tree.Folders[len(tree.Folders)-1], logs)
			if err != nil {
				logs.Errs.Printf("current path %s was not read correctly", path + "/" + elem.Name())
			}
		} else {
			currentHash, err := filestruct.GetMD5Hash(path + "/" + elem.Name(), logs)
			if err != nil {
				logs.Fatal.Fatal(err)
			}
			tree.Files = append(tree.Files, &filestruct.FilesInfo{
				Name: elem.Name(),
				Hash: currentHash,
				Mode: elem.Mode(),
			})

		}
	}
	logs.Info.Printf("%s read correctly\n", path)
	return nil
}

//MainPathScan scans main folder and files in it for creating upper-level branch of a file tree. Returns an error and
//an amount of ticker based on the weight of the main folder
func MainPathScan(path string, tree *filestruct.MainFolder, logs *logger.Logger) (int64, error) {
	var mainSize int64

	elems, err := ioutil.ReadDir(path)
	if err != nil {
		logs.Errs.Fatalf("current path %s was not read correctly", path)
	}

	for _, elem := range elems {
			if elem.IsDir() {
				tree.Folders = append(tree.Folders, &filestruct.FoldersInfo{
					Name: elem.Name(),
					Mode: elem.Mode(),
				})

				err = PathScan(path+"/"+elem.Name(), tree.Folders[len(tree.Folders)-1], logs)
				if err != nil {
					logs.Errs.Printf("current path %s was not read correctly", path+"/"+elem.Name())
				}
				mainSize += elem.Size()
			} else {
				currentHash, err := filestruct.GetMD5Hash(path + "/" + elem.Name(), logs)
				if err != nil {
					logs.Fatal.Fatal(err)
				}
				tree.Files = append(tree.Files, &filestruct.FilesInfo{
					Name: elem.Name(),
					Hash: currentHash,
					Mode: elem.Mode(),
				})
				mainSize += elem.Size()
			}
	}
	logs.Info.Printf("work with %s finished\n", path)
	ticker := (mainSize/gigabyte) * secPerGb + delta
	return ticker, nil
}
