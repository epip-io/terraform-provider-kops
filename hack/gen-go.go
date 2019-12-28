package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/epip-io/terraform-provider-kops/pkg/util/templater"
	"github.com/iancoleman/strcase"
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/util/pkg/reflectutils"
)

const (
	JsonNameTag string = "json"

	// TemplateDir is the directory where the Go Template live
	TemplateDir = "./hack/provider"

	// OutputDir is the directory to put the rendered templates
	OutputDir = "./pkg"
)

func main() {
	rInterfaces := []interface{}{
		kops.Cluster{},
		kops.InstanceGroup{},
		kops.SSHCredential{},
	}

	resources := &Resources{
		Elems:     make(map[string]*Elem),
		Resources: make(map[string]*Resource),
		Functions: make(map[string]*Function),
		Schemas:   make(map[string]*Schema),
	}

	for _, in := range rInterfaces {
		if err := resources.AddResource(in); err != nil {
			log.Fatalf("failed to add resource: %#v", err)
		}

		if err := resources.AddSchema(in); err != nil {
			log.Fatalf("failed to add schema to resource: %#v", err)
		}

		if err := resources.AddFunction(in); err != nil {
			log.Fatalf("failed to add function to resource: %#v", err)
		}
	}

	for _, r := range resources.Resources {
		for sn, s := range r.Schemas {
			if s.Shared {
				delete(r.Schemas, sn)
			}
		}
		for fn, f := range r.Functions {
			if f.Shared {
				delete(r.Functions, fn)
			}
		}
	}

	for n, s := range resources.Schemas {
		if !s.Shared {
			delete(resources.Schemas, n)
		}
	}
	for fn, f := range resources.Functions {
		if !f.Shared {
			delete(resources.Functions, fn)
		}
	}

	tplAbsDir, _ := filepath.Abs(TemplateDir)
	outAbsDir, _ := filepath.Abs(OutputDir)

	fs, err := ioutil.ReadDir(tplAbsDir)
	if err != nil {
		log.Fatalf("unable to access dir: %e", err)
	}

	v := make(map[string]interface{})

	reflectutils.JsonMergeStruct(&v, resources)

	tr := templater.New(false, func(dir string, fs []os.FileInfo) []string {
		n := len(fs)
		r := make([]string, n)
		for i, f := range fs {
			r[i] = dir + "/" + f.Name()
		}
		return r
	}(tplAbsDir, fs))
	tr.AddRender("values.yaml", tplAbsDir+"/values.yaml", v)

	for n, r := range resources.Resources {
		v := make(map[string]interface{})
		reflectutils.JsonMergeStruct(&v, r)

		tr.AddRender(
			outAbsDir+"/provider/expand_"+strcase.ToSnake(n)+".gen.go",
			fmt.Sprintf("%s/expand.tpl.go", tplAbsDir),
			map[string]interface{}{
				n: v,
			},
		)

		tr.AddRender(
			outAbsDir+"/provider/flatten_"+strcase.ToSnake(n)+".gen.go",
			fmt.Sprintf("%s/flatten.tpl.go", tplAbsDir),
			map[string]interface{}{
				n: v,
			},
		)

		tr.AddRender(
			outAbsDir+"/provider/resource_"+strcase.ToSnake(n)+".gen.go",
			fmt.Sprintf("%s/resource.tpl.go", tplAbsDir),
			map[string]interface{}{
				n: v,
			},
		)
	}

	tr.AddRender(
		fmt.Sprintf("%s/provider/expand.gen.go", outAbsDir),
		fmt.Sprintf("%s/expand.tpl.go", tplAbsDir),
		v,
	)

	tr.AddRender(
		fmt.Sprintf("%s/provider/flatten.gen.go", outAbsDir),
		fmt.Sprintf("%s/flatten.tpl.go", tplAbsDir),
		v,
	)

	tr.AddRender(
		fmt.Sprintf("%s/provider/schema.gen.go", outAbsDir),
		fmt.Sprintf("%s/resource.tpl.go", tplAbsDir),
		v,
	)

	tr.AddRender(
		fmt.Sprintf("%s/provider/provider.gen.go", outAbsDir),
		fmt.Sprintf("%s/provider.tpl.go", tplAbsDir),
		v,
	)

	tr.Renderer()
}

