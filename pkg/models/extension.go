package models

// ExtensionType means supported file extensions to parse
type ExtensionType string

const (
	// ExtensionToml is type of toml
	ExtensionToml ExtensionType = "toml"
	// ExtensionIni is type of ini
	ExtensionIni ExtensionType = "ini"
	// ExtensionYaml is type of yaml
	ExtensionYaml ExtensionType = "yaml"
)

// GetPossibleExtensionFile gets possible extension files
func GetPossibleExtensionFile(filename string) map[ExtensionType]string {
	return map[ExtensionType]string{
		ExtensionToml: filename + "." + string(ExtensionToml),
		ExtensionIni:  filename + "." + string(ExtensionIni),
		ExtensionYaml: filename + "." + string(ExtensionYaml),
	}
}
