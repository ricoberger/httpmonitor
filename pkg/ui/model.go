package ui

import (
	"fmt"
	"time"

	"github.com/ricoberger/httpmonitor/pkg/target"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type TickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

var tableBorder = table.Border{
	Top:    "─",
	Left:   "│",
	Right:  "│",
	Bottom: "─",

	TopRight:    "╮",
	TopLeft:     "╭",
	BottomRight: "╯",
	BottomLeft:  "╰",

	TopJunction:    "╥",
	LeftJunction:   "├",
	RightJunction:  "┤",
	BottomJunction: "╨",
	InnerJunction:  "╫",

	InnerDivider: "║",
}

type Model struct {
	ScreenWidth  int
	ScreenHeight int

	ActiveTargetIndex int
	Targets           []target.Client

	TargetsTable table.Model
	ResultsTable table.Model
}

func NewModel(targets []target.Client) Model {
	activeTargetIndex := -1
	if len(targets) == 1 {
		activeTargetIndex = 0
	}

	return Model{
		ScreenWidth:       0,
		ScreenHeight:      0,
		ActiveTargetIndex: activeTargetIndex,
		Targets:           targets,
		TargetsTable: table.New([]table.Column{
			table.NewFlexColumn("name", "Name", 100).WithStyle(lipgloss.NewStyle().Align(lipgloss.Left)),
			table.NewFlexColumn("dnslookup", "DNS Lookup", 20),
			table.NewFlexColumn("tcpconnection", "TCP Connection", 20),
			table.NewFlexColumn("tlshandshake", "TLS Handshake", 20),
			table.NewFlexColumn("serverprocessing", "Server Processing", 20),
			table.NewFlexColumn("contenttransfer", "Content Transfer", 20),
			table.NewFlexColumn("total", "Total", 20),
			table.NewFlexColumn("status", "Status", 20),
		}).
			Focused(true).
			HeaderStyle(lipgloss.NewStyle().Bold(true)).
			HighlightStyle(lipgloss.NewStyle().Background(lipgloss.Color("12")).Foreground(lipgloss.Color("8"))).
			Border(tableBorder).
			WithPageSize(10).
			WithStaticFooter("Page 1/1"),
		ResultsTable: table.New([]table.Column{
			table.NewFlexColumn("time", "Time", 100).WithStyle(lipgloss.NewStyle().Align(lipgloss.Left)),
			table.NewFlexColumn("dnslookup", "DNS Lookup", 20),
			table.NewFlexColumn("tcpconnection", "TCP Connection", 20),
			table.NewFlexColumn("tlshandshake", "TLS Handshake", 20),
			table.NewFlexColumn("serverprocessing", "Server Processing", 20),
			table.NewFlexColumn("contenttransfer", "Content Transfer", 20),
			table.NewFlexColumn("total", "Total", 20),
			table.NewFlexColumn("status", "Status", 20),
		}).
			Focused(true).
			HeaderStyle(lipgloss.NewStyle().Bold(true)).
			HighlightStyle(lipgloss.NewStyle().Background(lipgloss.Color("12")).Foreground(lipgloss.Color("8"))).
			Border(tableBorder).
			WithPageSize(10).
			WithStaticFooter("Page 1/1"),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tickEvery())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.TargetsTable, cmd = m.TargetsTable.Update(msg)
	cmds = append(cmds, cmd)
	m.ResultsTable, _ = m.ResultsTable.Update(msg)
	cmds = append(cmds, cmd)

	if m.ActiveTargetIndex == -1 {
		m.TargetsTable = m.TargetsTable.WithStaticFooter(fmt.Sprintf("Page %d/%d", m.TargetsTable.CurrentPage(), m.TargetsTable.MaxPages()))
	} else {
		m.ResultsTable = m.ResultsTable.WithStaticFooter(fmt.Sprintf("Page %d/%d", m.ResultsTable.CurrentPage(), m.ResultsTable.MaxPages()))
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)
		case "enter":
			if m.ActiveTargetIndex == -1 {
				m.ActiveTargetIndex = m.TargetsTable.GetHighlightedRowIndex()
			}
		case "esc":
			if len(m.Targets) > 1 {
				m.ActiveTargetIndex = -1
			}
		}
	case tea.WindowSizeMsg:
		m.ScreenWidth = msg.Width
		m.ScreenHeight = msg.Height

		m.TargetsTable = m.TargetsTable.WithTargetWidth(msg.Width).WithMinimumHeight(msg.Height).WithPageSize(msg.Height - 6)
		m.ResultsTable = m.ResultsTable.WithTargetWidth(msg.Width).WithMinimumHeight(msg.Height).WithPageSize(msg.Height - 6)

	case TickMsg:
		if m.ActiveTargetIndex == -1 {
			m.TargetsTable = m.TargetsTable.WithRows(m.generateTargetsTableRows())
		} else {
			m.ResultsTable = m.ResultsTable.WithRows(m.generateResultsTableRows())
		}
		cmds = append(cmds, tickEvery())
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.ScreenWidth == 0 || m.ScreenHeight == 0 {
		return "Loading..."
	}

	if m.ActiveTargetIndex == -1 {
		return m.TargetsTable.View()
	}

	return m.ResultsTable.View()
}

func (m *Model) generateTargetsTableRows() []table.Row {
	var rows []table.Row

	for _, target := range m.Targets {
		name := target.Name()
		result := target.LastResult()

		statusColor := lipgloss.Color("10")
		if result.StatusCode == 0 || result.StatusCode >= 500 {
			statusColor = lipgloss.Color("9")
		} else if result.StatusCode >= 400 {
			statusColor = lipgloss.Color("11")
		}

		rows = append(rows,
			table.NewRow(table.RowData{
				"name":             name,
				"dnslookup":        result.DNSLookup.String(),
				"tcpconnection":    result.TCPConnection.String(),
				"tlshandshake":     result.TLSHandshake.String(),
				"serverprocessing": result.ServerProcessing.String(),
				"contenttransfer":  result.ContentTransfer.String(),
				"total":            result.Total.String(),
				"status":           table.NewStyledCell(fmt.Sprintf("%d", result.StatusCode), lipgloss.NewStyle().Background(statusColor).Foreground(lipgloss.Color("8")).Bold(true)),
			}),
		)
	}

	return rows
}

func (m *Model) generateResultsTableRows() []table.Row {
	var rows []table.Row

	target := m.Targets[m.ActiveTargetIndex]
	results := target.Results()

	for i := len(results) - 1; i >= 0; i-- {
		result := results[i]
		time := result.StartTime().Format("2006-01-02 15:04:05")

		statusColor := lipgloss.Color("10")
		if result.StatusCode == 0 || result.StatusCode >= 500 {
			statusColor = lipgloss.Color("9")
		} else if result.StatusCode >= 400 {
			statusColor = lipgloss.Color("11")
		}

		rows = append(rows,
			table.NewRow(table.RowData{
				"time":             time,
				"dnslookup":        result.DNSLookup.String(),
				"tcpconnection":    result.TCPConnection.String(),
				"tlshandshake":     result.TLSHandshake.String(),
				"serverprocessing": result.ServerProcessing.String(),
				"contenttransfer":  result.ContentTransfer.String(),
				"total":            result.Total.String(),
				"status":           table.NewStyledCell(fmt.Sprintf("%d", result.StatusCode), lipgloss.NewStyle().Background(statusColor).Foreground(lipgloss.Color("8")).Bold(true)),
			}),
		)
	}

	return rows
}
