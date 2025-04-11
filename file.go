package ossfs

import (
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"

	"github.com/spf13/afero"
)

type File struct {
	name     string
	fs       *Fs
	openFlag int
	offset   int64
	isDir    bool

	// Whether the file is written.
	dirty bool

	// Whether the file is closed.
	closed bool

	// Whether the file is preloaded downto local.
	preloaded   bool
	preloadedFd afero.File

	mu sync.RWMutex
}

func NewOssFile(name string, flag int, fs *Fs) (*File, error) {
	return &File{
		name:        fs.normFileName(name),
		fs:          fs,
		openFlag:    flag,
		offset:      0,
		dirty:       false,
		closed:      false,
		isDir:       fs.isDir(fs.normFileName(name)),
		preloaded:   false,
		preloadedFd: nil,
	}, nil
}

func (f *File) preload() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	pfs := f.fs.preloadFs
	if _, err := pfs.Stat(f.name); err == nil {
		if e := pfs.Remove(f.name); e != nil {
			return e
		}
	}
	pfd, err := f.fs.preloadFs.Create(f.name)
	if err != nil {
		return err
	}

	r, clean, e := f.fs.manager.GetObject(f.fs.ctx, f.fs.bucketName, f.name)
	if e != nil {
		return e
	}
	defer clean()

	if _, err := io.Copy(pfd, r); err != nil {
		return err
	}

	if _, err := pfd.Seek(f.offset, io.SeekStart); err != nil {
		return err
	}

	f.preloadedFd = pfd
	f.preloaded = true
	return nil
}

func (f *File) getFileInfo() (os.FileInfo, error) {
	if f.dirty {
		if f.preloadedFd == nil {
			return nil, syscall.EACCES
		}
		return f.preloadedFd.Stat()
	}
	return f.fs.Stat(f.name)
}

func (f *File) isReadable() bool {
	return !f.closed && (f.openFlag&os.O_RDONLY != 0 || f.openFlag&os.O_RDWR != 0)
}

func (f *File) isWriteable() bool {
	return !f.closed && (f.openFlag&os.O_WRONLY != 0 || f.openFlag&os.O_RDWR != 0)
}

func (f *File) isAppendOnly() bool {
	return f.isWriteable() && f.openFlag&os.O_APPEND != 0
}

func (f *File) Read(p []byte) (int, error) {
	if !f.isReadable() || f.isDir {
		return 0, syscall.EPERM
	}
	n, err := f.ReadAt(p, f.offset)
	if err != nil {
		return 0, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.offset += int64(n)
	return n, err
}

func (f *File) ReadAt(p []byte, off int64) (int, error) {
	if !f.isReadable() || f.isDir {
		return 0, syscall.EPERM
	}
	reader, cleanUp, err := f.fs.manager.GetObjectPart(f.fs.ctx, f.fs.bucketName, f.name, off, off+int64(len(p)))
	if err != nil {
		return 0, err
	}
	defer cleanUp()
	return reader.Read(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if (!f.isReadable() && !f.isWriteable()) || f.isDir {
		return 0, syscall.EPERM
	}
	fi, err := f.getFileInfo()
	if err != nil {
		return 0, err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	max := fi.Size()
	var newOffset int64
	switch whence {
	case io.SeekCurrent:
		newOffset = f.offset + offset
	case io.SeekStart:
		newOffset = offset
	case io.SeekEnd:
		newOffset = max + offset
	default:
		return 0, fmt.Errorf("invalid whence value: %v", whence)
	}
	if newOffset < 0 || newOffset > max {
		return 0, afero.ErrOutOfRange
	}
	f.offset = newOffset
	return f.offset, nil
}

func (f *File) Write(p []byte) (int, error) {
	if !f.isWriteable() {
		return 0, syscall.EPERM
	}
	if f.isAppendOnly() {
		fi, err := f.getFileInfo()
		if err != nil {
			return 0, err
		}
		return f.doWriteAt(p, fi.Size())
	}
	// f.mu.Lock()
	// defer f.mu.Unlock()

	n, e := f.doWriteAt(p, f.offset)
	if e != nil {
		return 0, e
	}
	f.offset += int64(n)
	return n, e
}

func (f *File) doWriteAt(p []byte, off int64) (int, error) {
	if f.isDir {
		return 0, syscall.EPERM
	}

	if !f.preloaded {
		if err := f.preload(); err != nil {
			return 0, err
		}
	}

	n, e := f.preloadedFd.WriteAt(p, off)
	f.dirty = true
	if f.fs.autoSync {
		f.Sync()
	}
	return n, e
}

func (f *File) WriteAt(p []byte, off int64) (int, error) {
	if !f.isWriteable() || f.isAppendOnly() {
		return 0, syscall.EPERM
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.doWriteAt(p, off)
}

func (f *File) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	err := f.Sync()
	if err != nil {
		return err
	}
	delete(f.fs.openedFiles, f.name)
	if f.preloaded {
		err := f.fs.preloadFs.Remove(f.name)
		if err != nil {
			return err
		}
		err = f.preloadedFd.Close()
		if err != nil {
			return err
		}
		f.preloadedFd = nil
		f.preloaded = false
	}
	f.dirty = false
	f.closed = true
	return nil
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isReadable() {
		return nil, syscall.EPERM
	}

	fis, err := f.fs.manager.ListObjects(f.fs.ctx, f.fs.bucketName, f.fs.ensureAsDir(f.name), count)
	return fis, err
}

func (f *File) Readdirnames(n int) ([]string, error) {
	if !f.isReadable() {
		return nil, syscall.EPERM
	}

	fis, err := f.Readdir(n)
	if err != nil {
		return nil, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	var fNames []string
	for _, fi := range fis {
		fNames = append(fNames, fi.Name())
	}

	return fNames, nil
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.getFileInfo()
}

func (f *File) Sync() error {
	if f.preloaded {
		if _, err := f.fs.manager.PutObject(f.fs.ctx, f.fs.bucketName, f.name, f.preloadedFd); err != nil {
			return err
		}
	}
	return nil
}

func (f *File) Truncate(size int64) error {
	if !f.isWriteable() || f.isDir {
		return syscall.EPERM
	}
	_, err := f.WriteAt([]byte(""), 0)
	return err
}

func (f *File) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}
