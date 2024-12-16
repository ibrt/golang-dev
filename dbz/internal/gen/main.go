//go:generate go run .
package main

import (
	"fmt"
	"path/filepath"

	"github.com/ibrt/golang-utils/filez"
	"github.com/ibrt/golang-utils/gzipz"

	"github.com/ibrt/golang-dev/consolez"
	"github.com/ibrt/golang-dev/dbz/internal/assets"
	"github.com/ibrt/golang-dev/gtz"
	"github.com/ibrt/golang-dev/shellz"
)

func main() {
	consolez.DefaultCLI.Notice("dbz-gen", "setting up...")

	pluginFilePath := filez.MustAbs(filepath.Join("..", "assets", "sqlc", "plugin.wasm.gz"))
	dirPath := filez.MustCreateTempDir()

	defer filez.MustRemoveAll(dirPath)
	filez.MustChdir(dirPath)

	consolez.DefaultCLI.Notice("dbz-gen", "cloning sources...")
	shellz.NewCommand("git", "-c", "advice.detachedHead=false", "clone",
		"--quiet", "--depth", "1", "--branch", gtz.GoToolSQLCGenGo.GetVersion(),
		fmt.Sprintf("https://%v", gtz.GoToolSQLCGenGo.GetPackage())).
		SetDir(dirPath).
		MustRun()

	consolez.DefaultCLI.Notice("dbz-gen", "merging templates...")
	filez.MustExport(assets.SQLCTemplatesFS, assets.SQLCTemplatesFSRoot,
		filepath.Join(dirPath, "sqlc-gen-go", "internal", "templates"))

	consolez.DefaultCLI.Notice("dbz-gen", "building plugin...")
	shellz.NewCommand("go", "build", "-v", "-o", "plugin.wasm", "main.go").
		SetEnv("GOOS", "wasip1").
		SetEnv("GOARCH", "wasm").
		SetDir(filepath.Join(dirPath, "sqlc-gen-go", "plugin")).
		MustRun()

	consolez.DefaultCLI.Notice("dbz-gen", "storing artifacts...")
	buf := gzipz.MustCompress(filez.MustReadFile(filepath.Join(dirPath, "sqlc-gen-go", "plugin", "plugin.wasm")))
	filez.MustWriteFile(pluginFilePath, 0777, 0666, buf)
	consolez.DefaultCLI.Notice("dbz-gen", "generated", pluginFilePath)

}
