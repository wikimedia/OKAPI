package schema

// default license for the pages
const (
	LicenseIdentifier = "CC-BY-SA-3.0"
	LicenseName       = "Creative Commons Attribution Share Alike 3.0 Unported"
	LicenseURL        = "https://creativecommons.org/licenses/by-sa/3.0/"
)

// NewLicense create new license
func NewLicense() *License {
	return &License{
		Name:       LicenseName,
		Identifier: LicenseIdentifier,
		URL:        LicenseURL,
	}
}

// License schema
type License struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	URL        string `json:"url,omitempty"`
}
