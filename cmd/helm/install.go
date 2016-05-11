package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kubernetes/helm/pkg/helm"
	"github.com/kubernetes/helm/pkg/proto/hapi/release"
	"github.com/kubernetes/helm/pkg/timeconv"
)

const installDesc = `
This command installs a chart archive.

The install argument must be either a relative
path to a chart directory or the name of a
chart in the current working directory.
`

const (
	hostEnvVar  = "TILLER_HOST"
	defaultHost = ":44134"
)

// install flags & args
var (
	// installArg is the name or relative path of the chart to install
	installArg string
	// tillerHost overrides TILLER_HOST envVar
	tillerHost string
	// installDryRun performs a dry-run install
	installDryRun bool
)

var installCmd = &cobra.Command{
	Use:   "install [CHART]",
	Short: "install a chart archive.",
	Long:  installDesc,
	RunE:  runInstall,
}

func runInstall(cmd *cobra.Command, args []string) error {
	if err := checkArgsLength(1, len(args), "chart name"); err != nil {
		return err
	}
	installArg = args[0]
	setupInstallEnv()

	res, err := helm.InstallRelease(installArg, installDryRun)
	if err != nil {
		return prettyError(err)
	}

	printRelease(res.GetRelease())

	return nil
}

func printRelease(rel *release.Release) {
	if rel == nil {
		return
	}
	if flagVerbose {
		fmt.Printf("NAME:   %s\n", rel.Name)
		fmt.Printf("INFO:   %s %s\n", timeconv.String(rel.Info.LastDeployed), rel.Info.Status)
		fmt.Printf("CHART:  %s %s\n", rel.Chart.Metadata.Name, rel.Chart.Metadata.Version)
		fmt.Printf("MANIFEST: %s\n", rel.Manifest)
	} else {
		fmt.Println(rel.Name)
	}
}

func setupInstallEnv() {
	// The 'host' flag takes precendence uber alles.
	// If set, proceed with install using provided host
	// address. If unset, the 'TILLER_HOST' environment
	// variable (if set) is used, otherwise defaults to ":44134".
	if tillerHost == "" {
		tillerHost = os.Getenv(hostEnvVar)
		if tillerHost == "" {
			tillerHost = defaultHost
		}
	}

	helm.Config.ServAddr = tillerHost
}

func init() {
	installCmd.Flags().StringVar(&tillerHost, "host", "", "address of tiller server (default \":44134\")")
	installCmd.Flags().BoolVar(&installDryRun, "dry-run", false, "simulate an install")

	RootCommand.AddCommand(installCmd)
}
