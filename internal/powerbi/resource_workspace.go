package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceWorkspace represents a Power BI workspace
func ResourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: createWorkspace,
		Read:   readWorkspace,
		Update: updateWorkspace,
		Delete: deleteWorkspace,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the workspace.",
			},
			"capacity_to_use_display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display Name of the capacity to use for this workspace.",
			},
			"capacity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the capacity assigned to this workspace.",
			},
		},
	}
}

func createWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	resp, err := client.CreateGroup(powerbiapi.CreateGroupRequest{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(resp.ID)
	if d.Get("capacity_to_use_display_name") != "" {
		err := setCapacity(d, meta)
		if err != nil {
			return err
		}
	}
	return readWorkspace(d, meta)
}

func readWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	workspace, err := client.GetGroup(d.Id())
	if err != nil {
		return err
	}

	if workspace == nil {
		d.SetId("")
	} else {
		d.SetId(workspace.ID)
		d.Set("name", workspace.Name)
		capacityID := workspace.CapacityID
		d.Set("capacity_id", capacityID)
		if capacityID == "" {
			d.Set("capacity_to_use_display_name", "")
		} else {
			capacity, err := getCapacityById(client, capacityID)
			if err != nil {
				return err
			}
			d.Set("capacity_to_use_display_name", capacity.DisplayName)
		}
	}

	return nil
}

func updateWorkspace(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("capacity_to_use_display_name") {
		err := setCapacity(d, meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func setCapacity(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	workspaceId := d.Id()
	displayName := d.Get("capacity_to_use_display_name").(string)
	// setting to this capacity_id unassigns this workspace
	capacityId := "00000000-0000-0000-0000-000000000000"
	if displayName != "" {
		item, err := getCapacityByName(client, displayName)
		if err != nil {
			return err
		}
		capacityId = item.ID
	}
	err := client.GroupAssignToCapacity(workspaceId, powerbiapi.GroupAssignToCapacityRequest{
		CapacityID: capacityId,
	})
	if err != nil {
		return err
	}
	err = d.Set("capacity_to_use_display_name", displayName)
	if err != nil {
		return err
	}
	err = d.Set("capacity_id", "")
	if err != nil {
		return err
	}

	return nil
}

func getCapacityByName(client *powerbiapi.Client, displayName string) (*powerbiapi.GetCapacitiesResponseItem, error) {
	capacities, err := client.GetCapacities()
	if err != nil {
		return nil, err
	}
	var item *powerbiapi.GetCapacitiesResponseItem = nil
	for _, responseItem := range capacities.Value {
		if responseItem.DisplayName == displayName {
			item = &responseItem
			break
		}
	}
	if item == nil {
		return nil, fmt.Errorf("capacity not found for display name: %s", displayName)
	}
	return item, nil
}

func getCapacityById(client *powerbiapi.Client, ID string) (*powerbiapi.GetCapacitiesResponseItem, error) {
	capacities, err := client.GetCapacities()
	if err != nil {
		return nil, err
	}
	var item *powerbiapi.GetCapacitiesResponseItem = nil
	for _, responseItem := range capacities.Value {
		if responseItem.ID == ID {
			item = &responseItem
			break
		}
	}
	if item == nil {
		return nil, fmt.Errorf("capacity not found with ID: %s", ID)
	}
	return item, nil
}


func deleteWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	return client.DeleteGroup(d.Id())
}

func assignToCapacity(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	capacityID := d.Get("capacity_id").(string)
	if capacityID != "00000000-0000-0000-0000-000000000000" {
		var capacityObjFound bool

		capacityList, err := client.GetCapacities()
		if err != nil {
			return err
		}

		if len(capacityList.Value) >= 1 {
			for _, capacityObj := range capacityList.Value {
				if capacityObj.ID == capacityID {
					capacityObjFound = true
				}
			}
		}
		if capacityObjFound != true {
			return fmt.Errorf("Capacity id %s not found or logged-in user doesn't have capacity admin rights", capacityID)
		}
	}

	err := client.GroupAssignToCapacity(d.Id(), powerbiapi.GroupAssignToCapacityRequest{
		CapacityID: capacityID,
	})
	if err != nil {
		return err
	}

	return nil
}
