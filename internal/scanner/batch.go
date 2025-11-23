package scanner

func BatchByLines(files []Info, maxLines int) [][]Info {
	var result [][]Info
	var group []Info
	used := 0

	for _, f := range files {
		if used+f.Lines > maxLines && len(group) > 0 {
			result = append(result, group)
			group = nil
			used = 0
		}
		group = append(group, f)
		used += f.Lines
	}

	if len(group) > 0 {
		result = append(result, group)
	}

	return result
}
