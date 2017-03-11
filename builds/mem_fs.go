package builds

import (
	"fmt"
	"io"
	"path"
	"sort"
)

func pathDir(p string) string {
	ret := path.Dir(p)
	if ret == "." {
		return ""
	}
	return ret
}

// MemFS is a in memory file system.
type MemFS struct {
	dirs map[string]*memDir
}

// NewMemFS creates an empty memory file system.
func NewMemFS() *MemFS {
	dirs := make(map[string]*memDir)
	dirs[""] = newMemDir("")
	return &MemFS{
		dirs: dirs,
	}
}

// MakeDir makes an empty directory if it does not exist yet.
func (fs *MemFS) MakeDir(p string) error {
	if p == "" {
		return nil
	}
	if err := checkValidPath(p); err != nil {
		return err
	}
	for p != "" {
		_, ok := fs.dirs[p]
		if ok {
			return nil
		}
		fs.dirs[p] = newMemDir(p)
		p = pathDir(p)
	}
	return nil
}

// HasDir checks if a directory exists.
func (fs *MemFS) HasDir(p string) (bool, error) {
	if err := checkValidDir(p); err != nil {
		return false, err
	}
	_, ok := fs.dirs[p]
	return ok, nil
}

// Open opens a file for reading.
func (fs *MemFS) Open(p string) (*File, error) {
	if err := checkValidPath(p); err != nil {
		return nil, err
	}
	dir := fs.dirs[pathDir(p)]
	if dir == nil {
		return nil, fmt.Errorf("file %q not found", p)
	}
	name := path.Base(p)
	f := dir.open(name)
	if f == nil {
		return nil, fmt.Errorf("file %q not found", p)
	}
	return &File{
		Path:   p,
		Name:   name,
		Opener: f.Opener(),
	}, nil
}

// Create creates a memory file for writing.
func (fs *MemFS) Create(p string) (io.WriteCloser, error) {
	if err := checkValidPath(p); err != nil {
		return nil, err
	}
	dir := pathDir(p)
	if err := fs.MakeDir(dir); err != nil {
		return nil, err
	}

	mdir := fs.dirs[dir]
	return mdir.create(path.Base(p)), nil
}

// AddFile adds a file into the system.
func (fs *MemFS) AddFile(p string, bs []byte) error {
	wc, err := fs.Create(p)
	if err != nil {
		return err
	}
	defer wc.Close()
	if _, err := wc.Write(bs); err != nil {
		return err
	}
	return wc.Close()
}

// AddTextFile adds a text file into the system.
func (fs *MemFS) AddTextFile(p, s string) error {
	return fs.AddFile(p, []byte(s))
}

// ListFiles lists all the files in a
func (fs *MemFS) ListFiles(p string) ([]string, error) {
	if err := checkValidPath(p); err != nil {
		return nil, err
	}
	dir, ok := fs.dirs[p]
	if !ok {
		return nil, fmt.Errorf("directory %q not exist", p)
	}
	return dir.list(), nil
}

// ListDirs lists all sub directories of a directory
func (fs *MemFS) ListDirs(p string) ([]string, error) {
	if err := checkValidDir(p); err != nil {
		return nil, err
	}
	if fs.dirs[p] == nil {
		return nil, fmt.Errorf("directory %q not exist", p)
	}

	var ret []string
	for dir := range fs.dirs {
		if p == pathDir(dir) {
			ret = append(ret, dir)
		}
	}
	sort.Strings(ret)
	return ret, nil
}
