package ops

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/abdfnx/renio/core/options"
	"github.com/abdfnx/renio/tools"

	quickjs "github.com/abdfnx/qjs"
	"github.com/spf13/afero"
)

type FsDriver struct {
	Perms *options.Perms
	Fs    afero.Fs
}

var _ io.Reader = (*os.File)(nil)

func (fs *FsDriver) ReadFile(ctx *quickjs.Context, path quickjs.Value) quickjs.Value {
	data, err := afero.ReadFile(fs.Fs, path.String())
	tools.Check(err)
	return ctx.String(string(data))
}

func (fs *FsDriver) WriteFile(ctx *quickjs.Context, path quickjs.Value, content quickjs.Value) quickjs.Value {
	err := afero.WriteFile(fs.Fs, path.String(), []byte(content.String()), 0777)
	tools.Check(err)
	return ctx.Bool(true)
}

func (fs *FsDriver) Exists(ctx *quickjs.Context, path quickjs.Value) quickjs.Value {
	data, err := afero.Exists(fs.Fs, path.String())
	tools.Check(err)
	return ctx.Bool(data)
}

func (fs *FsDriver) DirExists(ctx *quickjs.Context, path quickjs.Value) quickjs.Value {
	data, err := afero.DirExists(fs.Fs, path.String())
	tools.Check(err)
	return ctx.Bool(data)
}

func (fs *FsDriver) Cwd(ctx *quickjs.Context) quickjs.Value {
	dir, err := os.Getwd()
	tools.Check(err)
	return ctx.String(dir)
}

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

func (fs *FsDriver) Stat(ctx *quickjs.Context, path quickjs.Value) quickjs.Value {
	entry, err := fs.Fs.Stat(path.String())
	tools.Check(err)
	f := FileInfo{
		Name:    entry.Name(),
		Size:    entry.Size(),
		Mode:    entry.Mode(),
		ModTime: entry.ModTime(),
		IsDir:   entry.IsDir(),
	}

	output, err := json.Marshal(f)
	tools.Check(err)
	return ctx.String(string(output))
}

func (fs *FsDriver) Remove(ctx *quickjs.Context, path quickjs.Value) quickjs.Value {
	err := fs.Fs.Remove(path.String())
	tools.Check(err)
	return ctx.Bool(true)
}

func (fs *FsDriver) Mkdir(ctx *quickjs.Context, path quickjs.Value) quickjs.Value {
	err := fs.Fs.Mkdir(path.String(), os.FileMode(0777))
	tools.Check(err)
	return ctx.Bool(true)
}

type walkFs struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
	Path    string
}

func (fs *FsDriver) Walk(ctx *quickjs.Context, pathDir quickjs.Value) quickjs.Value {
	var files []walkFs

	err := filepath.Walk(pathDir.String(), func(path string, info os.FileInfo, err error) error {
		data := walkFs{
			Name:    info.Name(),
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
			Path:    path,
		}

		if err != nil {
			tools.Check(err)
		}

		files = append(files, data)

		return nil
	})

	if err != nil {
		tools.Check(err)
	}

	output, err := json.Marshal(files)
	tools.Check(err)

	return ctx.String(string(output))
}
