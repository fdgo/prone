package domain

type UploadResult struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
}
