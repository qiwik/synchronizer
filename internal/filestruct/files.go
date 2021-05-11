package filestruct

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"synchronizer/internal/parameters/initial"
	"synchronizer/pkg/logger"
)

// FilesInfo consists of the information about passed file
type FilesInfo struct {
	Name string
	Exist bool
	Hash string
	Mode fs.FileMode
}

// FilesSearch compares files in the main folder and copy folder and performs actions on them.
// Returns true if work was done
func (m *MainFolder) FilesSearch(param initial.Parameters, copied *MainFolder, logs *logger.Logger) bool{
	if len(m.Files) == 0 {
		return true
	}

	for i := range m.Files {
		exist := false
		for j := range copied.Files {
			if m.Files[i].Name == copied.Files[j].Name && m.Files[i].Hash == copied.Files[j].Hash {
				exist = true
				copied.Files[j].Exist = true
				logs.Info.Printf("file %s didn't change\n", copied.Files[j].Name)
			}

			if m.Files[i].Name == copied.Files[j].Name && m.Files[i].Hash != copied.Files[j].Hash {
				exist = true
				err := os.Remove(param.CopyPath + "/" + m.Files[i].Name)
				if err != nil {
					logs.Errs.Printf("can't replace file %s\n", param.SourcePath + "/" + m.Files[i].Name)
					break
				}

				err = CopyFile(param.SourcePath + "/" + m.Files[i].Name, param.CopyPath + "/" + m.Files[i].Name, logs, m.Files[i].Mode)
				if err != nil {
					logs.Errs.Println(err)
				}
				copied.Files[j].Exist = true
			}
		}

		if exist != true {
			copied.Files = append(copied.Files, &FilesInfo{
				Name: m.Files[i].Name,
				Exist: true,
			})

			err := CopyFile(param.SourcePath + "/" + m.Files[i].Name, param.CopyPath + "/" + m.Files[i].Name, logs, m.Files[i].Mode)
			if err != nil {
				logs.Errs.Println(err)
			}
		}
	}
	DeleteFiles(param.CopyPath, copied.Files, logs)
	return true
}

// FilesSearch compares files in the embedded folders of main and copy folders and performs actions on them.
// Returns true if work was done
func (f *FoldersInfo) FilesSearch(openPath, createPath string, copied *FoldersInfo, logs *logger.Logger) bool {
	for i := range f.Files {
		exist := false
		for j := range copied.Files {
			if f.Files[i].Name == copied.Files[j].Name && f.Files[i].Hash == copied.Files[j].Hash {
				exist = true
				copied.Files[j].Exist = true
				logs.Info.Printf("file %s didn't change\n", copied.Files[j].Name)
			}

			if f.Files[i].Name == copied.Files[j].Name && f.Files[i].Hash != copied.Files[j].Hash {
				exist = true
				err := os.Remove(createPath + "/" + f.Files[i].Name)
				if err != nil {
					logs.Errs.Printf("can't replace file %s\n", createPath + "/" + f.Files[i].Name)
					break
				}

				err = CopyFile(openPath + "/" + f.Files[i].Name, createPath + "/" + f.Files[i].Name, logs, f.Files[i].Mode)
				if err != nil {
					logs.Errs.Println(err)
				}
				copied.Files[j].Exist = true
			}
		}

		if exist != true {
			copied.Files = append(copied.Files, &FilesInfo{
				Name: f.Files[i].Name,
				Exist: true,
			})

			err := CopyFile(openPath + "/" + f.Files[i].Name, createPath + "/" + f.Files[i].Name, logs, f.Files[i].Mode)
			if err != nil {
				logs.Errs.Println(err)
			}
		}
	}
	DeleteFiles(createPath, copied.Files, logs)
	return true
}

// CopyFile creates file in the copy folder
func CopyFile(openPath, createPath string, logs *logger.Logger, mode fs.FileMode) error {
	file, err := ioutil.ReadFile(openPath)
	if err != nil {
		return errors.New("can't open file" + openPath)
	}
	err = ioutil.WriteFile(createPath, file, mode)
	if err != nil {
		return errors.New("can't create file " + createPath)
	}
	logs.Info.Printf("file %s successfully copied\n", createPath)
	return nil
}

// DeleteFiles deletes file in the copy folder
func DeleteFiles(path string, files []*FilesInfo, logs *logger.Logger) {
	for f := range files {
		if files[f].Exist != true {
			err := os.Remove(path + "/" + files[f].Name)
			if err != nil {
				logs.Errs.Printf("can't remove file %s", path + "/" + files[f].Name)
				return
			}
			logs.Info.Printf("successfully removed file %s\n", path + "/" + files[f].Name)
		}
	}
}

// GetMD5Hash get hash sum of the files for tracking changes
func GetMD5Hash(name string, logs *logger.Logger) (string, error) {
	file, err := os.Open(name)
	if err != nil {
		return "", errors.New("can't open file" + name + "for create hash sum")
	}
	fileSum := md5.New()
	_, err = io.Copy(fileSum, file)
	if err != nil {
		return "", errors.New("can't copy " + name + " file's data to hash sum")
	}
	logs.Info.Printf("successfully created hash sum for %s\n", name)
	defer func() {
		err = file.Close()
		if err != nil {
			logs.Errs.Fatalf("can't close file %s\n", name)
		}
	}()
	return fmt.Sprintf("%X", fileSum.Sum(nil)), nil
}