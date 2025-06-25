package ui

import "github.com/charmbracelet/lipgloss"

var (
    userStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62")).Padding(0, 1)
    aiStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("48")).Background(lipgloss.Color("238")).Padding(0, 1)
    sepStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
    headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229"))
    statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Background(lipgloss.Color("236")).Padding(0,1)
)
