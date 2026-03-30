package provider_test

import (
	"github.com/datafy-io/terraform-provider-datafy/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"datafy": providerserver.NewProtocol6WithError(provider.New("test")()),
}
