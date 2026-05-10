package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	success = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	failure = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
	warn    = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	info    = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	dim     = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	heading = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)
	border  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cell    = lipgloss.NewStyle().Padding(0, 1)
	headRow = lipgloss.NewStyle().Padding(0, 1).Bold(true).Foreground(lipgloss.Color("99"))
)

func Success(format string, a ...any) {
	fmt.Fprintln(os.Stdout, success.Render(fmt.Sprintf(format, a...)))
}

func Error(w io.Writer, format string, a ...any) {
	fmt.Fprintln(w, failure.Render("error: ")+fmt.Sprintf(format, a...))
}

func Warn(format string, a ...any) {
	fmt.Fprintln(os.Stdout, warn.Render(fmt.Sprintf(format, a...)))
}

func Info(format string, a ...any) {
	fmt.Fprintln(os.Stdout, info.Render(fmt.Sprintf(format, a...)))
}

func Heading(s string) string { return heading.Render(s) }
func Dim(s string) string     { return dim.Render(s) }

// Table builds a styled table. statusCol, when >= 0, is the column index whose
// value should be colored green for SUCCESS / red for FAILURE.
func Table(headers []string, rows [][]string, statusCol int) string {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(border).
		Headers(headers...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headRow
			}
			if statusCol >= 0 && col == statusCol && row >= 0 && row < len(rows) {
				switch rows[row][col] {
				case "SUCCESS":
					return cell.Foreground(lipgloss.Color("42"))
				case "FAILURE":
					return cell.Foreground(lipgloss.Color("203"))
				}
			}
			return cell
		})
	return t.Render()
}
