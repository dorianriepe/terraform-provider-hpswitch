package hpswitch

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPSWITCH_HOST", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPSWITCH_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HPSWITCH_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hpswitch_vlan": resourceVlan(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hpswitch_vlan": dataSourceVlan(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// Configure the provider client -
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	host := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var diags diag.Diagnostics

	if (username != "") && (password != "") {
		c := client{host, username, password}
		return c, diags
	}

	c := client{"localhost", "admin", "admin"}
	return c, diags
}
