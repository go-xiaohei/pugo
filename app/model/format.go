package model

import "strings"

const (
	// FormatTOML mean toml format
	FormatTOML FormatType = 1
	// FormatINI mean ini format
	FormatINI FormatType = 2
	// FormatYAML FormatType = 3
)

// FormatType define type of front-matter format, meta file and language file
type FormatType int8

// ShouldMetaFiles return all filenames of meta file in all format
func ShouldMetaFiles() map[FormatType]string {
	return map[FormatType]string{
		FormatTOML: "meta.toml",
		FormatINI:  "meta.ini",
	}
}

// ShouldPostMetaFiles return all filenames of meta file for posts
func ShouldPostMetaFiles() map[FormatType]string {
	return map[FormatType]string{
		FormatTOML: "post.toml",
		FormatINI:  "post.ini",
	}
}

// ShouldPageMetaFiles return all filenames of meta file for pages
func ShouldPageMetaFiles() map[FormatType]string {
	return map[FormatType]string{
		FormatTOML: "page.toml",
		FormatINI:  "page.ini",
	}
}

// ShouldThemeMetaFiles return all filenames of meta file for theme
func ShouldThemeMetaFiles() map[FormatType]string {
	return map[FormatType]string{
		FormatTOML: "theme.toml",
		FormatINI:  "theme.ini",
	}
}

// ShouldFormatExtension return all extensions of all formats
func ShouldFormatExtension() map[FormatType]string {
	return map[FormatType]string{
		FormatTOML: ".toml",
		FormatINI:  ".ini",
	}
}

func detectFormat(str string) FormatType {
	str = strings.TrimSpace(str)
	if str == "toml" || str == ".toml" {
		return FormatTOML
	}
	if str == "ini" || str == ".ini" {
		return FormatINI
	}
	return 0
}
