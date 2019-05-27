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
		case ".yaml", ".yml":
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
// or {"a": {"b":[1,2]}} => {"a.b.0":1, "a.b.1": 2}
func flatten(data map[string]interface{}) map[string]string {
	ret := make(map[string]string)
	for k, v := range data {
		switch typed := v.(type) {
		case map[interface{}]interface{}:
			for fk, fv := range flatten(convertMap(typed)) {
				ret[fmt.Sprintf("%s.%s", k, fk)] = fv
			}
		case map[string]interface{}:
			for fk, fv := range flatten(typed) {
				ret[fmt.Sprintf("%s.%s", k, fk)] = fv
			}
		case []interface{}:
			for fk, fv := range flattenSlice(typed) {
				ret[fmt.Sprintf("%s.%s", k, fk)] = fv
			}
		default:
			ret[k] = fmt.Sprint(typed)
		}
	}
	return ret
}
func flattenSlice(data []interface{}) map[string]string {
	ret := make(map[string]string)
	for idx, v := range data {
		switch typed := v.(type) {
		case map[interface{}]interface{}:
			for fk, fv := range flatten(convertMap(typed)) {
				ret[fmt.Sprintf("%d.%s", idx, fk)] = fv
			}
		case map[string]interface{}:
			for fk, fv := range flatten(typed) {
				ret[fmt.Sprintf("%d.%s", idx, fk)] = fv
			}
		case []interface{}:
			for fk, fv := range flattenSlice(typed) {
				ret[fmt.Sprintf("%d.%s", idx, fk)] = fv
			}
		default:
			ret[fmt.Sprint(idx)] = fmt.Sprint(typed)
		}
	}
	return ret
}

func convertMap(originalMap map[interface{}]interface{}) map[string]interface{} {
	convertedMap := map[string]interface{}{}
	for key, value := range originalMap {
		convertedMap[key.(string)] = value
	}
	return convertedMap
}
