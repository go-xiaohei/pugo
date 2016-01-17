package migrate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo/app/helper"
	"github.com/go-xiaohei/pugo/app/migrate/def"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/yaml.v2"
)

const (
	// TypeHexo is type string of hexo migration
	TypeHexo = "Hexo"
)

var (
	_ Task = new(HexoTask)

	// ErrHexoSourceDirMissing means the hexo source directory is not found
	ErrHexoSourceDirMissing = fmt.Errorf("Hexo source directory is not found")
	// ErrHexoConfigMissing means the _config.yml is missing in source directory
	ErrHexoConfigMissing = fmt.Errorf("Hexo _config.yml is missing")
	// ErrHexoParsePostFail means get error when parse hexo's post
	ErrHexoParsePostFail = func(file string) error {
		return fmt.Errorf("%s need two block in hexo posts", file)
	}

	hexoBlockSeperate = "---"
	hexoReplacer      = strings.NewReplacer(" ", "-", ",", "-", ".", "-")
)

type (
	// HexoTask is migration of Hexo
	HexoTask struct {
		opt    *HexoOption
		result map[string]*bytes.Buffer
		err    error
	}
	// HexoOption is options for HexoTask
	HexoOption struct {
		Dest   string
		Source string
	}
	hexoSources struct {
		config *def.HexoConfig
		posts  []string
		pages  map[string]string
	}
)

// Is checks conf is supported to HexoTask
func (ht *HexoTask) Is(conf string) bool {
	return strings.HasPrefix(conf, "hexo://")
}

// New creates new HexoTask with cli context
func (ht *HexoTask) New(ctx *cli.Context) (Task, error) {
	u, err := url.Parse(ctx.String("src"))
	if err != nil {
		return nil, err
	}

	opt := &HexoOption{
		Dest: ctx.String("dest"),
	}
	opt.Source = strings.TrimSuffix(u.Path, "/")
	if opt.Source == "" {
		opt.Source = strings.TrimSuffix(u.Host, "/")
	}
	if !com.IsDir(opt.Source) {
		log15.Error("Hexo.Source.Dir.[" + opt.Source + "].Missing")
		return nil, ErrHexoSourceDirMissing
	}

	return &HexoTask{
		opt:    opt,
		result: make(map[string]*bytes.Buffer),
	}, nil
}

// Type returns HexoTask type string
func (ht *HexoTask) Type() string {
	return TypeHexo
}

// Do does HexoTask migration
func (ht *HexoTask) Do() (map[string]*bytes.Buffer, error) {
	var s *hexoSources
	if s, ht.err = ht.readSources(); ht.err != nil {
		return nil, ht.err
	}
	ht.err = ht.compileSources(s)
	return ht.result, ht.err
}

func (ht *HexoTask) readSources() (*hexoSources, error) {
	sources := &hexoSources{
		pages: make(map[string]string),
	}
	// test config.yml
	configFile := filepath.Join(ht.opt.Source, "_config.yml")
	if !com.IsFile(configFile) {
		return nil, ErrHexoConfigMissing
	}
	// parse config
	config, err := readConfig(configFile)
	if err != nil {
		return nil, err
	}
	sources.config = config

	// read source files
	srcDir := filepath.Join(ht.opt.Source, config.SourceDir)
	err = filepath.Walk(srcDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(p) != ".md" {
			return nil
		}
		p2, _ := filepath.Rel(srcDir, p)
		if strings.HasPrefix(p2, "_posts/") {
			sources.posts = append(sources.posts, p)
			return nil
		}
		base := filepath.Base(p2)
		if base == "index.md" {
			p2 = strings.TrimSuffix(p2, "/index.md")
		}
		sources.pages[p2] = p
		return nil
	})
	return sources, err
}

