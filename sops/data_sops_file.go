package sops

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/hashicorp/terraform/helper/schema"
	"go.mozilla.org/sops/decrypt"
	"gopkg.in/yaml.v2"
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
			"raw": &schema.Schema{
				Type:     schema.TypeString,
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
			return fmt.Errorf("Don't know how to decode file with extension %s, set input_type to json, yaml or raw as appropriate", ext)
		}
	}

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

// flatten flattens the nested struct.
//
// All keys will be joined by dot
// e.g. {"a": {"b":"c"}} => {"a.b":"c"}
func flatten(data map[string]interface{}) map[string]string {
	flattened := map[string]string{}
	for key, value := range data {
		switch t := value.(type) {
		case string, int:
			flattened[key] = fmt.Sprint(t)
		case map[interface{}]interface{}:
			f := flatten(convertMap(t))
			for k, v := range f {
				// Join all keys with dot
				flattened[fmt.Sprintf("%s.%s", key, k)] = v
			}
		default:
			fmt.Printf("unexpected type %T", t)
		}
	}
	return flattened
}

func convertMap(originalMap map[interface{}]interface{}) map[string]interface{} {
	convertedMap := map[string]interface{}{}
	for key, value := range originalMap {
		convertedMap[key.(string)] = value
	}
	return convertedMap
}
