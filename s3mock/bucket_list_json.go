// GENERATED FILE: Do not edit, all changes will be lost.

// Package s3mock contains the content of files from BucketListJSON. It
// is generated by:
//     https://github.com/aybabtme/gostatic
package s3mock

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"log"
)

// GetBucketListJSON will lookup the static assets. It returns a *bytes.Reader
// and true if found, false otherwise. The static assets contain exactly the
// following entries:
//
//   bucket_list.json
//
func GetBucketListJSON(filename string) (*bytes.Reader, bool) {
	data, ok := decompressedBucketListJSON[filename]
	return bytes.NewReader(data), ok
}

// ListBucketListJSON will return all the static assets sharing root
// BucketListJSON.
func ListBucketListJSON() map[string]*bytes.Reader {
	out := make(map[string]*bytes.Reader, len(decompressedBucketListJSON))
	for k, v := range decompressedBucketListJSON {
		out[k] = bytes.NewReader(v)
	}
	return out
}

var decompressedBucketListJSON = make(map[string][]byte)

func init() {

	var compressed = [...]struct {
		name   string
		gzip64 string
	}{
		{"bucket_list.json", `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`},
	}

	for _, file := range compressed {
		gzipdata, err := base64.StdEncoding.DecodeString(file.gzip64)
		if err != nil {
			log.Panicf("Couldn't decode base64 data for %q: %v", file.name, err)
		}
		gr, err := gzip.NewReader(bytes.NewBuffer(gzipdata))
		if err != nil {
			log.Panicf("Couldn't open gzip stream for data for %q: %v", file.name, err)
		}
		data, err := ioutil.ReadAll(gr)
		if err != nil {
			log.Panicf("Couldn't decompress gzip data in %q: %v", file.name, err)
		}
		decompressedBucketListJSON[file.name] = data
	}
}
