package container

import (
	"sort"
	"strings"

	engineapi "github.com/docker/docker/client"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/agent/secrets"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"golang.org/x/net/context"
)

type executor struct {
	client  engineapi.APIClient
	secrets exec.SecretsManager
}

// NewExecutor returns an executor from the docker client.
func NewExecutor(client engineapi.APIClient) exec.Executor {
	return &executor{
		client:  client,
		secrets: secrets.NewManager(),
	}
}

// Describe returns the underlying node description from the docker client.
func (e *executor) Describe(ctx context.Context) (*api.NodeDescription, error) {
	info, err := e.client.Info(ctx)
	if err != nil {
		return nil, err
	}

	plugins := map[api.PluginDescription]struct{}{}
	addPlugins := func(typ string, names []string) {
		for _, name := range names {
			plugins[api.PluginDescription{
				Type: typ,
				Name: name,
			}] = struct{}{}
		}
	}

	// add v1 plugins to 'plugins'
	addPlugins("Volume", info.Plugins.Volume)
	// Add builtin driver "overlay" (the only builtin multi-host driver) to
	// the plugin list by default.
	addPlugins("Network", append([]string{"overlay"}, info.Plugins.Network...))
	addPlugins("Authorization", info.Plugins.Authorization)

	// retrieve v2 plugins
	v2plugins, err := e.client.PluginList(ctx)
	if err != nil {
		log.L.WithError(err).Warning("PluginList operation failed")
	} else {
		// add v2 plugins to 'plugins'
		for _, plgn := range v2plugins {
			for _, typ := range plgn.Config.Interface.Types {
				if typ.Prefix == "docker" && plgn.Enabled {
					plgnTyp := typ.Capability
					if typ.Capability == "volumedriver" {
						plgnTyp = "Volume"
					} else if typ.Capability == "networkdriver" {
						plgnTyp = "Network"
					}
					plgnName := plgn.Name
					if plgn.Tag != "" {
						plgnName += ":" + plgn.Tag
					}
					plugins[api.PluginDescription{
						Type: plgnTyp,
						Name: plgnName,
					}] = struct{}{}
				}
			}
		}
	}

	pluginFields := make([]api.PluginDescription, 0, len(plugins))
	for k := range plugins {
		pluginFields = append(pluginFields, k)
	}
	sort.Sort(sortedPlugins(pluginFields))

	// parse []string labels into a map[string]string
	labels := map[string]string{}
	for _, l := range info.Labels {
		stringSlice := strings.SplitN(l, "=", 2)
		// this will take the last value in the list for a given key
		// ideally, one shouldn't assign multiple values to the same key
		if len(stringSlice) > 1 {
			labels[stringSlice[0]] = stringSlice[1]
		}
	}

	description := &api.NodeDescription{
		Hostname: info.Name,
		Platform: &api.Platform{
			Architecture: info.Architecture,
			OS:           info.OSType,
		},
		Engine: &api.EngineDescription{
			EngineVersion: info.ServerVersion,
			Labels:        labels,
			Plugins:       pluginFields,
		},
		Resources: &api.Resources{
			NanoCPUs:    int64(info.NCPU) * 1e9,
			MemoryBytes: info.MemTotal,
		},
	}

	return description, nil
}

func (e *executor) Configure(ctx context.Context, node *api.Node) error {
	return nil
}

// Controller returns a docker container controller.
func (e *executor) Controller(t *api.Task) (exec.Controller, error) {
	ctlr, err := newController(e.client, t, secrets.Restrict(e.secrets, t))
	if err != nil {
		return nil, err
	}

	return ctlr, nil
}

func (e *executor) SetNetworkBootstrapKeys([]*api.EncryptionKey) error {
	return nil
}

func (e *executor) Secrets() exec.SecretsManager {
	return e.secrets
}

type sortedPlugins []api.PluginDescription

func (sp sortedPlugins) Len() int { return len(sp) }

func (sp sortedPlugins) Swap(i, j int) { sp[i], sp[j] = sp[j], sp[i] }

func (sp sortedPlugins) Less(i, j int) bool {
	if sp[i].Type != sp[j].Type {
		return sp[i].Type < sp[j].Type
	}
	return sp[i].Name < sp[j].Name
}
