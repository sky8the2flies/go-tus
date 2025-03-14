package tus

import (
	"bytes"
)

type Uploader struct {
	client     *Client
	url        string
	upload     *Upload
	offset     int64
	curoffset  int64
	aborted    bool
	uploadSubs []chan Upload
	notifyChan chan bool
}

// Subscribes to progress updates.
func (u *Uploader) NotifyUploadProgress(c chan Upload) {
	u.uploadSubs = append(u.uploadSubs, c)
}

// Abort aborts the upload process.
// It doens't abort the current chunck, only the remaining.
func (u *Uploader) Abort() {
	u.aborted = true
}

// IsAborted returns true if the upload was aborted.
func (u *Uploader) IsAborted() bool {
	return u.aborted
}

// Url returns the upload url.
func (u *Uploader) Url() string {
	return u.url
}

// Offset returns the current offset uploaded.
func (u *Uploader) Offset() int64 {
	return u.curoffset
}

// Upload uploads the entire body to the server.
func (u *Uploader) Upload() error {
	for u.curoffset < u.upload.size && !u.aborted {
		err := u.UploadChunck()

		if err != nil {
			return err
		}
	}

	u.upload.stream = nil

	return nil
}

// UploadChunck uploads a single chunck.
func (u *Uploader) UploadChunck() error {
	data := make([]byte, u.offset)
	_, err := u.upload.stream.Seek(u.curoffset, 0)
	if err != nil {
		return err
	}

	size, err := u.upload.stream.Read(data)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(data[:size])
	defer body.Reset()
	u.curoffset, err = u.client.uploadChunck(u.url, body, int64(size), u.offset)
	if err != nil {
		return err
	}

	u.curoffset -= u.offset

	u.upload.updateProgress(u.curoffset)

	u.notifyChan <- true

	data = nil

	return nil
}

// Waits for a signal to broadcast to all subscribers
func (u *Uploader) broadcastProgress() {
	for _ = range u.notifyChan {
		for _, c := range u.uploadSubs {
			c <- *u.upload
		}
	}
}

// NewUploader creates a new Uploader.
func NewUploader(client *Client, url string, upload *Upload, offset int64, curoffset int64) *Uploader {
	notifyChan := make(chan bool)

	uploader := &Uploader{
		client,
		url,
		upload,
		offset,
		curoffset,
		false,
		nil,
		notifyChan,
	}

	go uploader.broadcastProgress()

	return uploader
}
