package ossfs

import (
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/messikiller/afero-oss/internal/mocks"
	"github.com/messikiller/afero-oss/internal/utils"
)

func getMockedFs() *Fs {
	fs := NewOssFs("test-ak", "test-sk", "test-region", "test-bucket")
	fs.manager = &mocks.ObjectManager{}
	return fs
}

func getMockedFile(name string, flag int, fs *Fs) *File {
	f, _ := NewOssFile(name, flag, fs)
	return f
}

func TestNewOssFile(t *testing.T) {
	t.Run("create new file with read flag", func(t *testing.T) {
		fs := &Fs{}
		file, err := NewOssFile("testfile", os.O_RDONLY, fs)
		assert.NoError(t, err)
		assert.Equal(t, "testfile", file.name)
		assert.Equal(t, os.O_RDONLY, file.openFlag)
		assert.Equal(t, fs, file.fs)
		assert.False(t, file.dirty)
		assert.False(t, file.closed)
		assert.False(t, file.isDir)
		assert.False(t, file.preloaded)
		assert.Nil(t, file.preloadedFd)
	})

	t.Run("create new file with write flag", func(t *testing.T) {
		fs := &Fs{}
		file, err := NewOssFile("testfile", os.O_WRONLY, fs)
		assert.NoError(t, err)
		assert.Equal(t, "testfile", file.name)
		assert.Equal(t, os.O_WRONLY, file.openFlag)
		assert.Equal(t, fs, file.fs)
		assert.False(t, file.dirty)
		assert.False(t, file.closed)
		assert.False(t, file.isDir)
		assert.False(t, file.preloaded)
		assert.Nil(t, file.preloadedFd)
	})

	t.Run("create new directory", func(t *testing.T) {
		fs := &Fs{}
		file, err := NewOssFile("testdir/", os.O_RDONLY, fs)
		assert.NoError(t, err)
		assert.Equal(t, "testdir/", file.name)
		assert.Equal(t, os.O_RDONLY, file.openFlag)
		assert.Equal(t, fs, file.fs)
		assert.False(t, file.dirty)
		assert.False(t, file.closed)
		assert.True(t, file.isDir)
		assert.False(t, file.preloaded)
		assert.Nil(t, file.preloadedFd)
	})

	t.Run("normalize file name", func(t *testing.T) {
		fs := &Fs{}
		file, err := NewOssFile("/path/testfile", os.O_RDONLY, fs)
		assert.NoError(t, err)
		assert.Equal(t, "path/testfile", file.name)
	})
}

func TesFiletRead(t *testing.T) {
	t.Run("Read with unreadable flag return error", func(t *testing.T) {
		fs := getMockedFs()
		f := getMockedFile("testfile", os.O_WRONLY, fs)

		p := make([]byte, 0)
		_, e := f.Read(p)

		assert.Error(t, e)
		assert.NotNil(t, e)
	})

	t.Run("Read on directory return error", func(t *testing.T) {
		fs := getMockedFs()
		f := getMockedFile("testdir", os.O_RDONLY, fs)
		f.isDir = true

		p := make([]byte, 10)
		_, e := f.Read(p)

		assert.Error(t, e)
		assert.Equal(t, syscall.EPERM, e)
	})

	t.Run("Read on closed file return error", func(t *testing.T) {
		fs := getMockedFs()
		f := getMockedFile("testfile", os.O_RDONLY, fs)
		f.closed = true

		p := make([]byte, 10)
		_, e := f.Read(p)

		assert.Error(t, e)
		assert.Equal(t, syscall.EPERM, e)
	})

	t.Run("Successful read updates offset", func(t *testing.T) {
		fs := getMockedFs()
		var cu utils.CleanUp = func() {}
		mockManager := fs.manager.(*mocks.ObjectManager)
		mockManager.On("GetObjectPart", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(strings.NewReader("testdata"), cu, nil)

		f := getMockedFile("testfile", os.O_RDONLY, fs)
		p := make([]byte, 8)
		n, err := f.Read(p)

		assert.NoError(t, err)
		assert.Equal(t, 8, n)
		assert.Equal(t, int64(8), f.offset)
	})

	t.Run("ReadAt error propagates", func(t *testing.T) {
		fs := getMockedFs()
		mockManager := fs.manager.(*mocks.ObjectManager)
		mockManager.On("GetObjectPart", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil, syscall.EIO)

		f := getMockedFile("testfile", os.O_RDONLY, fs)
		p := make([]byte, 8)
		_, err := f.Read(p)

		assert.Error(t, err)
		assert.Equal(t, syscall.EIO, err)
	})
}

func TestFileReadAt(t *testing.T) {
	t.Run("ReadAt with unreadable flag return error", func(t *testing.T) {
		fs := getMockedFs()
		f := getMockedFile("testfile", os.O_WRONLY, fs)

		p := make([]byte, 0)
		_, e := f.ReadAt(p, 0)

		assert.Error(t, e)
		assert.NotNil(t, e)
	})

	t.Run("ReadAt on dir return error", func(t *testing.T) {
		fs := getMockedFs()
		f := getMockedFile("/path/to/dir/", os.O_WRONLY, fs)

		p := make([]byte, 0)
		_, e := f.ReadAt(p, 0)

		assert.Error(t, e)
		assert.NotNil(t, e)
	})

	t.Run("ReadAt success", func(t *testing.T) {
		fs := getMockedFs()
		f := getMockedFile("testfile", os.O_RDONLY, fs)

		p := make([]byte, 4)

		var cu utils.CleanUp = func() {}
		off := int64(5)
		m := &mocks.ObjectManager{}
		m.
			On("GetObjectPart", f.fs.ctx, f.fs.bucketName, f.name, off, off+int64(len(p))).
			Return(strings.NewReader("test result"), cu, nil)
		fs.manager = m

		n, e := f.ReadAt(p, off)

		assert.Nil(t, e)
		assert.Equal(t, 4, n)
		assert.Equal(t, "test", string(p))
	})
}
