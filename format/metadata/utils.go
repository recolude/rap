package metadata

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func getProperties(data interface{}) (Property, error) {
	switch d := data.(type) {
	case map[string]interface{}:
		x, foundX := d["x"]
		y, foundY := d["y"]
		z, foundZ := d["z"]
		if len(d) == 2 && foundX && foundY {
			return NewVector2Property(x.(float64), y.(float64)), nil
		}
		if len(d) == 3 && foundX && foundY && foundZ {
			return NewVector3Property(x.(float64), y.(float64), z.(float64)), nil
		}
		nestedProps := make(map[string]Property)
		for key, val := range d {
			np, err := getProperties(val)
			if err != nil {
				return nil, err
			}
			nestedProps[key] = np
		}
		return NewMetadataProperty(NewBlock(nestedProps)), nil
	case []interface{}:
		if _, ok := d[0].(bool); ok {
			bools := make([]bool, 0, len(d))
			for _, item := range d {
				b, _ := item.(bool)
				bools = append(bools, b)
			}
			return NewBoolArrayProperty(bools), nil
		}
		props := make([]Property, 0, len(d))
		for _, item := range d {
			propItem, err := getProperties(item)
			if err != nil {
				return nil, err
			}
			props = append(props, propItem)
		}
		return ArrayProperty{props: props, originalBaseCode: props[0].Code()}, nil
	case string:
		if strings.HasPrefix(d, HEX_PREFIX) && len(d) == 4 {
			var p ByteProperty
			if err := json.Unmarshal([]byte(fmt.Sprintf(`"%v"`, d)), &p); err != nil {
				return nil, err
			}
			return p, nil
		}
		if strings.HasPrefix(d, HEX_PREFIX) {
			b, err := hex.DecodeString(strings.TrimPrefix(d, HEX_PREFIX))
			if err != nil {
				return nil, err
			}
			return NewBinaryArrayProperty(b), nil
		}
		if _, err := time.Parse(time.RFC3339Nano, d); err == nil {
			var p TimeProperty
			if err := json.Unmarshal([]byte(fmt.Sprintf(`"%v"`, d)), &p); err != nil {
				return nil, err
			}
			return p, nil
		}
		var p StringProperty
		if err := json.Unmarshal([]byte(fmt.Sprintf(`"%v"`, d)), &p); err != nil {
			return nil, err
		}
		return p, nil
	case float64:
		if d == float64(int64(d)) {
			var p Int32Property
			if err := json.Unmarshal([]byte(fmt.Sprintf(`%v`, d)), &p); err != nil {
				return nil, err
			}
			return p, nil
		}
		var p Float32Property
		if err := json.Unmarshal([]byte(fmt.Sprintf(`%v`, d)), &p); err != nil {
			return nil, err
		}
		return p, nil
	case bool:
		var p BoolProperty
		if err := json.Unmarshal([]byte(fmt.Sprintf(`%v`, d)), &p); err != nil {
			return nil, err
		}
		return p, nil
	}
	return nil, nil
}
