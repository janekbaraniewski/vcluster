package banner

import (
	"fmt"

	"github.com/loft-sh/log"
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
)

const asciiArt = `
             __           __           
 _   _______/ /_  _______/ /____  _____
| | / / ___/ / / / / ___/ __/ _ \/ ___/
| |/ / /__/ / /_/ (__  ) /_/  __/ /    
|___/\___/_/\__,_/____/\__/\___/_/     
`

func PrintSuccessMessageVClusterConnect(vClusterName, kubeConfigPath string, logger log.Logger) {
	logger.WriteString(logrus.InfoLevel, ansi.Color(asciiArt, "green+b"))
	logger.WriteString(logrus.InfoLevel, fmt.Sprintf("\n  Virtual Cluster %s is running and was set up as the default context\n", ansi.Color(vClusterName, "white+b")))
	logger.WriteString(logrus.InfoLevel, fmt.Sprintf("  in your kubeconfig (%s)\n\n", kubeConfigPath))
	logger.WriteString(logrus.InfoLevel, "  To get started:\n")
}
