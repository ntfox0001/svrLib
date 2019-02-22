package util

// go没有模板，特殊slice只能拷贝一份到自己代码里了
func SSliceDel(s []string, id int) []string {
	if s == nil {
		return s
	}
	if len(s) <= id {
		return s
	}

	t := append(s[:id], s[id+1:]...)
	return t
}
func ISliceDel(s []int, id int) []int {
	if s == nil {
		return s
	}
	if len(s) <= id {
		return s
	}

	t := append(s[:id], s[id+1:]...)
	return t
}
