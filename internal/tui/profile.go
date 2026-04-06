package tui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"toofan/internal/theme"
)

type profileData struct {
	Tests     int
	Time      time.Duration
	Best      map[string]map[int]float64 // mode -> dur -> wpm
	Recent    []testEntry
	Activity  map[string]int
	RecentAvg float64
}

type testEntry struct {
	Date   time.Time
	WPM    float64
	Dur    int
	Acc    float64
	Mode   string
	Raw    float64
	Errors int
}

func loadProfile() profileData {
	pd := profileData{
		Best:     make(map[string]map[int]float64),
		Activity: make(map[string]int),
	}
	pd.Best["words"] = make(map[int]float64)
	pd.Best["code"] = make(map[int]float64)

	home, err := os.UserHomeDir()
	if err != nil {
		return pd
	}
	dataDir := filepath.Join(home, ".toofan")

	f, err := os.Open(filepath.Join(dataDir, "results.txt"))
	if err != nil {
		return pd
	}
	defer f.Close()

	var all []testEntry
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		e, ok := parseResultLine(sc.Text())
		if !ok {
			continue
		}
		all = append(all, e)
		pd.Tests++
		pd.Time += time.Duration(e.Dur) * time.Second

		mode := e.Mode
		if strings.HasPrefix(e.Mode, "code:") {
			mode = "code"
		}
		if pd.Best[mode] == nil {
			pd.Best[mode] = make(map[int]float64)
		}
		if e.WPM > pd.Best[mode][e.Dur] {
			pd.Best[mode][e.Dur] = e.WPM
		}
		pd.Activity[e.Date.Format("2006-01-02")]++
	}

	if len(all) > 80 {
		pd.Recent = all[len(all)-80:]
	} else {
		pd.Recent = all
	}

	if len(all) >= 10 {
		sum := 0.0
		for i := len(all) - 10; i < len(all); i++ {
			sum += all[i].WPM
		}
		pd.RecentAvg = sum / 10.0
	} else if len(all) > 0 {
		sum := 0.0
		for _, e := range all {
			sum += e.WPM
		}
		pd.RecentAvg = sum / float64(len(all))
	}

	return pd
}

func parseResultLine(line string) (testEntry, bool) {
	parts := strings.Split(line, "|")
	if len(parts) < 5 {
		return testEntry{}, false
	}

	date, err := time.Parse("2006-01-02 15:04", strings.TrimSpace(parts[0]))
	if err != nil {
		return testEntry{}, false
	}

	wpmStr := strings.TrimSpace(parts[1])
	wpmStr = strings.TrimSuffix(wpmStr, "wpm")
	wpmStr = strings.TrimSpace(wpmStr)
	wpm, _ := strconv.ParseFloat(wpmStr, 64)

	accStr := strings.TrimSpace(parts[2])
	accStr = strings.TrimSuffix(accStr, "%")
	accStr = strings.TrimSpace(accStr)
	acc, _ := strconv.ParseFloat(accStr, 64)

	durStr := strings.TrimSpace(parts[3])
	durStr = strings.TrimSuffix(durStr, "s")
	durStr = strings.TrimSpace(durStr)
	dur, _ := strconv.Atoi(durStr)

	modeStr := strings.TrimSpace(parts[4])

	var raw float64
	var errors int
	if len(parts) >= 6 {
		rawStr := strings.TrimSpace(parts[5])
		rawStr = strings.TrimSuffix(rawStr, "raw")
		rawStr = strings.TrimSpace(rawStr)
		raw, _ = strconv.ParseFloat(rawStr, 64)
	}
	if len(parts) >= 7 {
		errStr := strings.TrimSpace(parts[6])
		errStr = strings.TrimSuffix(errStr, "err")
		errStr = strings.TrimSpace(errStr)
		errors, _ = strconv.Atoi(errStr)
	}

	return testEntry{Date: date, WPM: wpm, Dur: dur, Acc: acc, Mode: modeStr, Raw: raw, Errors: errors}, true
}

func (m model) handleProfile(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.active = screenTyping
	return m, nil
}

