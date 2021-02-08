package powerbi

import (
	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourcePBIX represents a Power BI PBIX file
func DataCapacity() *schema.Resource {
	return &schema.Resource{
		Read: readCapacity,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the PowerBI Embedded instance to filter on",
			},
			"capacity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The capacity_id of the PowerBI Embedded instance",
			},

		},
	}
}

func readCapacity(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	capacities, err := client.GetCapacities()
	name := d.Get("name")
	if err != nil {
		panic(err)
		return err
	}

	if capacities == nil {
		d.SetId("")
		err := d.Set("capacity_id", "no capacities")
		if err != nil { return err }
	} else {
		var item *powerbiapi.GetCapacitiesResponseItem = nil
		for _, responseItem := range capacities.Value {
			if responseItem.DisplayName == name {
				item = &responseItem
				break
			}
		}
		if item == nil {
			d.SetId("")
			err := d.Set("capacity_id", "not found in list")
			if err != nil { return err }
		} else {
			err := d.Set("capacity_id", item.ID)
			if err != nil { return err }
			err = d.Set("name", name)
			if err != nil { return err }
		}
	}

	return nil
}