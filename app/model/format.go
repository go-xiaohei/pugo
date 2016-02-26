package model

const (
	FormatTOML FormatType = 1
	FormatINI  FormatType = 2
	// FormatYAML FormatType = 3
)

type FormatType int8

func ShouldMetaFiles() map[FormatType]string {
	return map[FormatType]string{
		FormatTOML: "meta.toml",
		FormatINI:  "meta.ini",
	}
}
