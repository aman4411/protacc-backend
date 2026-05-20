package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	driveFilePathRegex   = regexp.MustCompile(`(?i)/file/d/([a-zA-Z0-9_-]+)`)
	driveOpenRegex       = regexp.MustCompile(`(?i)[?&]id=([a-zA-Z0-9_-]+)`)
	docsPathRegex        = regexp.MustCompile(`(?i)/document/d/([a-zA-Z0-9_-]+)`)
	sheetsPathRegex      = regexp.MustCompile(`(?i)/spreadsheets/d/([a-zA-Z0-9_-]+)`)
	driveFolderPathRegex = regexp.MustCompile(`(?i)/folders/([a-zA-Z0-9_-]+)`)
)

type GoogleDriveLink struct {
	FileID      string `json:"fileId"`
	IsFolder    bool   `json:"isFolder"`
	EmbedURL    string `json:"embedUrl"`
	DownloadURL string `json:"downloadUrl"`
	ViewURL     string `json:"viewUrl"`
}

func ParseGoogleDriveURL(rawURL string) (*GoogleDriveLink, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil, fmt.Errorf("Google Drive URL is required")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format")
	}

	host := strings.ToLower(parsed.Host)
	if !strings.Contains(host, "google.com") && !strings.Contains(host, "googleusercontent.com") {
		return nil, fmt.Errorf("URL must be a Google Drive or Google Docs link")
	}

	fullURL := rawURL
	if !strings.HasPrefix(strings.ToLower(fullURL), "http") {
		fullURL = "https://" + fullURL
	}

	var fileID string
	isFolder := false

	if match := driveFolderPathRegex.FindStringSubmatch(fullURL); len(match) > 1 {
		fileID = match[1]
		isFolder = true
	} else if match := driveFilePathRegex.FindStringSubmatch(fullURL); len(match) > 1 {
		fileID = match[1]
	} else if match := docsPathRegex.FindStringSubmatch(fullURL); len(match) > 1 {
		fileID = match[1]
	} else if match := sheetsPathRegex.FindStringSubmatch(fullURL); len(match) > 1 {
		fileID = match[1]
	} else if match := driveOpenRegex.FindStringSubmatch(fullURL); len(match) > 1 {
		fileID = match[1]
	}

	if fileID == "" {
		return nil, fmt.Errorf("could not extract file or folder ID from Google Drive link")
	}

	link := &GoogleDriveLink{
		FileID:  fileID,
		IsFolder: isFolder,
		ViewURL: fullURL,
	}

	if isFolder {
		link.EmbedURL = fmt.Sprintf("https://drive.google.com/embeddedfolderview?id=%s", fileID)
		link.DownloadURL = fmt.Sprintf("https://drive.google.com/drive/folders/%s", fileID)
	} else {
		link.EmbedURL = fmt.Sprintf("https://drive.google.com/file/d/%s/preview", fileID)
		link.DownloadURL = fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID)
		if !strings.Contains(fullURL, "/view") && !strings.Contains(fullURL, "/preview") {
			link.ViewURL = fmt.Sprintf("https://drive.google.com/file/d/%s/view", fileID)
		}
	}

	return link, nil
}

func NormalizeGoogleDriveURL(rawURL string) (string, error) {
	link, err := ParseGoogleDriveURL(rawURL)
	if err != nil {
		return "", err
	}
	return link.ViewURL, nil
}
