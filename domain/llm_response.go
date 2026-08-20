package domain

type GeneratedFile struct {
	FilePath string `json:"filePath"`
	Response string `json:"response"`
}

type CodeResponse struct {
	Files []GeneratedFile `json:"files"`
}
