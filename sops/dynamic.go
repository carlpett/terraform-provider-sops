package sops

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"gopkg.in/yaml.v3"
)

// dynamicToGo converts an arbitrary Terraform value into plain Go values for marshalling.
func dynamicToGo(ctx context.Context, dyn types.Dynamic) (any, error) {
	tfVal, err := dyn.ToTerraformValue(ctx)
	if err != nil {
		return nil, err
	}
	return tftypesValueToGo(tfVal)
}

func tftypesValueToGo(v tftypes.Value) (any, error) {
	if !v.IsKnown() {
		return nil, fmt.Errorf("value is not known")
	}
	if v.IsNull() {
		return nil, nil
	}

	typ := v.Type()
	switch {
	case typ.Is(tftypes.String):
		var s string
		if err := v.As(&s); err != nil {
			return nil, err
		}
		return s, nil
	case typ.Is(tftypes.Bool):
		var b bool
		if err := v.As(&b); err != nil {
			return nil, err
		}
		return b, nil
	case typ.Is(tftypes.Number):
		var n *big.Float
		if err := v.As(&n); err != nil {
			return nil, err
		}
		if n.IsInt() {
			i, _ := n.Int64()
			return i, nil
		}
		f, _ := n.Float64()
		return f, nil
	}

	switch typ.(type) {
	case tftypes.List, tftypes.Set, tftypes.Tuple:
		var elems []tftypes.Value
		if err := v.As(&elems); err != nil {
			return nil, err
		}
		result := make([]any, len(elems))
		for i, e := range elems {
			converted, err := tftypesValueToGo(e)
			if err != nil {
				return nil, err
			}
			result[i] = converted
		}
		return result, nil
	case tftypes.Map, tftypes.Object:
		var elems map[string]tftypes.Value
		if err := v.As(&elems); err != nil {
			return nil, err
		}
		result := make(map[string]any, len(elems))
		for k, e := range elems {
			converted, err := tftypesValueToGo(e)
			if err != nil {
				return nil, err
			}
			result[k] = converted
		}
		return result, nil
	}

	return nil, fmt.Errorf("unsupported type %s", typ)
}

// marshalDynamicValue marshals value (from dynamicToGo) as JSON or YAML.
func marshalDynamicValue(format string, value any) ([]byte, error) {
	switch format {
	case "json":
		return json.Marshal(value)
	case "yaml":
		return yaml.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported format %q", format)
	}
}
