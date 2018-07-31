package helper

// https://golang.org/src/path/path.go?s=4371:4399#L158
func FileAndExt(path string) (string, string) {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[:i], path[i:]
		}
	}
	return "", ""
}
