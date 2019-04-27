package sops

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const configTestDataSourceSopsFile_basic = `
data "sops_file" "test_basic" {
  source_file = "%s/test-fixtures/basic.yaml"
}`

func TestDataSourceSopsFile_basic(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_basic, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "data.bool", "true"),
				),
			},
		},
	})
}

const configTestDataSourceSopsFile_nested = `
data "sops_file" "test_nested" {
  source_file = "%s/test-fixtures/nested.yaml"
}`

func TestDataSourceSopsFile_nested(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_nested, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_nested", "data.db.user", "foo"),
					resource.TestCheckResourceAttr("data.sops_file.test_nested", "data.db.password", "bar"),
				),
			},
		},
	})
}

const configTestDataSourceSopsFile_raw = `
data "sops_file" "test_raw" {
  source_file = "%s/test-fixtures/raw.txt"
  input_type = "raw"
}`

func TestDataSourceSopsFile_raw(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_raw, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_raw", "raw", "Hello raw world!"),
				),
			},
		},
	})
}

const configTestDataSourceSopsFile_simplelist = `
data "sops_file" "test_list" {
  source_file = "%s/test-fixtures/simple-list.yaml"
}`

func TestDataSourceSopsFile_simplelist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_simplelist, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.0", "val1"),
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.1", "val2"),
				),
			},
		},
	})
}

const configTestDataSourceSopsFile_complexlist = `
data "sops_file" "test_list" {
  source_file = "%s/test-fixtures/complex-list.yaml"
}`

func TestDataSourceSopsFile_complexlist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_complexlist, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.0.name", "foo"),
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.0.index", "0"),
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.1.name", "bar"),
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.1.index", "1"),
				),
			},
		},
	})
}

func TestFlattening(t *testing.T) {
	tc := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]string
	}{
		{
			name: "all data types become strings",
			input: map[string]interface{}{
				"a_string":   "foo",
				"an_integer": 12,
				"a_bool":     true,
				"a_float":    1.1,
			},
			expected: map[string]string{
				"a_string":   "foo",
				"an_integer": "12",
				"a_bool":     "true",
				"a_float":    "1.1",
			},
		},
		{
			name: "dicts are unnested",
			input: map[string]interface{}{
				"a_dict": map[string]interface{}{"foo": "bar"},
			},
			expected: map[string]string{
				"a_dict.foo": "bar",
			},
		},
		{
			name: "lists are unpacked with index keys",
			input: map[string]interface{}{
				"a_list": []interface{}{1, 2},
			},
			expected: map[string]string{
				"a_list.0": "1",
				"a_list.1": "2",
			},
		},
		{
			name: "deep nesting",
			/*
				This test corresponds to this yaml structure:
				foo:
				- a: 1
				  b:
				    c:
				    - d: 2
			*/
			input: map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"a": 1,
						"b": map[string]interface{}{
							"c": []interface{}{
								map[string]interface{}{"d": 2},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"foo.0.a":       "1",
				"foo.0.b.c.0.d": "2",
			},
		},
	}
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			output := flatten(c.input)
			if !reflect.DeepEqual(c.expected, output) {
				t.Errorf("Unexpected flattening output, expected %v, got %v", c.expected, output)
			}
		})
	}
}
