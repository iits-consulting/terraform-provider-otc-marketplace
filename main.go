package main

import (
	"context"
	"log"
	"terraform-provider-otc-marketplace/internal/provider_marketplace"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/iits-consulting/otc-marketplace",
	}

	err := providerserver.Serve(context.Background(), provider_marketplace.New(), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
