package gs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Uploader struct {
	ProjectId string
	client    *http.Client
}

type File struct {
	Path   string
	Bucket string
	Object []byte
}

func NewUploader(scope, id string) *Uploader {
	return &Uploader{
		ProjectId: id,
		client:    newOauthClient(scope),
	}
}

func (u *Uploader) getMimetype(filename string) string {
	extension := strings.Split(filename, ".")[1]

	switch extension {
	case "jpg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "html":
		return "text/html"
	case "css":
		return "text/css"
	case "js":
		return "text/javascript"
	}

	return "text/plain"
}

func (u *Uploader) Do(f *File) error {
	urls := fmt.Sprintf("http://storage.googleapis.com/%s/%s", f.Bucket, f.Path)
	params := make(url.Values)

	urls += "?" + params.Encode()

	req, err := http.NewRequest("PUT", urls, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	req.Header.Set("Content-Type", u.getMimetype(f.Path))
	req.Header.Set("x-goog-project-id", u.ProjectId)

	body := ioutil.NopCloser(bytes.NewReader(f.Object))
	req.Body = body
	req.ContentLength = int64(len(f.Object))

	resp, err := u.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
