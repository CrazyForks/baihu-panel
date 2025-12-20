package static

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var distFS embed.FS

// GetFileSystem 返回嵌入的静态文件系统
func GetFileSystem() http.FileSystem {
	subFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return http.FS(subFS)
}

// GetFS 返回嵌入的 fs.FS
func GetFS() fs.FS {
	subFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return subFS
}

// ReadFile 读取嵌入的文件
func ReadFile(name string) ([]byte, error) {
	return distFS.ReadFile("dist/" + name)
}
