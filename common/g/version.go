package g

import "fmt"

var (
	BinaryName string
	Version    string
)

func VersionInfo() string {
	return fmt.Sprintf("%s", Version)
}

func HbsInfo() string {
	return fmt.Sprintf("%s.%s", BinaryName, Version)
}