func cleanLegacyLang(lang string) string {
	for _, suffix := range []string{"15s", "30s", "60s", "120s"} {
		lang = strings.TrimSuffix(lang, suffix)
	}
	switch lang {
	case "javascript":
		return "js"
	case "typescript":
		return "ts"
	}
	if len(lang) > 9 {
		return lang[:9]
	}
	return lang
}

func rank(wpm float64) string {
	switch {
	case wpm >= 120:
		return "phantom"
	case wpm >= 80:
		return "demon"
	case wpm >= 50:
		return "survivor"
	case wpm >= 30:
		return "warrior"
	default:
		return "grandma"
	}
}

func (m model) bestRow(label string, data map[int]float64, dim, val, hi lipgloss.Style) string {
	const cw = 5
	const lw = 7
	cells := []string{col(lw, hi.Render(label))}
	for _, d := range []int{15, 30, 60, 120} {
		if wpm, ok := data[d]; ok {
			cells = append(cells, col(cw, val.Render(fmt.Sprintf("%.0f", wpm))))
		} else {
			cells = append(cells, col(cw, dim.Render("-")))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, cells...)
}

func (m model) viewProfile(p theme.Palette) string {
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	val := lipgloss.NewStyle().Foreground(p.Typed).Bold(true)
	hi := lipgloss.NewStyle().Foreground(p.Accent)

	rankLabel := ""
	if m.prof.Tests > 0 {
		rankLabel = " · " + rank(m.prof.RecentAvg)
	}
	title := val.Render("_toofan" + rankLabel)

	fullWidth := 76
	if m.width > 0 && m.width < 82 {
		fullWidth = m.width - 6
	}
	paneWidth := (fullWidth - 2) / 3 // 2 gaps of 1 char each

	paneStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.Foreground).
		Padding(1, 2)

	hours := int(m.prof.Time.Hours())
	mins := int(m.prof.Time.Minutes()) % 60
	timeStr := fmt.Sprintf("%dm", mins)
	if hours > 0 {
		timeStr = fmt.Sprintf("%dh %dm", hours, mins)
	}

	avgStr := dim.Render("-")
	if m.prof.Tests > 0 {
		avgStr = val.Render(fmt.Sprintf("%.0f wpm", m.prof.RecentAvg))
	}

	avgAccStr := dim.Render("-")
	if len(m.prof.Recent) > 0 {
		var totalAcc float64
		for _, e := range m.prof.Recent {
			totalAcc += e.Acc
		}
		avgAccStr = val.Render(fmt.Sprintf("%.0f%%", totalAcc/float64(len(m.prof.Recent))))
	}

	overview := lipgloss.JoinVertical(lipgloss.Left,
		hi.Render("overview"),
		"",
		dim.Render("tests  ")+val.Render(fmt.Sprintf("%d", m.prof.Tests)),
		dim.Render("time   ")+val.Render(timeStr),
		dim.Render("avg    ")+avgStr,
		dim.Render("acc    ")+avgAccStr,
	)

	const cw = 5
	const lw = 7

	durLabels := lipgloss.JoinHorizontal(lipgloss.Left,
		col(lw, ""),
		col(cw, dim.Render("15s")),
		col(cw, dim.Render("30s")),
		col(cw, dim.Render("60s")),
		col(cw, dim.Render("120s")),
	)

	wordsLine := m.bestRow("words", m.prof.Best["words"], dim, val, hi)
	codeLine := m.bestRow("code", m.prof.Best["code"], dim, val, hi)

	bests := lipgloss.JoinVertical(lipgloss.Left,
		hi.Render("personal bests"),
		"",
		durLabels,
		wordsLine,
		codeLine,
	)

	cur := rank(m.prof.RecentAvg)
	type tier struct {
		name string
		wpm  int
	}
	tiers := []tier{
		{"grandma", 0}, {"warrior", 30}, {"survivor", 50}, {"demon", 80}, {"phantom", 120},
	}

	var rankLines []string
	for _, t := range tiers {
		line := fmt.Sprintf("%-8s %3d+", t.name, t.wpm)
		if t.name == cur {
			rankLines = append(rankLines, hi.Render("●")+" "+val.Render(line))
		} else {
			rankLines = append(rankLines, dim.Render("○ "+line))
		}
	}

	ranks := lipgloss.JoinVertical(lipgloss.Left,
		hi.Render("ranks"),
		"",
		strings.Join(rankLines, "\n"),
	)

	// Render all boxes first to measure actual heights
	overviewBox := paneStyle.Width(paneWidth).Render(overview)
	bestBox := paneStyle.Width(paneWidth).Render(bests)
	ranksBox := paneStyle.Width(paneWidth).Render(ranks)

	// Match heights
	maxH := lipgloss.Height(overviewBox)
	if h := lipgloss.Height(bestBox); h > maxH {
		maxH = h
	}
	if h := lipgloss.Height(ranksBox); h > maxH {
		maxH = h
	}

	overviewBox = paneStyle.Width(paneWidth).Height(maxH - 2).Render(overview)
	bestBox = paneStyle.Width(paneWidth).Height(maxH - 2).Render(bests)
	ranksBox = paneStyle.Width(paneWidth).Height(maxH - 2).Render(ranks)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, overviewBox, " ", bestBox, " ", ranksBox)

	var histRows []string
	header := lipgloss.JoinHorizontal(lipgloss.Left,
		col(6, hi.Render("wpm")),
		col(6, hi.Render("raw")),
		col(6, hi.Render("acc")),
		col(5, hi.Render("err")),
		col(7, hi.Render("type")),
		col(9, hi.Render("lang")),
		col(6, hi.Render("dur")),
		col(13, hi.Render("date")),
	)
	histRows = append(histRows, header, "")

	limit := 10
	if len(m.prof.Recent) < limit {
		limit = len(m.prof.Recent)
	}

	for i := len(m.prof.Recent) - 1; i >= len(m.prof.Recent)-limit; i-- {
		e := m.prof.Recent[i]
		dstr := e.Date.Format("02 Jan 15:04")

		modeType := "words"
		modeLang := "english"
		if strings.HasPrefix(e.Mode, "code:") {
			modeType = "code"
			modeLang = cleanLegacyLang(strings.TrimPrefix(e.Mode, "code:"))
		}

		durStr := "∞"
		if e.Dur > 0 {
			durStr = fmt.Sprintf("%ds", e.Dur)
		}

		row := lipgloss.JoinHorizontal(lipgloss.Left,
			col(6, val.Render(fmt.Sprintf("%.0f", e.WPM))),
			col(6, dim.Render(fmt.Sprintf("%.0f", e.Raw))),
			col(6, dim.Render(fmt.Sprintf("%.0f%%", e.Acc))),
			col(5, dim.Render(fmt.Sprintf("%d", e.Errors))),
			col(7, dim.Render(modeType)),
			col(9, dim.Render(modeLang)),
			col(6, dim.Render(durStr)),
			col(13, dim.Render(dstr)),
		)
		histRows = append(histRows, row)
	}

	histBox := paneStyle.Width(fullWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			hi.Render("recent tests"),
			"",
			lipgloss.JoinVertical(lipgloss.Left, histRows...),
		),
	)

	heatmapStr := heatGrid(m.prof.Activity, p, fullWidth)
	heatBox := paneStyle.Width(fullWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			hi.Render("activity map"),
			"",
			heatmapStr,
		),
	)

	body := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		topRow,
		"",
		histBox,
		"",
		heatBox,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body)
}

func heatGrid(activity map[string]int, p theme.Palette, width int) string {
	now := time.Now()
	days := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
	dows := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}

	c0 := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))
	c1 := lipgloss.NewStyle().Foreground(p.Accent)
	dim := lipgloss.NewStyle().Foreground(p.Foreground)

	weeks := (width - 14) / 2
	if weeks < 1 {
		weeks = 1
	}
	if weeks > 26 {
		weeks = 26
	}

	var rows []string
	for i, dow := range dows {
		var row strings.Builder
		row.WriteString(dim.Render(fmt.Sprintf("%3s  ", days[i])))
		for w := weeks - 1; w >= 0; w-- {
			d := now.AddDate(0, 0, -w*7)
			for d.Weekday() != dow {
				d = d.AddDate(0, 0, -1)
			}
			if activity[d.Format("2006-01-02")] > 0 {
				row.WriteString(c1.Render("■") + " ")
			} else {
				row.WriteString(c0.Render("■") + " ")
			}
		}
		rows = append(rows, row.String())
	}
	return strings.Join(rows, "\n")
}
