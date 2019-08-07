package sops

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"go.mozilla.org/sops/decrypt"
	"gopkg.in/yaml.v2"
)

func dataSourceExternal() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceExternalRead,

		Schema: map[string]*schema.Schema{
			"input_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"data": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"raw": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceExternalRead(d *schema.ResourceData, meta interface{}) error {
	source := d.Get("source").(string)
	content, err := ioutil.ReadAll(strings.NewReader(source))
	if err != nil {
		return err
	}

	format := d.Get("input_type").(string)
	return readData(content, format, d)
}

func readData(content []byte, format string, d *schema.ResourceData) error {
	cleartext, err := decrypt.Data(content, format)
	if err != nil {
		return err
	}

	if format != "raw" {
		var data map[string]interface{}
		switch format {
		case "json":
			err = json.Unmarshal(cleartext, &data)
		case "yaml":
			err = yaml.Unmarshal(cleartext, &data)
		default:
			return fmt.Errorf("Don't know how to unmarshal format %s", format)
		}
		if err != nil {
			return err
		}

		err = d.Set("data", flatten(data))
		if err != nil {
			return err
		}
	} else {
		err = d.Set("raw", string(cleartext))
		if err != nil {
			return err
		}
	}

	d.SetId("-")
	return nil
}
