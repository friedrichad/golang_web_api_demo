package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Map[T, U any](data []T, f func(T) U) []U {
	res := make([]U, 0, len(data))
	for _, e := range data {
		res = append(res, f(e))
	}
	return res
}

func Pointer[T any](val T) *T {
	return &val
}

func StrToInt(val string, defaultVal int) int {
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

func SliceStrToInt(val []string) []int {
	return Map(val, func(s string) int {
		return StrToInt(s, 0)
	})
}

func Contains(array []string, value string) bool {
	set := Slice2Set(array)
	return set[value]
}

func Slice2Set(array []string) map[string]bool {
	set := make(map[string]bool)
	for _, v := range array {
		set[v] = true
	}
	return set
}

func AnyContains(array []string, value []string) bool {
	set := Slice2Set(array)
	for _, v := range value {
		if set[v] {
			return true
		}
	}
	return false
}

func Filter[T any](data []T, f func(T) bool) []T {
	arr := make([]T, 0, len(data))
	for _, e := range data {
		if f(e) {
			arr = append(arr, e)
		}
	}
	return arr
}

func BuildTree[T any](slice []T, isParent func(m T) bool, compare func(p T, c T) bool, setTree func(p *T, c []T)) []T {
	tree := Filter(slice, isParent)
	for i, _ := range tree {
		buildTree(&tree[i], slice, compare, setTree)
	}
	return tree
}

func buildTree[T any](parent *T, slice []T, compare func(p T, c T) bool, setTree func(p *T, c []T)) {
	var children = Filter(slice, func(c T) bool {
		return compare(*parent, c)
	})
	if len(children) == 0 {
		return
	}
	for i, _ := range children {
		buildTree(&children[i], slice, compare, setTree)
	}
	setTree(parent, children)
}

func PageSize(c *gin.Context) (int, int) {
	page := c.Query("page")
	size := c.Query("size")
	return StrToInt(page, 0), StrToInt(size, 10)
}

func DateTimeFormat(t time.Time) string {
	return t.Format("02/01/2006 15:04:05")
}

func DateFormat(t time.Time) string {
	return t.Format("02/01/2006")
}

func TrunDate(t *time.Time, toStart bool) *time.Time {
	if t == nil || t.IsZero() {
		return nil
	}
	if toStart {
		return Pointer(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()))
	}
	return Pointer(time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()))
}
