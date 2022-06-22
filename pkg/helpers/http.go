package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

func HttpGetContentRange(header string) (int, int, int, error) {
	tmp := strings.Split(header, " ")
	if len(tmp) != 2 {
		return 0, 0, 0, fmt.Errorf("invalid content-range header")
	}

	if tmp[0] != "bytes" {
		return 0, 0, 0, fmt.Errorf("invalid content-range header")
	}

	rangeAndSize := strings.Split(tmp[1], "/")
	if len(rangeAndSize) != 2 {
		return 0, 0, 0, fmt.Errorf("invalid content-range header")
	}

	rangeParts := strings.Split(rangeAndSize[0], "-")
	if len(rangeParts) != 2 {
		return 0, 0, 0, fmt.Errorf("invalid content-range header")
	}

	rangeStart, err := strconv.Atoi(rangeParts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid content-range header")
	}

	rangeEnd, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid content-range header")
	}

	fileSize, err := strconv.Atoi(rangeAndSize[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid content-range header")
	}

	return rangeStart, rangeEnd, fileSize, nil
}

func HttpSplitPath(path string) []string {
	return StringArrayDeleteEmpty(strings.Split(path, "/"))
}
