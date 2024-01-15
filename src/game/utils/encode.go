package utils

import "io"
import "bytes"
import "encoding/base64"
import "compress/gzip"
import "errors"

import "github.com/tinne26/luckyfeet/src/game/utils/ch426"

func GzipAndEncodeAsBase64(data []byte) (string, error) {
	outBuffer := bytes.NewBuffer(nil)
	writer := gzip.NewWriter(outBuffer)
	n, err := writer.Write(data)
	if err != nil { return "", err }
	if n != len(data) { return "", errors.New("short write") }
	err = writer.Close()
	if err != nil { return "", err }
	gzippedData := outBuffer.Bytes()
	return base64.StdEncoding.EncodeToString(gzippedData), nil
}

func DecodeFromBase64AndUngzip(data string) ([]byte, error) {
	gzippedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil { return nil, err }
	reader, err := gzip.NewReader(bytes.NewBuffer(gzippedBytes))
	if err != nil { return nil, err }
	defer reader.Close()
	return io.ReadAll(reader)
}

func GzipAndEncodeAsCh426(data []byte) (string, error) {
	outBuffer := bytes.NewBuffer(nil)
	writer := gzip.NewWriter(outBuffer)
	n, err := writer.Write(data)
	if err != nil { return "", err }
	if n != len(data) { return "", errors.New("short write") }
	err = writer.Close()
	if err != nil { return "", err }
	gzippedData := outBuffer.Bytes()
	return ch426.Encode(gzippedData), nil
}

func DecodeFromCh426AndUngzip(data string) ([]byte, error) {
	gzippedBytes, err := ch426.Decode(data)
	if err != nil { return nil, err }
	reader, err := gzip.NewReader(bytes.NewBuffer(gzippedBytes))
	if err != nil { return nil, err }
	defer reader.Close()
	return io.ReadAll(reader)
}
