package builder

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"pugo/pkg/models"

	"github.com/pelletier/go-toml"
)

var (
	// ErrMissingMetadataFile is error of missing metadata file
	ErrMissingMetadataFile = errors.New("missing metadata file")
	// ErrUnsupportedMetaExtension is error of parsing unsupported metadata extension
	ErrUnsupportedMetaExtension = errors.New("unsupported meta extension file")
)

// ParseMeta parses metadata file in content directory
func ParseMeta(contentDir string) (*models.MetaFile, error) {
	files := models.GetPossibleExtensionFile("meta")
	for extType, file := range files {
		metaFile := filepath.Join(contentDir, file)
		if _, err := os.Stat(metaFile); err != nil {
			continue
		}
		parser, err := newMetaParser(extType)
		if err != nil {
			return nil, err
		}
		return parser.ParseFile(metaFile)
	}
	return nil, ErrMissingMetadataFile
}

type metaParser interface {
	ParseFile(fpath string) (*models.MetaFile, error)
}

func newMetaParser(extType models.ExtensionType) (metaParser, error) {
	switch extType {
	case models.ExtensionToml:
		return &tomlMetaParser{}, nil
	}
	return nil, ErrUnsupportedMetaExtension
}

var (
	_ metaParser = (*tomlMetaParser)(nil)
)

type tomlMetaParser struct{}

func (tp *tomlMetaParser) ParseFile(fpath string) (*models.MetaFile, error) {
	fileBytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	var metaFileData models.MetaFile
	return &metaFileData, toml.Unmarshal(fileBytes, &metaFileData)
}
