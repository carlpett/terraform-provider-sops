package sops

import (
	"fmt"
	"io/ioutil"
	"path"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	if err := validateInputType(format); err != nil {
		return err
	}

	if os.Getenv("SOPS_SKIP_KMS_ASSUME_ROLE") == "1" {
		regex := regexp.MustCompile("(?m)[\r\n]+^[[:space:]]+\"?role\"?:[[:space:]]+\"?arn:aws:iam::.+\"?,?")
		withoutRoles := regex.ReplaceAll(content, []byte(""))
		return readData(withoutRoles, format, d)
	} else {
		return readData(content, format, d)
	}
}
