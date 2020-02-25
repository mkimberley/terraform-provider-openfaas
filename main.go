package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/mkimberley/terraform-provider-openfaas/openfaas"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: openfaas.Provider})
}
