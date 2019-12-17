package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
	"k8s.io/kops/pkg/apis/kops"

	templater "github.com/epip-io/terraform-provider-kops/pkg/util/templater"
)

const (
	// TemplateDir is the directory where the Go Template live
	TemplateDir = "./hack/api"

	// OutputDir is the directory to put the rendered templates
	OutputDir = "./pkg"
)

func main() {
	schema := generateTemplateValues([]*structs.Struct{
		structs.New(kops.Cluster{}),
		structs.New(kops.InstanceGroup{}),
		// structs.New(v1alpha2.SSHCredential{}),
	})

	tplAbsDir, _ := filepath.Abs(TemplateDir)
	outAbsDir, _ := filepath.Abs(OutputDir)

	fs, err := ioutil.ReadDir(tplAbsDir)
	if err != nil {
		log.Fatalf("unable to access dir: %e", err)
	}

	tr := templater.New(false, func(dir string, fs []os.FileInfo) []string {
		n := len(fs)
		r := make([]string, n)
		for i, f := range fs {
			r[i] = fmt.Sprintf("%s/%s", dir, f.Name())
		}
		return r
	}(tplAbsDir, fs))

	// tr.AddRender("values.yaml", fmt.Sprintf("%s/values.yaml", tplAbsDir), schema)

	os.MkdirAll(fmt.Sprintf("%s", outAbsDir), os.ModePerm)
	for _, r := range schema["Resources"].([]string) {
		types := schema["Resource"].(map[string]interface{})
		types = types[r].(map[string]interface{})

		for _, t := range types {
			ts := t.(map[string]interface{})

			tnp := strings.Split(ts["Type"].(string), ".")
			tn := tnp[len(tnp)-1]

			rv := map[string]interface{}{
				"Type":   tn,
				"Schema": ts,
			}

			tr.AddRender(fmt.Sprintf("%s/convert/expand_%s.gen.go", outAbsDir, strcase.ToSnake(tn)), fmt.Sprintf("%s/expand.tpl.go", tplAbsDir), rv)
			tr.AddRender(fmt.Sprintf("%s/convert/flatten_%s.gen.go", outAbsDir, strcase.ToSnake(tn)), fmt.Sprintf("%s/flatten.tpl.go", tplAbsDir), rv)
			tr.AddRender(fmt.Sprintf("%s/api/schema_%s.gen.go", outAbsDir, strcase.ToSnake(tn)), fmt.Sprintf("%s/schema.tpl.go", tplAbsDir), rv)
		}

		tr.AddRender(fmt.Sprintf("%s/provider/datasource.gen.go", outAbsDir), fmt.Sprintf("%s/datasource.tpl.go", tplAbsDir), schema)
		tr.AddRender(fmt.Sprintf("%s/provider/resources.gen.go", outAbsDir), fmt.Sprintf("%s/resource.tpl.go", tplAbsDir), schema)
		tr.AddRender(fmt.Sprintf("%s/provider/provider.gen.go", outAbsDir), fmt.Sprintf("%s/provider.tpl.go", tplAbsDir), schema)
	}

	tr.Renderer()
}

func generateTemplateValues(r []*structs.Struct) map[string]interface{} {
	schema := make(map[string]interface{})
	seen := make([]string, 0)
	types := make([]string, 0)
	schema = make(map[string]interface{})
	for _, s := range r {
		if _, ok := schema["Resources"]; !ok {
			schema["Resources"] = make([]string, 0)
		}
		if _, ok := schema["Resource"]; !ok {
			schema["Resource"] = make(map[string]interface{}, 0)
		}
		if _, ok := schema["Resource"].(map[string]interface{})[s.Name()]; !ok {
			schema["Resources"] = append(schema["Resources"].([]string), s.Name())
			schema["Resource"].(map[string]interface{})[s.Name()] = generateFields(s, &seen, &types)
		}
		schema["Types"] = types
	}

	return schema
}

