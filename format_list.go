package ytdl

import "sort"

// FormatList is a slice of formats with filtering functionality
type FormatList []Format

// Filter filters
func (formats FormatList) Filter(key FormatKey, values []interface{}) FormatList {
	var dst FormatList
	for _, v := range values {
		for _, f := range formats {
			if interfaceToString(f.ValueForKey(key)) == interfaceToString(v) {
				dst = append(dst, f)
			}
		}
	}

	return dst
}

// Extremes filters by extremes
func (formats FormatList) Extremes(key FormatKey, best bool) FormatList {
	dst := formats.Copy()
	if len(dst) > 1 {
		dst.Sort(key, best)
		first := dst[0]
		var i int
		for i = 0; i < len(dst)-1; i++ {
			if first.CompareKey(dst[i+1], key) != 0 {
				break
			}
		}
		i++
		dst = dst[0:i]
	}
	return dst
}

// Best filters by extremes for the best
func (formats FormatList) Best(key FormatKey) FormatList {
	return formats.Extremes(key, true)
}

// Worst filters by extremes for the worst
func (formats FormatList) Worst(key FormatKey) FormatList {
	return formats.Extremes(key, false)
}

// Sort sorts by format key
func (formats FormatList) Sort(key FormatKey, reverse bool) {
	wrapper := formatsSortWrapper{formats, key}
	if !reverse {
		sort.Stable(wrapper)
	} else {
		sort.Stable(sort.Reverse(wrapper))
	}
}

// Subtract filters by subtracting one FormatList by another
func (formats FormatList) Subtract(other FormatList) FormatList {
	var dst FormatList
	for _, f := range formats {
		include := true
		for _, f2 := range other {
			if f2.Itag == f.Itag {
				include = false
				break
			}
		}
		if include {
			dst = append(dst, f)
		}
	}
	return dst
}

// Copy copies the FormatList and returns
func (formats FormatList) Copy() FormatList {
	dst := make(FormatList, len(formats))
	copy(dst, formats)
	return dst
}

type formatsSortWrapper struct {
	formats FormatList
	key     FormatKey
}

func (s formatsSortWrapper) Len() int {
	return len(s.formats)
}

func (s formatsSortWrapper) Less(i, j int) bool {
	return s.formats[i].CompareKey(s.formats[j], s.key) < 0
}

func (s formatsSortWrapper) Swap(i, j int) {
	s.formats[i], s.formats[j] = s.formats[j], s.formats[i]
}
