package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	reDefineTag   = regexp.MustCompile("{{ ?define \"([^\"]*)\" ?\"?([a-zA-Z0-9]*)?\"? ?}}")
	reTemplateTag = regexp.MustCompile("{{ ?template \"([^\"]*)\" ?([^ ]*)? ?}}")
)

type (
	// Theme object, maintains a sort of templates for whole site data
	Theme struct {
		dir        string
		lock       sync.Mutex
		funcMap    template.FuncMap
		templates  map[string]*template.Template
		extensions []string

		cache               []*namedTemplate
		regularTemplateDefs []string
	}
	namedTemplate struct {
		Name string
		Src  string
	}
)

// NewTheme returns new theme with dir, functions and extensions
func NewTheme(dir string, funcMap template.FuncMap, extension []string) *Theme {
	theme := &Theme{
		dir:        dir,
		funcMap:    funcMap,
		extensions: extension,
	}
	theme.funcMap["Include"] = func(tpl string, data interface{}) template.HTML {
		var buf bytes.Buffer
		if err := theme.Execute(&buf, tpl, data); err != nil {
			return template.HTML("<!-- template " + tpl + " error:" + err.Error() + "-->")
		}
		return template.HTML(string(buf.Bytes()))
	}
	return theme
}

// Load loads templates
func (th *Theme) Load() error {
	return th.loadTemplates()
}

// changes from https://github.com/go-macaron/renders/blob/master/renders.go#L43,
// thanks a lot
func (th *Theme) loadTemplates() error {
	th.lock.Lock()
	defer th.lock.Unlock()

	templates := make(map[string]*template.Template)

	err := filepath.Walk(th.dir, func(p string, fi os.FileInfo, err error) error {
		r, err := filepath.Rel(th.dir, p) // get relative path
		if err != nil {
			return err
		}
		ext := getExt(r)
		for _, extension := range th.extensions {
			if ext == extension {
				if err := th.add(p); err != nil {
					return err
				}
				for _, t := range th.regularTemplateDefs {
					found := false
					defineIdx := 0
					// From the beginning (which should) most specifc we look for definitions
					for _, nt := range th.cache {
						nt.Src = reDefineTag.ReplaceAllStringFunc(nt.Src, func(raw string) string {
							parsed := reDefineTag.FindStringSubmatch(raw)
							name := parsed[1]
							if name != t {
								return raw
							}
							// Don't touch the first definition
							if !found {
								found = true
								return raw
							}
							defineIdx++

							return fmt.Sprintf("{{ define \"%s_invalidated_#%d\" }}", name, defineIdx)
						})
					}
				}

				var (
					baseTmpl *template.Template
					i        int
				)

				for _, nt := range th.cache {
					var currentTmpl *template.Template
					if i == 0 {
						baseTmpl = template.New(nt.Name)
						currentTmpl = baseTmpl
					} else {
						currentTmpl = baseTmpl.New(nt.Name)
					}

					template.Must(currentTmpl.Funcs(th.funcMap).Parse(nt.Src))
					i++
				}
				tname := generateTemplateName(th.dir, p)
				templates[tname] = baseTmpl

				// Make sure we empty the cache between runs
				th.cache = th.cache[0:0]

				break
				//return nil
			}
		}
		return nil
	})
	th.templates = templates
	return err
}

func (th *Theme) add(path string) error {
	// Get file content
	tplSrc, err := getFileContent(path)
	if err != nil {
		return err
	}
	tplName := generateTemplateName(th.dir, path)
	// Make sure template is not already included
	alreadyIncluded := false
	for _, nt := range th.cache {
		if nt.Name == tplName {
			alreadyIncluded = true
			break
		}
	}
	if alreadyIncluded {
		return nil
	}

	// Add to the cache
	nt := &namedTemplate{
		Name: tplName,
		Src:  tplSrc,
	}
	th.cache = append(th.cache, nt)

	// Check for any template block
	for _, raw := range reTemplateTag.FindAllString(nt.Src, -1) {
		parsed := reTemplateTag.FindStringSubmatch(raw)
		templatePath := parsed[1]
		ext := getExt(templatePath)
		if !strings.Contains(templatePath, ext) {
			th.regularTemplateDefs = append(th.regularTemplateDefs, templatePath)
			continue
		}

		// Add this template and continue looking for more template blocks
		th.add(filepath.Join(th.dir, templatePath))
	}
	return nil
}

// Execute executes template by name with data,
// write into a Writer
func (th *Theme) Execute(w io.Writer, name string, data interface{}) error {
	tpl := th.Template(name)
	if tpl == nil {
		return errorFileMissing.New(name)
	}
	return tpl.ExecuteTemplate(w, name, data)
}

// Static gets static dir in the theme
func (th *Theme) Static() string {
	return path.Join(th.dir, "static")
}

// Template gets template by name
func (th *Theme) Template(name string) *template.Template {
	return th.templates[name]
}

func generateTemplateName(base, path string) string {
	//name := (r[0 : len(r)-len(ext)])
	return filepath.ToSlash(path[len(base)+1:])
}

func getFileContent(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	s := string(b)
	if len(s) < 1 {
		return "", errors.New("render: template file is empty")
	}
	return s, nil
}

func getExt(s string) string {
	if strings.Index(s, ".") == -1 {
		return ""
	}
	return "." + strings.Join(strings.Split(s, ".")[1:], ".")
}