func generateFields(st *structs.Struct, seen *[]string, types *[]string) map[string]interface{} {
	rf := make(map[string]interface{})

	(*seen) = append((*seen), st.Name())

	for _, f := range st.Fields() {
		if !f.IsExported() {
			continue
		}

		if f.Tag("json") == "" {
			continue
		}

		if f.Tag("json") == ",inline" {
			continue
		}

		i := f.Value()

		if types != nil {
			done := false
			for _, tn := range *types {
				if reflect.TypeOf(i).String() == tn {
					done = true
				}
			}

			if done {
				continue
			}
		}

		rf[f.Name()] = map[string]interface{}{
			"Name":     strcase.ToSnake(strings.Split(f.Tag("json"), ",")[0]),
			"Kind":     f.Kind(),
			"Type":     reflect.TypeOf(i).String(),
			"Required": true,
			"IsPtr": func() bool {
				if f.Kind() == reflect.Ptr {
					return true
				}
				return false
			}(),
			"IsSlice": func() bool {
				if f.Kind() == reflect.Slice {
					return true
				}

				return false
			}(),
			"First":  nil,
			"Fields": nil,
			"Top":    false,
			"Seen":   false,
		}

		switch f.Kind() {
		case reflect.Ptr, reflect.Slice:
			v := reflect.New(reflect.TypeOf(i).Elem()).Elem()
			k := v.Kind()

			i = v.Interface()

			switch rf[f.Name()].(map[string]interface{})["Kind"] {
			case reflect.Ptr:
				rf[f.Name()].(map[string]interface{})["First"] = "ptr"
			case reflect.Slice:
				rf[f.Name()].(map[string]interface{})["First"] = "slice"
			}

			rf[f.Name()].(map[string]interface{})["Type"] = reflect.TypeOf(i).String()
			rf[f.Name()].(map[string]interface{})["Kind"] = k
		}

		switch rf[f.Name()].(map[string]interface{})["Kind"] {
		case reflect.Ptr, reflect.Slice:
			v := reflect.New(reflect.TypeOf(i).Elem()).Elem()
			k := v.Kind()

			i = v.Interface()

			switch rf[f.Name()].(map[string]interface{})["Kind"] {
			case reflect.Ptr:
				rf[f.Name()].(map[string]interface{})["IsPtr"] = true
			case reflect.Slice:
				rf[f.Name()].(map[string]interface{})["IsSlice"] = true
			}

			rf[f.Name()].(map[string]interface{})["Type"] = reflect.TypeOf(i).String()
			rf[f.Name()].(map[string]interface{})["Kind"] = k
		}

		switch rf[f.Name()].(map[string]interface{})["Kind"] {
		case reflect.Struct:
			s := structs.New(i)

			done := false
			for _, stn := range *seen {
				if s.Name() == stn {
					done = true
				}
			}

			if done {
				rf[f.Name()].(map[string]interface{})["Seen"] = true
			} else {
				rf[f.Name()].(map[string]interface{})["Fields"] = generateFields(s, seen, nil)
			}
		}

		rf[f.Name()].(map[string]interface{})["Kind"] = rf[f.Name()].(map[string]interface{})["Kind"].(reflect.Kind).String()

		if strings.Contains(rf[f.Name()].(map[string]interface{})["Name"].(string), "cid_rs") {
			rf[f.Name()].(map[string]interface{})["Name"] = strings.Replace(rf[f.Name()].(map[string]interface{})["Name"].(string), "cid_rs", "cidrs", -1)
		}

		if strings.Contains(rf[f.Name()].(map[string]interface{})["Name"].(string), "k_8_s") {
			rf[f.Name()].(map[string]interface{})["Name"] = strings.Replace(rf[f.Name()].(map[string]interface{})["Name"].(string), "k_8_s", "k8s", -1)
		}

		if strings.Contains(f.Tag("json"), "omitempty") {
			rf[f.Name()].(map[string]interface{})["Required"] = false
		}

		if types != nil {
			rf[f.Name()].(map[string]interface{})["Top"] = true
			(*types) = append((*types), rf[f.Name()].(map[string]interface{})["Type"].(string))
		}
	}

	return rf
}
