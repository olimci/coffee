package coffee

import (
	"fmt"
	"slices"
	"strings"
)

func (m *model) View() string {
	if m.height <= 0 {
		return strings.Join(slices.Concat(
			m.sectionLines(m.header),
			m.sectionLines(m.body),
			m.sectionLines(m.footer),
		), "\n") + "\n"
	}

	focused, _ := m.focused()
	header := renderSection(m.header, focused, m.height, m.width)

	footerMax := max(0, m.height-len(header))
	footer := renderSection(m.footer, focused, footerMax, m.width)

	bodyMax := max(0, m.height-len(header)-len(footer))
	body := renderSection(m.body, focused, bodyMax, m.width)

	lines := make([]string, 0, len(header)+len(body)+len(footer))
	lines = append(lines, header...)
	lines = append(lines, body...)

	if m.altScreen && len(footer) > 0 {
		for range max(0, m.height-len(lines)-len(footer)) {
			lines = append(lines, "")
		}
	}

	lines = append(lines, footer...)
	return strings.Join(lines, "\n") + "\n"
}

func (m *model) sectionLines(items []item) []string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		if item.entry != nil {
			lines = append(lines, item.entry.submodel.View())
		} else {
			lines = append(lines, item.text)
		}
	}
	return lines
}

func renderSection(items []item, focused *submodelEntry, height, width int) []string {
	lines, focusStart, focusEnd := flattenItems(items, focused)
	if height <= 0 {
		return nil
	}
	if len(lines) <= height {
		return lines
	}

	start := len(lines) - height
	end := len(lines)

	if focusStart >= 0 {
		focusHeight := focusEnd - focusStart
		switch {
		case focusHeight >= height:
			start = focusStart
			end = start + height
		case focusEnd > end:
			end = focusEnd
			start = end - height
		case focusStart < start:
			start = focusStart
			end = start + height
		}
	}

	if start < 0 {
		start = 0
	}
	if end > len(lines) {
		end = len(lines)
	}
	if end-start > height {
		end = start + height
	}

	out := slices.Clone(lines[start:end])
	omittedBefore := start
	omittedAfter := len(lines) - end

	switch {
	case omittedBefore > 0 && omittedAfter > 0 && len(out) == 1:
		out[0] = omissionLine(omittedBefore+omittedAfter, width)
	case omittedBefore > 0 && len(out) > 0:
		out[0] = omissionLine(omittedBefore, width)
		if omittedAfter > 0 && len(out) > 1 {
			out[len(out)-1] = omissionLine(omittedAfter, width)
		}
	case omittedAfter > 0 && len(out) > 0:
		out[len(out)-1] = omissionLine(omittedAfter, width)
	}

	return out
}

func flattenItems(items []item, focused *submodelEntry) ([]string, int, int) {
	lines := make([]string, 0, len(items))
	focusStart, focusEnd := -1, -1

	for _, item := range items {
		itemLines := itemLines(item)
		start := len(lines)
		lines = append(lines, itemLines...)

		if item.entry != nil && item.entry == focused {
			focusStart = start
			focusEnd = len(lines)
		}
	}

	return lines, focusStart, focusEnd
}

func itemLines(item item) []string {
	if item.entry != nil {
		return strings.Split(item.entry.submodel.View(), "\n")
	}
	return strings.Split(item.text, "\n")
}

func omissionLine(count, width int) string {
	label := fmt.Sprintf("(%d lines omitted)", count)
	if width <= len(label)+4 {
		return label
	}

	left := (width - len(label) - 4) / 2
	right := width - len(label) - 4 - left
	return strings.Repeat("─", left) + "┤ " + label + " ├" + strings.Repeat("─", right)
}
