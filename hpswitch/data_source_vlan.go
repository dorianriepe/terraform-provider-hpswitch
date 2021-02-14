package hpswitch

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVlan() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVlanRead,
		Schema: map[string]*schema.Schema{
			"vlan": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tagged_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	c := m.(client)

	vlanTag := d.Get("vlan").(string)

	vlanID, vlanDescription, taggedPortsMaps := c.readVlan(vlanTag)

	if vlanID == "" && vlanDescription == "" && taggedPortsMaps == nil {
		d.SetId("")
		return diags
	}

	if err := d.Set("vlan", vlanID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", vlanDescription); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tagged_ports", taggedPortsMaps); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
