package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/modfin/terraform-provider-ceph/ceph"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ceph.Provider,
	})
}
