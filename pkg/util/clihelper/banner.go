package clihelper

import (
	"fmt"
	"strings"

	"github.com/loft-sh/log"
	"github.com/sirupsen/logrus"
)

// vCluster ASCII art logo - a stylized "v" shape
var vclusterLogo = []string{
	` _____`,
	`\_____ / /`,
	`  ___ / /`,
	` \   \ \_/`,
	`  \   \`,
	`   \  /`,
	`    \/`,
}

// PrintConnectSuccessBanner prints the ASCII art banner with vcluster connection success message
func PrintConnectSuccessBanner(vclusterName, contextName string, log log.Logger) {
	// Build the text lines that appear next to the ASCII art
	textLines := []string{
		"",
		fmt.Sprintf("Virtual Cluster %s is running and was set up as the default context in your", vclusterName),
		"kubeconfig (~/.kube/config)",
		"",
		"To get started:",
		"- Use `kubectl` (or any other client) to access the Virtual Cluster",
		"- Use `vcluster disconnect` to return to your previous Kubernetes context",
	}

	// Print empty line before banner
	log.WriteString(logrus.InfoLevel, "\n")

	// Print the banner with logo and text side by side
	maxLogoWidth := 0
	for _, line := range vclusterLogo {
		if len(line) > maxLogoWidth {
			maxLogoWidth = len(line)
		}
	}

	// Calculate padding for alignment
	padding := maxLogoWidth + 3 // Add some space between logo and text

	// Print each line
	maxLines := len(vclusterLogo)
	if len(textLines) > maxLines {
		maxLines = len(textLines)
	}

	for i := 0; i < maxLines; i++ {
		var logoLine string
		var textLine string

		if i < len(vclusterLogo) {
			logoLine = vclusterLogo[i]
		}
		if i < len(textLines) {
			textLine = textLines[i]
		}

		// Pad the logo line to align text
		paddedLogo := logoLine + strings.Repeat(" ", padding-len(logoLine))
		log.WriteString(logrus.InfoLevel, paddedLogo+textLine+"\n")
	}
}

// PrintConnectSuccessBannerPortForwarding prints the banner for port-forwarding mode
func PrintConnectSuccessBannerPortForwarding(vclusterName, contextName string, log log.Logger) {
	// Build the text lines for port-forwarding mode
	textLines := []string{
		"",
		fmt.Sprintf("Virtual Cluster %s is running and was set up as the default context in your", vclusterName),
		"kubeconfig (~/.kube/config)",
		"",
		"To get started:",
		"- Use `kubectl` (or any other client) to access the Virtual Cluster",
		"- Use CTRL+C to return to your previous Kubernetes context",
	}

	// Print empty line before banner
	log.WriteString(logrus.InfoLevel, "\n")

	// Print the banner with logo and text side by side
	maxLogoWidth := 0
	for _, line := range vclusterLogo {
		if len(line) > maxLogoWidth {
			maxLogoWidth = len(line)
		}
	}

	// Calculate padding for alignment
	padding := maxLogoWidth + 3

	// Print each line
	maxLines := len(vclusterLogo)
	if len(textLines) > maxLines {
		maxLines = len(textLines)
	}

	for i := 0; i < maxLines; i++ {
		var logoLine string
		var textLine string

		if i < len(vclusterLogo) {
			logoLine = vclusterLogo[i]
		}
		if i < len(textLines) {
			textLine = textLines[i]
		}

		// Pad the logo line to align text
		paddedLogo := logoLine + strings.Repeat(" ", padding-len(logoLine))
		log.WriteString(logrus.InfoLevel, paddedLogo+textLine+"\n")
	}
}
