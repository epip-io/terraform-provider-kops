package templater

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

// Render respresents a file, it's template and values
type Render struct {
	BasePath string
	Template string
	Values   map[string]interface{}
}

// File respresents a template file
type File struct {
	Dir      string
	BasePath string
	Content  string
}

// Engine is a simple template renderer
type Engine struct {
	Strict bool

	Files map[string]File

	Render map[string]Render
}

// New returns a new Engine
func New(strict bool, files []string) *Engine {
	return &Engine{
		Strict: strict,
		Files:  loadFiles(files),
	}
}

// AddRender adds a Render
func (e *Engine) AddRender(file, tplFile string, values map[string]interface{}) error {
	if v, ok := e.Files[tplFile]; ok {
		if e.Render == nil {
			e.Render = make(map[string]Render)
		}
		absFile, err := filepath.Abs(file)
		if err != nil {
			log.Fatalf("unable to find absolute path for file: %s: %v", file, err)
		}

		e.Render[absFile] = Render{
			BasePath: filepath.Dir(absFile),
			Template: v.Content,
			Values:   values,
		}
	}

	return fmt.Errorf("template doesn't exist: %v", tplFile)
}

// Renderer attempts render Go templates
func (e *Engine) Renderer() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("rendering template failed: %v", r)
		}
	}()

	t := template.New("gotpl")

	e.initFuncMap(t)

	for n, f := range e.Files {
		if strings.HasPrefix(filepath.Base(n), "_") {
			if _, err := t.New(n).Parse(f.Content); err != nil {
				return cleanupParseError(n, err)
			}
		}
	}

	for n, f := range e.Render {
		if _, err := os.Stat(f.BasePath); os.IsNotExist(err) {
			os.MkdirAll(f.BasePath, os.ModePerm)
		}

		if _, err := t.New(n).Parse(f.Template); err != nil {
			return cleanupParseError(n, err)
		}

		fh, _ := os.OpenFile(n, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err := t.ExecuteTemplate(fh, n, f.Values); err != nil {
			log.Fatalf("unable to write file: %s", n)
		}
	}

	return nil
}

// initFuncMap setups the function map to extend text/template functionality
func (e *Engine) initFuncMap(t *template.Template) {
	funcMap := FuncMap()

	t.Funcs(funcMap)
}

// loadFiles loads the Files into a map of Render
func loadFiles(fns []string) map[string]File {
	fs := make(map[string]File)

	for _, f := range fns {
		absFile, err := filepath.Abs(f)
		if err != nil {
			log.Fatalf("unable to find absolute path for file: %s: %v", absFile, err)
		}

		basePath := filepath.Dir(absFile)

		fc, err := ioutil.ReadFile(absFile)
		if err != nil {
			log.Fatalf("unable to load file (%s): %v", absFile, err)
		}

		fs[absFile] = File{
			Dir:      filepath.Base(basePath),
			BasePath: basePath,
			Content: fmt.Sprintf(
				"// Code generated by engine.go; DO NOT EDIT.\n\n%s",
				string(fc),
			),
		}
	}

	return fs
}

func cleanupParseError(filename string, err error) error {
	tokens := strings.Split(err.Error(), ": ")
	if len(tokens) == 1 {
		// This might happen if a non-templating error occurs
		return fmt.Errorf("parse error in (%s): %s", filename, err)
	}
	// The first token is "template"
	// The second token is either "filename:lineno" or "filename:lineNo:columnNo"
	location := tokens[1]
	// The remaining tokens make up a stacktrace-like chain, ending with the relevant error
	errMsg := tokens[len(tokens)-1]
	return fmt.Errorf("parse error at (%s): %s", string(location), errMsg)
}

func cleanupExecError(filename string, err error) error {
	if _, isExecError := err.(template.ExecError); !isExecError {
		return err
	}

	tokens := strings.SplitN(err.Error(), ": ", 3)
	if len(tokens) != 3 {
		// This might happen if a non-templating error occurs
		return fmt.Errorf("execution error in (%s): %s", filename, err)
	}

	// The first token is "template"
	// The second token is either "filename:lineno" or "filename:lineNo:columnNo"
	location := tokens[1]

	parts := tokens[2]
	if len(parts) >= 2 {
		return fmt.Errorf("execution error at (%s): %s", string(location), parts)
	}

	return err
}
