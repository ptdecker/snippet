/*
 * Custom file server to prevent browsing of directories but preserve
 * traditional index.html behavior
 *
 * c.f. https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
 */

package main

import (
	"net/http"
	"strings"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := nfs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
