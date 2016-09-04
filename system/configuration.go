package system

// Configuration stores configuration data
type Configuration struct {
	PublicPath   string `json:"public_path"`
	TemplatePath string `json:"template_path"`
	UploadPath   string `json:"upload_path"`
	Bind         string `json:"bind"`
}
