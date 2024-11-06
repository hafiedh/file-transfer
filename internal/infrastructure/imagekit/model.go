package imagekit

type (
	UploadRequest struct {
		FileName          string `json:"fileName"`
		UseUniqueFileName string `json:"useUniqueFileName"`
		Folder            string `json:"folder"`
		Iat               int64  `json:"iat"`
		Exp               int64  `json:"exp"`
	}

	UploadResponse struct {
		VersionInfo  VersionInfo `json:"versionInfo"`
		FileID       string      `json:"fileId"`
		Name         string      `json:"name"`
		FilePath     string      `json:"filePath"`
		URL          string      `json:"url"`
		FileType     string      `json:"fileType"`
		ThumbnailURL string      `json:"thumbnailUrl"`
		Size         int64       `json:"size"`
		Height       int64       `json:"height"`
		Width        int64       `json:"width"`
		Orientation  int64       `json:"orientation"`
	}

	VersionInfo struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	BadResponse struct {
		Message string `json:"message"`
		Help    string `json:"help"`
	}
)
