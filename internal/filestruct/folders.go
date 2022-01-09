package filestruct

import (
	"context"
	"errors"
	"github.com/qiwik/synchronizer/internal/parameters/initial"
	"github.com/qiwik/synchronizer/pkg/logger"
	"io/fs"
	"os"
)

// FoldersInfo consists of the information about passed embedded folders and files and folders in it
type FoldersInfo struct {
	Folders []*FoldersInfo
	Files   []*FilesInfo
	Name    string
	Exist   bool
	Mode    fs.FileMode
}

// MainFolder consists of the information about folders and files in it
type MainFolder struct {
	Folders []*FoldersInfo
	Files   []*FilesInfo
}

// FoldersSearch compares folders in main and copy folders, and performs actions on them
func (m *MainFolder) FoldersSearch(param initial.Parameters, copied *MainFolder, logs *logger.Logger, ctxs context.Context) {
	ctx, cancel := context.WithCancel(ctxs)
	go func() {
		ans := m.FilesSearch(param, copied, logs)
		if ans {
			cancel()
		}
	}()

	if len(m.Folders) == 0 {
		return
	}

	for i := range m.Folders {
		exist := false
		for j := range copied.Folders {
			if m.Folders[i].Name == copied.Folders[j].Name {
				exist = true
				copied.Folders[j].Exist = true
				m.Folders[i].FoldersSearch(param.SourcePath+"/"+m.Folders[i].Name, param.CopyPath+"/"+m.Folders[i].Name,
					copied.Folders[j], logs, ctxs)
			}
		}

		if exist != true {
			copied.Folders = append(copied.Folders, &FoldersInfo{
				Name:  m.Folders[i].Name,
				Exist: true,
			})

			err := MakeDir(param.CopyPath+"/"+m.Folders[i].Name, logs, m.Folders[i].Mode)
			if err != nil {
				logs.Errs.Println(err)
			}

			if len(m.Folders[i].Folders) != 0 {
				m.Folders[i].FoldersSearch(param.SourcePath+"/"+m.Folders[i].Name,
					param.CopyPath+"/"+m.Folders[i].Name, copied.Folders[len(copied.Folders)-1], logs, ctxs)
			}
		}
	}

	<-ctx.Done()
	DeleteTree(param.CopyPath, copied.Folders, logs)
}

//FoldersSearch compares embedded folders of the main and copy folders, and performs actions on them
func (f *FoldersInfo) FoldersSearch(openPath, createPath string, copied *FoldersInfo, logs *logger.Logger, ctxs context.Context) {
	ctx, cancel := context.WithCancel(ctxs)
	go func() {
		ans := f.FilesSearch(openPath, createPath, copied, logs)
		if ans {
			cancel()
		}
	}()

	for i := range f.Folders {
		exist := false
		for j := range copied.Folders {
			if f.Folders[i].Name == copied.Folders[j].Name {
				exist = true
				copied.Folders[j].Exist = true
				f.Folders[i].FoldersSearch(openPath+"/"+f.Folders[i].Name, createPath+"/"+f.Folders[i].Name,
					copied.Folders[j], logs, ctxs)
			}
		}

		if exist != true {
			copied.Folders = append(copied.Folders, &FoldersInfo{
				Name:  f.Folders[i].Name,
				Exist: true,
			})

			err := MakeDir(createPath+"/"+f.Folders[i].Name, logs, f.Folders[i].Mode)
			if err != nil {
				logs.Errs.Println(err)
			}

			if len(f.Folders[i].Files) != 0 {
				f.Folders[i].FilesSearch(openPath+"/"+f.Folders[i].Name,
					createPath+"/"+f.Folders[i].Name, copied.Folders[len(copied.Folders)-1], logs)
			}

			if len(f.Folders[i].Folders) != 0 {
				f.Folders[i].FoldersSearch(openPath+"/"+f.Folders[i].Name,
					createPath+"/"+f.Folders[i].Name, copied.Folders[len(copied.Folders)-1], logs, ctxs)
			}
		}
	}

	<-ctx.Done()
	DeleteTree(createPath, copied.Folders, logs)
}

// MakeDir creates directory in the copy folder if don't exist
func MakeDir(dirName string, logs *logger.Logger, mode fs.FileMode) error {
	err := os.Mkdir(dirName, mode)
	if err != nil {
		return errors.New("can't do make directory operation with " + dirName)
	}
	logs.Info.Printf("directory %s created successfully\n", dirName)
	return nil
}

// DeleteTree deletes folders, embedded folders and files if they exist only in a copy folder
func DeleteTree(createPath string, copied []*FoldersInfo, logs *logger.Logger) {
	for k := range copied {
		if copied[k].Exist == false {
			if len(copied[k].Files) != 0 {
				for i := range copied[k].Files {
					err := os.Remove(createPath + "/" + copied[k].Name + "/" + copied[k].Files[i].Name)
					if err != nil {
						logs.Errs.Printf("can't remove %s file",
							createPath+"/"+copied[k].Name+"/"+copied[k].Files[i].Name)
						break
					}
					logs.Info.Printf("file %s successfully removed\n",
						createPath+"/"+copied[k].Name+"/"+copied[k].Files[i].Name)
				}
			}

			if len(copied[k].Folders) != 0 {
				DeleteTree(createPath+"/"+copied[k].Name, copied[k].Folders, logs)
			}

			err := os.Remove(createPath + "/" + copied[k].Name)
			if err != nil {
				logs.Errs.Printf("can't remove %s folder", createPath+"/"+copied[k].Name)
				break
			}
			logs.Info.Printf("%s was successfully removed\n", createPath+"/"+copied[k].Name)
		}
	}
}
