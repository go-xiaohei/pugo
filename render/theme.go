package render

import (
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
	re_defineTag   *regexp.Regexp = regexp.MustCompile("{{ ?define \"([^\"]*)\" ?\"?([a-zA-Z0-9]*)?\"? ?}}")
	re_templateTag *regexp.Regexp = regexp.MustCompile("{{ ?template \"([^\"]*)\" ?([^ ]*)? ?}}")

	ErrTemplateMissing = errors.New("template-missing")
)

type (
	// theme object, maintains a sort of templates for whole site data
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

// new theme with dir, functions and extenions
func NewTheme(dir string, funcMap template.FuncMap, extension []string) *Theme {
	return &Theme{
		dir:        dir,
		funcMap:    funcMap,
		extensions: extension,
	}
}

// load templates
func (th *Theme) Load() error {
	return th.loadTemplates()
}

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
						nt.Src = re_defineTag.ReplaceAllStringFunc(nt.Src, func(raw string) string {
							parsed := re_defineTag.FindStringSubmatch(raw)
							name := parsed[1]
							if name != t {
								return raw
							}
							// Don't touch the first definition
							if !found {
								found = true
								return raw
							}

							defineIdx += 1

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

func (t *Theme) add(path string) error {
	// Get file content
	tplSrc, err := getFileContent(path)
	if err != nil {
		return err
	}
	tplName := generateTemplateName(t.dir, path)
	// Make sure template is not already included
	alreadyIncluded := false
	for _, nt := range t.cache {
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
	t.cache = append(t.cache, nt)

	// Check for any template block
	for _, raw := range re_templateTag.FindAllString(nt.Src, -1) {
		parsed := re_templateTag.FindStringSubmatch(raw)
		templatePath := parsed[1]
		ext := getExt(templatePath)
		if !strings.Contains(templatePath, ext) {
			t.regularTemplateDefs = append(t.regularTemplateDefs, templatePath)
			continue
		}

		// Add this template and continue looking for more template blocks
		t.add(filepath.Join(t.dir, templatePath))
	}
	return nil
}

// execute template by name with data,
// write into a Writer
func (t *Theme) Execute(w io.Writer, name string, data interface{}) error {
	tpl := t.Template(name)
	if tpl == nil {
		return ErrTemplateMissing
	}
	return tpl.ExecuteTemplate(w, name, data)
}

// static dir in the theme
func (t *Theme) Static() string {
	return path.Join(t.dir, "static")
}

// get template by name
func (t *Theme) Template(name string) *template.Template {
	return t.templates[name]
}

func generateTemplateName(base, path string) string {
	//name := (r[0 : len(r)-len(ext)])
	return filepath.ToSlash(path[len(base)+1:])
}

func getFileContent(path string) (string, error) {
	// Read the file content of the template
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
