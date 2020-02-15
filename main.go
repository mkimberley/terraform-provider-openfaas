package main

import (
	"github.com/mkimberley/terraform-provider-openfaas/openfaas"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: openfaas.Provider})
}
