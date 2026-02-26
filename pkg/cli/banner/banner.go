package banner

import (
	"fmt"

	"github.com/loft-sh/log"
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
)

func PrintSuccessMessageVClusterConnect(vClusterName, kubeConfigPath string, logger log.Logger) {
	logger.WriteString(logrus.InfoLevel, fmt.Sprintf(`
%s     Virtual Cluster %s is running and was set up as the default context in your
%s      kubeconfig (%s)
%s
%s        To get started:
%s         - Use %s (or any other client) to access the Virtual Cluster
               - Use %s to return to your previous Kubernetes context
`,
		ansi.Color(`\_____/ /`, "green+b"),
		ansi.Color(vClusterName, "white+b"),
		ansi.Color(`   ___/ /`, "green+b"),
		kubeConfigPath,
		ansi.Color(`  \  \ \_/`, "green+b"),
		ansi.Color(`   \  \`, "green+b"),
		ansi.Color(`    \/`, "green+b"),
		ansi.Color("`kubectl`", "white+b"),
		ansi.Color("`vcluster disconnect`", "white+b"),
	))
}