type Elem struct {
	Name     string
	Tag      string
	Type     string
	Function *Function
	Schema   *Schema
}

type Function struct {
	Name      string
	Type      string
	Interface string
	Elems     map[string]*Elem
	Shared    bool
}

type Schema struct {
	Name     string
	Type     string
	SubType  string
	Required bool
	Elems    map[string]*Elem
	Shared   bool
}

type Resource struct {
	Name      string
	Elems     map[string]*Elem
	Functions map[string]*Function
	Schemas   map[string]*Schema

	reflectType reflect.Type
}

type Resources struct {
	Elems     map[string]*Elem
	Resources map[string]*Resource
	Functions map[string]*Function
	Schemas   map[string]*Schema
}

func (rs Resources) AddResource(in interface{}) error {
	t := reflect.TypeOf(in)

	rs.Resources[t.Name()] = &Resource{
		Name:        t.Name(),
		Elems:       make(map[string]*Elem),
		Functions:   make(map[string]*Function),
		Schemas:     make(map[string]*Schema),
		reflectType: t,
	}

	if err := reflectutils.ReflectRecursive(reflect.New(t), func(path string, field *reflect.StructField, v reflect.Value) error {
		pathSlice := strings.Split(path, ".")[1:]

		if len(pathSlice) > 0 && field != nil {
			e := &Elem{
				Name: field.Name,
				Tag:  field.Tag.Get(JsonNameTag),
				Type: field.Type.String(),
			}

			if err := rs.BuildSchema(e, field); err != nil {
				return err
			}

			if err := rs.BuildFunction(e, field); err != nil {
				return err
			}
			if e.Schema != nil {
				if len(pathSlice) > 1 {
					var ok bool
					for i := 0; i < len(pathSlice); i++ {
						if _, ok = rs.Resources[t.Name()].Elems[pathSlice[i]]; ok {
							break
						}
					}

					if !ok {
						rs.Resources[t.Name()].Elems[e.Name] = e
					}
				} else {
					rs.Resources[t.Name()].Elems[e.Name] = e
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (rs Resources) AddFunction(in interface{}) error {
	t := reflect.TypeOf(in)
	rNum := len(rs.Resources)

	for fn, f := range rs.Functions {
		if rNum == 1 {
			rs.Resources[t.Name()].Functions[fn] = f
			continue
		}

		var ok bool
		for n, r := range rs.Resources {
			if _, ok = r.Functions[fn]; ok && n != t.Name() {
				break
			}
		}

		if !ok {
			f.Shared = false
			rs.Resources[t.Name()].Functions[fn] = f
		}
	}

	return nil
}

func (rs Resources) AddSchema(in interface{}) error {
	t := reflect.TypeOf(in)
	rNum := len(rs.Resources)

	for sn, s := range rs.Schemas {
		if rNum == 1 {
			rs.Resources[t.Name()].Schemas[sn] = s
			continue
		}

		var ok bool
		for n, r := range rs.Resources {
			if _, ok = r.Schemas[sn]; ok && n != t.Name() {
				break
			}
		}

		if !ok {
			s.Shared = false
			rs.Resources[t.Name()].Schemas[sn] = s
		}
	}

	return nil
}

func (rs Resources) BuildFunction(e *Elem, f *reflect.StructField) error {
	var fun *Function
	var ok bool

	if f.Name == "" {
		return nil
	}

	fName := rs.BuildFunctionName(f.Type, f.Tag.Get(JsonNameTag))

	if fun, ok = rs.Functions[fName]; ok {
		if len(rs.Resources) > 1 {
			fun.Shared = true
		}
		e.Function = fun
	} else {
		e.Function = &Function{
			Name:  fName,
			Elems: rs.BuildFunctionElems(f.Type),
			Type:  f.Type.String(),
		}
		rs.Functions[fName] = e.Function
	}

	return nil
}

func (rs Resources) BuildFunctionElems(t reflect.Type) map[string]*Elem {
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		return rs.BuildFunctionElems(t.Elem())
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	r := make(map[string]*Elem)

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath != "" {
			continue
		}

		var e *Elem

		if _, ok := rs.Schemas[t.Name()]; ok {
			e, _ = rs.Schemas[t.Name()].Elems[t.Field(i).Name]
		}

		if e == nil {
			e = &Elem{
				Name: t.Field(i).Name,
				Tag:  t.Field(i).Tag.Get(JsonNameTag),
				Type: t.Field(i).Type.String(),
			}
		}

		f := t.Field(i)

		if err := rs.BuildFunction(e, &f); err != nil {
			log.Fatalf("Failed to Build Schema: %#v; %#v", e, err)
		}

		r[e.Name] = e
	}

	return r
}

func (rs Resources) BuildFunctionName(t reflect.Type, tag string) string {
	switch t.Kind() {
	case reflect.Ptr:
		return rs.BuildFunctionName(t.Elem(), tag)
	case reflect.Struct:
		return t.Name()
	case reflect.Slice:
		return rs.BuildFunctionName(t.Elem(), tag) + "Slice"
	case reflect.Map:
		return rs.BuildFunctionName(t.Elem(), tag) + "Map"
	default:
		return strcase.ToCamel(t.Name())
	}
}

func (rs Resources) BuildSchema(e *Elem, f *reflect.StructField) error {
	fTag := f.Tag.Get(JsonNameTag)
	fJName := strings.Split(fTag, ",")[0]

	if fJName == "" {
		return nil
	}

	var s *Schema
	var ok bool

	sName := rs.BuildSchemaName(f.Type, fTag)
	if s, ok = rs.Schemas[sName]; ok {
		if len(rs.Resources) > 1 {
			s.Shared = true
		}
		e.Schema = s
	} else {
		e.Schema = &Schema{
			Name:     sName,
			Elems:    rs.BuildSchemaElems(f.Type),
			Required: rs.IsSchemaRequired(fTag),
			Type:     rs.BuildSchemaType(f.Type),
		}
		rs.Schemas[sName] = e.Schema
	}

	return nil
}

func (rs Resources) BuildSchemaElems(t reflect.Type) map[string]*Elem {
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		return rs.BuildSchemaElems(t.Elem())
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	r := make(map[string]*Elem)

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath != "" {
			continue
		}

		e := &Elem{
			Name: t.Field(i).Name,
			Tag:  t.Field(i).Tag.Get(JsonNameTag),
			Type: t.Field(i).Type.String(),
		}

		f := t.Field(i)
		if err := rs.BuildSchema(e, &f); err != nil {
			log.Fatalf("Failed to Build Schema: %#v; %#v", e, err)
		}

		r[e.Name] = e
	}

	return r
}

func (rs Resources) BuildSchemaName(t reflect.Type, tag string) string {
	switch t.Kind() {
	case reflect.Ptr:
		return rs.BuildSchemaName(t.Elem(), tag)
	case reflect.Struct:
		return t.Name()
	case reflect.Slice:
		return rs.BuildSchemaName(t.Elem(), tag) + "Slice"
	case reflect.Map:
		return rs.BuildSchemaName(t.Elem(), tag) + "Map"
	case reflect.Float32, reflect.Float64, reflect.Int64:
		if rs.IsSchemaRequired(tag) {
			return "RequiredFloat"
		}

		return "OptionalFloat"
	case reflect.Int, reflect.Int16, reflect.Int32:
		if rs.IsSchemaRequired(tag) {
			return "RequiredInt"
		}

		return "OptionalInt"
	default:
		if rs.IsSchemaRequired(tag) {
			return "Required" + strcase.ToCamel(t.Name())
		}

		return "Optional" + strcase.ToCamel(t.Name())
	}
}

func (rs Resources) BuildSchemaType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return rs.BuildSchemaType(t.Elem())
	case reflect.Slice, reflect.Struct:
		return "List"
	case reflect.Map:
		return "Map"
	case reflect.Int, reflect.Int16, reflect.Int32:
		return "Int"
	case reflect.Int64, reflect.Float32, reflect.Float64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "Float"
	default:
		return strcase.ToCamel(t.Kind().String())
	}
}

func (rs Resources) BuildTypeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return "*" + rs.BuildTypeName(t.Elem())
	case reflect.Slice:
		return "[]" + rs.BuildTypeName(t.Elem())
	case reflect.Map:
		return "map[" + rs.BuildTypeName(t.Key()) + "]" + rs.BuildTypeName(t.Elem())
	default:
		return t.Name()
	}
}

func (rs Resources) IsSchemaRequired(tag string) bool {
	if strings.Contains(tag, "omitempty") {
		return false
	}

	return true
}
