// Package cmd provides functionally common to various work CLIs
package cmd

import (
	"fmt"
	"net/url"
	"os"
	"os/user"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kubework/work"
)

// NewVersionCmd returns a new `version` command to be used as a sub-command to root
func NewVersionCmd(cliName string) *cobra.Command {
	var short bool
	versionCmd := cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print version information"),
		Run: func(cmd *cobra.Command, args []string) {
			version := work.GetVersion()
			fmt.Printf("%s: %s\n", cliName, version)
			if short {
				return
			}
			fmt.Printf("  BuildDate: %s\n", version.BuildDate)
			fmt.Printf("  GitCommit: %s\n", version.GitCommit)
			fmt.Printf("  GitTreeState: %s\n", version.GitTreeState)
			if version.GitTag != "" {
				fmt.Printf("  GitTag: %s\n", version.GitTag)
			}
			fmt.Printf("  GoVersion: %s\n", version.GoVersion)
			fmt.Printf("  Compiler: %s\n", version.Compiler)
			fmt.Printf("  Platform: %s\n", version.Platform)
		},
	}
	versionCmd.Flags().BoolVar(&short, "short", false, "print just the version number")
	return &versionCmd
}

// MustIsDir returns whether or not the given filePath is a directory. Exits if path does not exist
func MustIsDir(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return fileInfo.IsDir()
}

// MustHomeDir returns the home directory of the user
func MustHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// IsURL returns whether or not a string is a URL
func IsURL(u string) bool {
	var parsedURL *url.URL
	var err error

	parsedURL, err = url.ParseRequestURI(u)
	if err == nil {
		if parsedURL != nil && parsedURL.Host != "" {
			return true
		}
	}
	return false
}

// ParseLabels turns a string representation of a label set into a map[string]string
func ParseLabels(labelSpec interface{}) (map[string]string, error) {
	labelString, isString := labelSpec.(string)
	if !isString {
		return nil, fmt.Errorf("expected string, found %v", labelSpec)
	}
	if len(labelString) == 0 {
		return nil, fmt.Errorf("no label spec passed")
	}
	labels := map[string]string{}
	labelSpecs := strings.Split(labelString, ",")
	for ix := range labelSpecs {
		labelSpec := strings.Split(labelSpecs[ix], "=")
		if len(labelSpec) != 2 {
			return nil, fmt.Errorf("unexpected label spec: %s", labelSpecs[ix])
		}
		if len(labelSpec[0]) == 0 {
			return nil, fmt.Errorf("unexpected empty label key")
		}
		labels[labelSpec[0]] = labelSpec[1]
	}
	return labels, nil
}
