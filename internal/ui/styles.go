package ui

import "github.com/charmbracelet/lipgloss"

// Color palette - Modern and accessible colors
var (
    // Primary colors
    primaryColor     = lipgloss.Color("#6366f1")   // Modern indigo
    secondaryColor   = lipgloss.Color("#10b981")   // Emerald green
    accentColor      = lipgloss.Color("#f59e0b")   // Warm amber
    errorColor       = lipgloss.Color("#ef4444")   // Clear red
    successColor     = lipgloss.Color("#22c55e")   // Success green
    
    // Neutral colors
    backgroundDark   = lipgloss.Color("#1f2937")   // Dark background
    backgroundLight  = lipgloss.Color("#f9fafb")   // Light background
    textPrimary      = lipgloss.Color("#f9fafb")   // Primary text (light)
    textSecondary    = lipgloss.Color("#9ca3af")   // Secondary text
    textMuted        = lipgloss.Color("#6b7280")   // Muted text
    borderColor      = lipgloss.Color("#374151")   // Border color
    
    // Message colors
    userBg           = lipgloss.Color("#3b82f6")   // User message background
    userText         = lipgloss.Color("#ffffff")   // User message text
    aiBg             = lipgloss.Color("#1f2937")   // AI message background
    aiText           = lipgloss.Color("#e5e7eb")   // AI message text
    aiAccent         = lipgloss.Color("#10b981")   // AI accent color
)

// Enhanced styles with better hierarchy and spacing
var (
    // Header styles
    headerStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(textPrimary).
        Background(primaryColor).
        Padding(0, 2).
        MarginBottom(1)
    
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(primaryColor).
        MarginBottom(1)
    
    // Message styles with better visual separation
    userStyle = lipgloss.NewStyle().
        Foreground(userText).
        Background(userBg).
        Padding(1, 2).
        MarginTop(1).
        MarginBottom(1).
        Bold(true)
    
    aiStyle = lipgloss.NewStyle().
        Foreground(aiText).
        Background(aiBg).
        Padding(1, 2).
        MarginTop(1).
        MarginBottom(1).
        BorderLeft(true).
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(aiAccent)
    
    // Status and info styles
    statusStyle = lipgloss.NewStyle().
        Foreground(textSecondary).
        Background(backgroundDark).
        Padding(0, 2).
        Bold(false)
    
    infoStyle = lipgloss.NewStyle().
        Foreground(accentColor).
        Bold(true)
    
    // Interactive elements
    buttonStyle = lipgloss.NewStyle().
        Foreground(textPrimary).
        Background(secondaryColor).
        Padding(0, 2).
        MarginLeft(1).
        MarginRight(1)
    
    buttonInactiveStyle = lipgloss.NewStyle().
        Foreground(textMuted).
        Background(borderColor).
        Padding(0, 2).
        MarginLeft(1).
        MarginRight(1)
    
    // Utility styles
    separatorStyle = lipgloss.NewStyle().
        Foreground(borderColor).
        MarginTop(1).
        MarginBottom(1)
    
    loadingStyle = lipgloss.NewStyle().
        Foreground(accentColor).
        Bold(true)
    
    errorStyle = lipgloss.NewStyle().
        Foreground(errorColor).
        Bold(true).
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(errorColor)
    
    successStyle = lipgloss.NewStyle().
        Foreground(successColor).
        Bold(true)
    
    // Help and shortcuts
    helpStyle = lipgloss.NewStyle().
        Foreground(textSecondary).
        Italic(true)
    
    shortcutStyle = lipgloss.NewStyle().
        Foreground(accentColor).
        Bold(true)
)
