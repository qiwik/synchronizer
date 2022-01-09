package copytest

import (
	"github.com/qiwik/synchronizer/internal/filestruct"
	"github.com/qiwik/synchronizer/pkg/logger"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

var (
	openPath             = "in.txt"
	copyPath             = "out.txt"
	mode     fs.FileMode = 0777
)

var f *logger.Logger

func init() {
	logF, _ := logger.LogFileInit()
	f = logger.LogInit(logF)
}

func BenchmarkFileCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filestruct.CopyFile(openPath, copyPath, f, mode)
		os.Remove(copyPath)
	}
}

func BenchmarkMakeDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filestruct.MakeDir("test", f, mode)
		os.Remove("test")
	}
}

func TestMD5Hash(t *testing.T) {
	req := require.New(t)
	hash, err := filestruct.GetMD5Hash("in.txt", f)
	req.NoError(err)
	req.NotEqual("", hash)
}

func TestMD5HashErr(t *testing.T) {
	req := require.New(t)
	hash, err := filestruct.GetMD5Hash("", f)
	req.Error(err)
	req.Equal("", hash)
}

func TestMakeDir(t *testing.T) {
	req := require.New(t)
	err := filestruct.MakeDir("testing", f, mode)
	req.DirExists("testing")
	req.NoError(err)
	os.Remove("testing")
}

func TestMakeDirErr(t *testing.T) {
	req := require.New(t)
	err := filestruct.MakeDir("", f, mode)
	req.Error(err)
}

func TestCopyFile(t *testing.T) {
	req := require.New(t)
	err := filestruct.CopyFile("in.txt", "out.txt", f, mode)
	req.FileExists("out.txt")
	req.NoError(err)
	os.Remove("out.txt")
}

func TestCopyFileErrOpen(t *testing.T) {
	req := require.New(t)
	err := filestruct.CopyFile("", "out.txt", f, mode)
	req.Error(err)
}

func TestCopyFileErrCreate(t *testing.T) {
	req := require.New(t)
	err := filestruct.CopyFile("in.txt", "", f, mode)
	req.Error(err)
}

func TestLog(t *testing.T) {
	req := require.New(t)
	_, err := logger.LogFileInit()
	req.FileExists("log.txt")
	req.NoError(err)
	os.Remove("log.txt")
}
