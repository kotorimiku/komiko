package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ParamToUint(str string) (uint, error) {
	id64, err := strconv.ParseUint(str, 10, 32)
	return uint(id64), err
}

func RelToAbs(path string, relPath string) string {
	return strings.ReplaceAll(filepath.Clean(filepath.Join(filepath.Dir(relPath), string(path))), "\\", "/")
}

func Contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

var chineseNums = map[rune]int{
	'零': 0, '一': 1, '二': 2, '三': 3, '四': 4,
	'五': 5, '六': 6, '七': 7, '八': 8, '九': 9,
	'十': 10, '百': 100, '千': 1000, '万': 10000,
}

// 中文数字转阿拉伯数字（支持十、二十、三十一等简单格式）
func ChineseToInt(s string) int {
	total := 0
	unit := 1
	num := 0
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		c := runes[i]
		if val, ok := chineseNums[c]; ok {
			if val >= 10 {
				if val > unit {
					unit = val
				} else {
					unit *= val
				}
				if num == 0 {
					num = 1
				}
				total += num * unit
				num = 0
				unit = 1
			} else {
				num = val
			}
		}
	}
	total += num * unit
	return total
}

var patterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(?:第)?(\d+\.\d+|\d+)\s*[卷话章集]`),           // 第1卷、1话、第1章
	regexp.MustCompile(`(?i)卷\s*(\d+\.\d+|\d+)`),                      // 卷3
	regexp.MustCompile(`(?i)vol(?:ume)?[._-]?\s*(\d+\.\d+|\d+)`),      // vol.4, volume 4, vol-4, vol_4
	regexp.MustCompile(`(?i)t(?:ome)?[._-]?\s*(\d+\.\d+|\d+)`),        // tome 1 t.1 t-1 t_1
	regexp.MustCompile(`(?i)(?:第)?([一二三四五六七八九十百千万零]+)\s*[卷话章集]`),      // 第十二卷
	regexp.MustCompile(`(?i)卷([一二三四五六七八九十百千万零]+)`),                    // 卷十
	regexp.MustCompile(`(?i)vol(?:ume)?[._-]?\s*([一二三四五六七八九十百千万零]+)`), // vol.十
	regexp.MustCompile(`(?i)\[(\d+\.\d+|\d+)\]`),                      // [1]
	regexp.MustCompile(`(?i)\((\d+\.\d+|\d+)\)`),                      // (1)
	regexp.MustCompile(`(?i)(?:^|.|\s)(\d+\.\d+|\d+)[.\s]`),           // .1.
	regexp.MustCompile(`(?i)\[([一二三四五六七八九十百千万零]+)\]`),                 // [十]
	regexp.MustCompile(`(?i)\(([一二三四五六七八九十百千万零]+)\)`),                 // (十)
}

func ExtractVolume(s string) (float32, bool) {
	s = strings.ToLower(s)

	for _, re := range patterns {
		if match := re.FindStringSubmatch(s); len(match) >= 2 {
			if num, err := strconv.ParseFloat(match[1], 32); err == nil {
				return float32(num), true
			} else {
				return float32(ChineseToInt(match[1])), true
			}
		}
	}

	return 0, false
}

func FileName(path string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	nameOnly := filename[:len(filename)-len(ext)]
	return nameOnly
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// func structToMapExcludeZero[T any](obj *T) map[string]interface{} {
// 	result := make(map[string]interface{})
// 	val := reflect.ValueOf(obj).Elem()
// 	typ := val.Type()

// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Field(i)
// 		fieldType := typ.Field(i)

// 		// 跳过未导出字段
// 		if !field.CanInterface() {
// 			continue
// 		}

// 		// 判断是否为零值
// 		if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
// 			continue
// 		}

// 		// 获取 gorm tag
// 		tag := fieldType.Tag.Get("gorm")
// 		name := parseGormTagToColumn(tag, fieldType.Name)
// 		result[name] = field.Interface()
// 	}

// 	return result
// }

// func parseGormTagToColumn(tag string, defaultName string) string {
// 	// 简单解析 tag，如 `gorm:"column:username"` -> "username"
// 	parts := strings.Split(tag, ";")
// 	for _, part := range parts {
// 		if strings.HasPrefix(part, "column:") {
// 			return strings.TrimPrefix(part, "column:")
// 		}
// 	}
// 	return defaultName
// }
