package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
)

// CompressToBase64 compresses data using gzip and encodes to base64
func CompressToBase64(data string) (string, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write([]byte(data)); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// DecompressFromBase64 decodes base64 and decompresses gzip data
func DecompressFromBase64(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	gz, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return "", err
	}
	defer gz.Close()
	result, err := io.ReadAll(gz)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
