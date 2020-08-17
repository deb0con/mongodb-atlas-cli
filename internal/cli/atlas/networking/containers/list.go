// Copyright 2020 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package containers

import (
	"github.com/mongodb/mongocli/internal/cli"
	"github.com/mongodb/mongocli/internal/config"
	"github.com/mongodb/mongocli/internal/flag"
	"github.com/mongodb/mongocli/internal/output"
	"github.com/mongodb/mongocli/internal/store"
	"github.com/mongodb/mongocli/internal/usage"
	"github.com/spf13/cobra"
	atlas "go.mongodb.org/atlas/mongodbatlas"
)

type ListOpts struct {
	cli.GlobalOpts
	cli.ListOpts
	provider string
	store    store.ContainersLister
}

func (opts *ListOpts) initStore() error {
	var err error
	opts.store, err = store.New(config.Default())
	return err
}

var listTemplate = `ID	PROVIDER	REGION	ATLAS CIDR	PROVISIONED{{range .}}
{{.ID}}	{{.ProviderName}}	{{if .RegionName}}{{.RegionName}}{{else}}{{.Region}}{{end}}	{{.AtlasCIDRBlock}}	{{.Provisioned}}{{end}}
`

func (opts *ListOpts) Run() error {
	var r []atlas.Container
	var err error
	if opts.provider == "" {
		r, err = opts.store.AllContainers(opts.ConfigProjectID(), opts.NewListOptions())
	} else {
		listOpts := opts.newContainerListOptions()
		r, err = opts.store.ContainersByProvider(opts.ConfigProjectID(), listOpts)
	}
	if err != nil {
		return err
	}

	return output.Print(config.Default(), listTemplate, r)
}

func (opts *ListOpts) newContainerListOptions() *atlas.ContainersListOptions {
	return &atlas.ContainersListOptions{
		ListOptions: *opts.NewListOptions(),
	}
}

// mongocli atlas networking container(s) list [--projectId projectId] [--orgId orgId] [--page N] [--limit N] [--minDate minDate] [--maxDate maxDate]
func ListBuilder() *cobra.Command {
	opts := &ListOpts{}
	cmd := &cobra.Command{
		Use:     "list",
		Short:   listContainers,
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.initStore()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}
	cmd.Flags().StringVar(&opts.provider, flag.Provider, "", usage.Provider)
	cmd.Flags().IntVar(&opts.PageNum, flag.Page, 0, usage.Page)
	cmd.Flags().IntVar(&opts.ItemsPerPage, flag.Limit, 0, usage.Limit)

	cmd.Flags().StringVar(&opts.ProjectID, flag.ProjectID, "", usage.ProjectID)

	return cmd
}