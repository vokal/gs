/****************************************************************************
 * Copyright (c) 2013, VOKAL Interactive, LLC
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the software nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY VOKAL INTERACTIVE, LLC ''AS IS'' AND ANY
 * EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL VOKAL INTERACTIVE, LLC BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 ****************************************************************************/

package gs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (u *Uploader) Do(f *File) error {
	urls := fmt.Sprintf("http://storage.googleapis.com/%s/%s", f.Bucket, f.Path)
	params := make(url.Values)

	urls += "?" + params.Encode()

	req, err := http.NewRequest("PUT", urls, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	req.Header.Set("Content-Type", http.DetectContentType(f.Object))
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
