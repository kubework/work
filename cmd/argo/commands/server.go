package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/kubework/pkg/stats"
	log "github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/kubework/work/cmd/work/commands/client"
	wfclientset "github.com/kubework/work/pkg/client/clientset/versioned"
	"github.com/kubework/work/server/apiserver"
	"github.com/kubework/work/util/help"
)

func NewServerCommand() *cobra.Command {
	var (
		authMode          string
		configMap         string
		port              int
		baseHRef          string
		namespaced        bool   // --namespaced
		managedNamespace  string // --managed-namespace
		enableOpenBrowser bool
	)

	var command = cobra.Command{
		Use:   "server",
		Short: "Start the Work Server",
		Example: fmt.Sprintf(`
See %s`, help.WorkSever),
		RunE: func(c *cobra.Command, args []string) error {
			stats.RegisterStackDumper()
			stats.StartStatsTicker(5 * time.Minute)

			config, err := client.Config.ClientConfig()
			if err != nil {
				return err
			}
			config.Burst = 30
			config.QPS = 20.0

			namespace, _, err := client.Config.Namespace()
			if err != nil {
				return err
			}

			kubeConfig := kubernetes.NewForConfigOrDie(config)
			wflientset := wfclientset.NewForConfigOrDie(config)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if !namespaced && managedNamespace != "" {
				log.Warn("ignoring --managed-namespace because --namespaced is false")
				managedNamespace = ""
			}
			if namespaced && managedNamespace == "" {
				managedNamespace = namespace
			}

			log.WithFields(log.Fields{
				"authMode":         authMode,
				"namespace":        namespace,
				"managedNamespace": managedNamespace,
				"baseHRef":         baseHRef}).
				Info()

			opts := apiserver.WorkServerOpts{
				BaseHRef:         baseHRef,
				Namespace:        namespace,
				WfClientSet:      wflientset,
				KubeClientset:    kubeConfig,
				RestConfig:       config,
				AuthMode:         authMode,
				ManagedNamespace: managedNamespace,
				ConfigName:       configMap,
			}
			err = opts.ValidateOpts()
			if err != nil {
				return err
			}
			browserOpenFunc := func(url string) {}
			if enableOpenBrowser {
				browserOpenFunc = func(url string) {
					log.Infof("Work UI is available at %s", url)
					err := open.Run(url)
					if err != nil {
						log.Warnf("Unable to open the browser. %v", err)
					}
				}
			}
			apiserver.NewWorkServer(opts).Run(ctx, port, browserOpenFunc)
			return nil
		},
	}

	command.Flags().IntVarP(&port, "port", "p", 2746, "Port to listen on")
	defaultBaseHRef := os.Getenv("BASE_HREF")
	if defaultBaseHRef == "" {
		defaultBaseHRef = "/"
	}
	command.Flags().StringVar(&baseHRef, "basehref", defaultBaseHRef, "Value for base href in index.html. Used if the server is running behind reverse proxy under subpath different from /. Defaults to the environment variable BASE_HREF.")
	command.Flags().StringVar(&authMode, "auth-mode", "server", "API server authentication mode. One of: client|server|hybrid")
	command.Flags().StringVar(&configMap, "configmap", "workflow-controller-configmap", "Name of K8s configmap to retrieve workflow controller configuration")
	command.Flags().BoolVar(&namespaced, "namespaced", false, "run as namespaced mode")
	command.Flags().StringVar(&managedNamespace, "managed-namespace", "", "namespace that watches, default to the installation namespace")
	command.Flags().BoolVarP(&enableOpenBrowser, "browser", "b", false, "enable automatic launching of the browser [local mode]")
	return &command
}