func (ht *HexoTask) compileSources(s *hexoSources) error {
	var (
		fileName      string
		fileBytes     []byte
		fileBytesPart [][]byte
		err           error
		timeLayout    = helper.TimeLayoutReplacer.Replace(fmt.Sprintf("%s %s", s.config.DateFormat, s.config.TimeFormat))
	)

	println(timeLayout)

	// parse posts
	for _, postFile := range s.posts {
		if fileBytes, err = ioutil.ReadFile(postFile); err != nil {
			return err
		}

		fileBytesPart = bytes.SplitN(fileBytes, []byte(hexoBlockSeperate), 2)
		if len(fileBytesPart) != 2 {
			return ErrHexoParsePostFail(postFile)
		}
		b, meta, err := fixHexoMeta(fileBytesPart[0], "", timeLayout)
		if err != nil {
			return err
		}
		b.WriteString("\n")
		b.Write(fileBytesPart[1])

		t, err := time.Parse(timeLayout, meta.Created)
		if err != nil {
			return err
		}
		fileName = fmt.Sprintf("post/%s/%s.md", t.Format("2006"), hexoReplacer.Replace(meta.Title))
		ht.result[fileName] = b
		log15.Debug("Hexo.Generate.[" + fileName + "]")

		fileName, fileBytes, fileBytesPart = "", nil, nil
	}

	// parse pages
	for slug, pageFile := range s.pages {
		if fileBytes, err = ioutil.ReadFile(pageFile); err != nil {
			return err
		}

		fileBytesPart = bytes.SplitN(fileBytes, []byte(hexoBlockSeperate), 2)
		if len(fileBytesPart) != 2 {
			return ErrHexoParsePostFail(pageFile)
		}
		b, meta, err := fixHexoMeta(fileBytesPart[0], slug, timeLayout)
		if err != nil {
			return err
		}
		b.WriteString("\n")
		b.Write(fileBytesPart[1])

		fileName = fmt.Sprintf("page/%s.md", hexoReplacer.Replace(meta.Title))
		ht.result[fileName] = b
		log15.Debug("Hexo.Generate.[" + fileName + "]")

		fileName, fileBytes, fileBytesPart = "", nil, nil
	}

	// write meta
	metaBuf := bytes.NewBufferString("[meta]\n")
	metaBuf.WriteString(fmt.Sprintf("title = %s\n", s.config.Title))
	metaBuf.WriteString(fmt.Sprintf("subtitle = %s\n", s.config.Subtitle))
	metaBuf.WriteString("keyword = \n")
	metaBuf.WriteString(fmt.Sprintf("desc = %s\n", s.config.Desc))

	u, err := url.Parse(s.config.URL)
	if err != nil {
		return err
	}
	metaBuf.WriteString(fmt.Sprintf("domain = %s\n", u.Host))
	metaBuf.WriteString(fmt.Sprintf("root = %s\n", s.config.URL))
	metaBuf.WriteString(fmt.Sprintf("lang = %s\n", s.config.Language))
	metaBuf.WriteString("cover = \n\n")

	metaBuf.WriteString(hexoDefaultNav)
	metaBuf.WriteString("\n\n")

	if s.config.Author != "" {
		metaBuf.WriteString("[author]\n-:owner\n\n")
		metaBuf.WriteString(fmt.Sprintf("[author.owner]\nname=%s\n\n", s.config.Author))
	}

	metaBuf.WriteString(migrateMetaExtraString)
	ht.result["meta.ini"] = metaBuf
	log15.Debug("Hexo.Generate.[meta.ini]")

	return nil
}

func readConfig(file string) (*def.HexoConfig, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	config := &def.HexoConfig{}
	if err = yaml.Unmarshal(bytes, config); err != nil {
		return nil, err
	}
	return config, nil
}

func fixHexoMeta(b []byte, slug, timeLayout string) (*bytes.Buffer, *def.HexoMeta, error) {
	var (
		meta = &def.HexoMeta{}
		err  error
	)
	if err = yaml.Unmarshal(b, meta); err != nil {
		return nil, nil, err
	}
	buf := bytes.NewBufferString("```ini\n")
	buf.WriteString(fmt.Sprintf("title = %s\n", meta.Title))
	if slug == "" {
		slug = hexoReplacer.Replace(meta.Title)
	}
	buf.WriteString(fmt.Sprintf("slug = %s\n", slug))
	if c, err := time.Parse(timeLayout, meta.Created); err != nil {
		buf.WriteString(fmt.Sprintf("date = %s\n", meta.Created))
	} else {
		buf.WriteString(fmt.Sprintf("date = %s\n", c.Format("2006-01-02 15:04:05")))
	}
	if meta.Updated != "" {
		if c, err := time.Parse(timeLayout, meta.Created); err != nil {
			buf.WriteString(fmt.Sprintf("update_date = %s\n", meta.Updated))
		} else {
			buf.WriteString(fmt.Sprintf("update_date = %s\n", c.Format("2006-01-02 15:04:05")))
		}
	}
	if str, ok := meta.Tags.(string); ok {
		buf.WriteString(fmt.Sprintf("tags = %s\n", str))
	}
	if s, ok := meta.Tags.([]interface{}); ok {
		var tags []string
		for _, str := range s {
			tags = append(tags, fmt.Sprint(str))
		}
		buf.WriteString(fmt.Sprintf("tags = %s\n", strings.Join(tags, ",")))
	}
	buf.WriteString("\n```\n")
	return buf, meta, nil
}

var (
	// default navigation for Hexo
	hexoDefaultNav = `[nav]
-:home
-:archive

[nav.home]
link = /
title = Home
i18n = home
hover = home

[nav.archive]
link = /archive
title = Archive
i18n = archive
hover = archive`
)
