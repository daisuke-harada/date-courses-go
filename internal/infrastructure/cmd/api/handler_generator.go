//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/api/gen"
	"github.com/samber/lo"
)

var re *regexp.Regexp = regexp.MustCompile("([a-z0-9])([A-Z])")

const handlerRootPath = "../handler"

// Templates are loaded from the repository's existing templates directory.
// This generator is expected to be invoked from `internal/infrastructure/cmd/api/gen`
// via `go generate`, so paths are relative to that directory.
const handlerTemplatePath = "../../../../../templates/handler.tmpl"
const constructorTemplatePath = "../../../../../templates/handler_constructor.tmpl"

func toSnakeCase(str string) string {
	snake := re.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func extractMethodNames() []reflect.Method {
	// apiパッケージからServerInterfaceインターフェースの型を取得
	// reflect.TypeOf((*api.ServerInterface)(nil))は、ServerInterface型のnilポインタを取得
	// nilポインタを取得するのは、インターフェースの型情報を取得するために具体的なインスタンスを作成する必要がないから
	// .Elem()は、そのポインタが指す要素の型を取得
	serverInterfaceType := reflect.TypeOf((*gen.ServerInterface)(nil)).Elem()

	var methods []reflect.Method
	// NumMethod()でインターフェイスに定義されているメソッドの数を取得する
	for i := 0; i < serverInterfaceType.NumMethod(); i++ {
		method := serverInterfaceType.Method(i)
		methods = append(methods, method)
	}

	return methods
}

func main() {
	methods := extractMethodNames()

	tmpl, err := template.ParseFiles(handlerTemplatePath)
	if err != nil {
		log.Fatalf("Failed to parse template: %s", err)
	}

	err = os.MkdirAll(handlerRootPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %s", err)
	}

	// handler.goないのコードを生成する
	createHandlerConstractor(methods)
	for _, method := range methods {
		fileName := toSnakeCase(method.Name) + ".go"
		filePath := filepath.Join(handlerRootPath, fileName)
		if fileNoExists(filePath) {
			createFile(filePath, tmpl, method)
		}
	}
}

func fileNoExists(filePath string) bool {
	// エラーは、指定されたパスが存在しない場合や、アクセス権限がない場合に返されるため、ファイルの存在可否を確認できる
	_, err := os.Stat(filePath)
	return os.IsNotExist(err)
}

func createHandlerConstractor(methods []reflect.Method) {
	handlerNames := lo.Map(methods, func(method reflect.Method, _ int) string {
		return method.Name
	})

	constractorTmpl, err := template.ParseFiles(constructorTemplatePath)
	if err != nil {
		log.Fatalf("Failed to parse template: %s", err)
	}

	// handlerのconstractorを作成
	constractorFile, err := os.Create(filepath.Join(handlerRootPath, "handler.go"))
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer constractorFile.Close()

	constractorTmpl.Execute(constractorFile, struct {
		HandlerNames []string
	}{
		HandlerNames: handlerNames,
	})
}

func createFile(filePath string, tmpl *template.Template, method reflect.Method) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	var args []struct {
		Name string
		Type string
	}

	for j := 0; j < method.Type.NumIn(); j++ {
		argType := method.Type.In(j)
		args = append(args, struct {
			Name string
			Type string
		}{
			Name: fmt.Sprintf("arg%d", j),
			Type: argType.String(),
		})
	}

	if err := tmpl.Execute(file, struct {
		MethodName string
		Args       []struct {
			Name string
			Type string
		}
	}{
		MethodName: method.Name,
		Args:       args,
	}); err != nil {
		log.Fatalf("Failed to execute template: %s", err)
	}

	log.Printf("Generated %s\n", filePath)
}
