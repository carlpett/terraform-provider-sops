package sops

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/hashicorp/terraform/helper/schema"
	"go.mozilla.org/sops/decrypt"
)

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFileRead,

		Schema: map[string]*schema.Schema{
			"input_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_file": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"data": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceFileRead(d *schema.ResourceData, meta interface{}) error {
	sourceFile := d.Get("source_file").(string)
	content, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	var format string
	if input_type := d.Get("input_type").(string); input_type != "" {
		format = input_type
	} else {
		switch ext := path.Ext(sourceFile); ext {
		case ".json":
			format = "json"
		case ".yaml":
			format = "yaml"
		default:
			return fmt.Errorf("Don't know how to decode file with extension %s", ext)
		}
	}

	cleartext, err := decrypt.Data(content, format)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(cleartext, &data)
	if err != nil {
		return err
	}

	err = d.Set("data", data)
	if err != nil {
		return err
	}

	d.SetId("-")
	return nil
}
