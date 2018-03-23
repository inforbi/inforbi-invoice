package data

import "bufio"

func ReceiveBlob(reader bufio.Reader, length int) ([]byte) {
	file := make([]byte, 0, length)
	buf := make([]byte, 256)
	read := 0

	for read < length {
		n, err := reader.Read(buf)
		read += n
		file = append(file, buf[:n]...)
		if err != nil {

		}
	}
	return file
}
