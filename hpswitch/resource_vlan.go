package hpswitch

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVlanCreate,
		ReadContext:   resourceVlanRead,
		UpdateContext: resourceVlanUpdate,
		DeleteContext: resourceVlanDelete,
		Schema: map[string]*schema.Schema{
			"vlan": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "VLAN",
			},
			"tagged_ports": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceVlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vlanId := d.Get("vlan").(string)
	description := d.Get("description").(string)
	taggedPorts := d.Get("tagged_ports").([]interface{})

	fmt.Println(vlanId)

	c := m.(client)

	id := c.setVlan(vlanId, description, taggedPorts)

	d.SetId(id)

	resourceVlanRead(ctx, d, m)

	return diags
}

func resourceVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	c := m.(client)

	vlanTag := d.Get("vlan").(string)

	vlanID, vlanDescription, taggedPortsMaps := c.readVlan(vlanTag)

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

func resourceVlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Warning Message Summary",
		Detail:   "This is the detailed warning message from resourceVlanUpdate",
	})

	vlanID := d.Get("vlan").(string)
	description := d.Get("description").(string)
	taggedPorts := d.Get("tagged_ports").([]interface{})

	fmt.Println(vlanID)

	c := m.(client)

	id := c.setVlan(vlanID, description, taggedPorts)

	d.SetId(id)

	resourceVlanRead(ctx, d, m)

	return diags
}

func resourceVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
