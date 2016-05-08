package models

import (
	"bufio"
	"os"
)

func SplitFile(fn string) ([]byte, []byte, error) {
	var err error
	var file *os.File

	if file, err = os.Open(fn); err != nil {
		return nil, nil, err
	}
	defer file.Close()
	var h []byte
	var b []byte
	in := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "" && in {
			in = false
			continue
		}
		if in {
			h = append(h, scanner.Bytes()...)
			h = append(h, '\n')
		} else {
			b = append(b, scanner.Bytes()...)
			b = append(b, '\n')
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, nil, err
	}
	return h, b, nil
}
