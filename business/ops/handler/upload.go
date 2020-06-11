package handler

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
	"os"
)

type upFile struct {
	*os.File
	hash.Hash
}

func (f upFile) Write(p []byte) (int, error) {
	if _, err := f.File.Write(p); err != nil {
		return -1, err
	}
	return f.Hash.Write(p)
}

// Upload 上传文件
func Upload(w http.ResponseWriter, r *http.Request) {
	f, h, err := r.FormFile("file_data")
	if err != nil {
		jFailed(w, http.StatusBadRequest, err.Error())
		return
	}
	defer f.Close()

	u := &upFile{Hash: md5.New()}
	url := "static/upload/" + h.Filename
	if u.File, err = os.Create(url); err != nil {
		jFailed(w, http.StatusBadRequest, err.Error())
		return
	}
	defer u.Close()
	io.Copy(u, f)

	jSuccess(w, map[string]string{"url": url,
		"md5": hex.EncodeToString(u.Hash.Sum(nil))})
}
