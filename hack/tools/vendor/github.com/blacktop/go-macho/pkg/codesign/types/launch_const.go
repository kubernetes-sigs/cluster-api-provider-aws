package types

import (
	"encoding/asn1"
	"fmt"
)

// LaunchContraints is the ASN.1 DER structure for launch constraints
type LaunchContraints struct {
	Count        int64          `json:"appl"`
	CCAT         int64          `json:"ccat"`
	COMP         int64          `json:"comp"`
	Requirements map[string]any `json:"reqs"`
	Version      int64          `json:"vers"`
}

const typeAppl = "application,tag:16"
const typeCont = "tag:16"

type contraint struct {
	Raw asn1.RawContent
	Key string `asn1:"utf8"`
	Val asn1.RawValue
}

type launchContraints struct {
	Raw    asn1.RawContent
	Count  int64
	Fields []contraint `asn1:"tag:16"`
}

// Entitlement https://developer.apple.com/documentation/security/defining_launch_environment_and_library_constraints#4180790
type Entitlement struct {
	Raw       asn1.RawContent
	Operation int64
	Value     asn1.RawValue
}

func checkInteger(bytes []byte) error {
	if len(bytes) == 0 {
		return fmt.Errorf("empty integer")
	}
	if len(bytes) == 1 {
		return nil
	}
	if (bytes[0] == 0 && bytes[1]&0x80 == 0) || (bytes[0] == 0xff && bytes[1]&0x80 == 0x80) {
		return fmt.Errorf("integer not minimally-encoded")
	}
	return nil
}

func parseInt64(bytes []byte) (ret int64, err error) {
	err = checkInteger(bytes)
	if err != nil {
		return
	}
	if len(bytes) > 8 {
		// We'll overflow an int64 in this case.
		err = fmt.Errorf("integer too large")
		return
	}
	for bytesRead := 0; bytesRead < len(bytes); bytesRead++ {
		ret <<= 8
		ret |= int64(bytes[bytesRead])
	}

	// Shift up and down in order to sign extend the result.
	ret <<= 64 - uint8(len(bytes))*8
	ret >>= 64 - uint8(len(bytes))*8
	return
}

func parseBool(bytes []byte) (ret bool, err error) {
	if len(bytes) != 1 {
		err = fmt.Errorf("invalid boolean")
		return
	}
	// DER demands that "If the encoding represents the boolean value TRUE,
	// its single contents octet shall have all eight bits set to one."
	// Thus only 0 and 255 are valid encoded values.
	switch bytes[0] {
	case 0:
		ret = false
	case 0xff:
		ret = true
	default:
		err = fmt.Errorf("invalid boolean")
	}
	return
}

func getValue(val asn1.RawValue) (any, error) {
	switch val.Tag {
	case asn1.TagBoolean:
		return parseBool(val.Bytes)
	case asn1.TagInteger:
		return parseInt64(val.Bytes)
	case asn1.TagUTF8String:
		return string(val.Bytes), nil
	default:
		return nil, fmt.Errorf("getValue: unsupported asn1 raw value tag %d (notify author)", val.Tag)
	}
}

func parseReqs(data []byte) (req map[string]any, err error) {
	var prop contraint

	req = make(map[string]any)

	for len(data) > 0 {
		var peek asn1.RawValue
		if _, err = asn1.Unmarshal(data, &peek); err != nil {
			return nil, fmt.Errorf("failed to ASN.1 parse launch contraint property: %v", err)
		}

		data, err = asn1.Unmarshal(data, &prop)
		if err != nil {
			return nil, fmt.Errorf("failed to ASN.1 parse launch contraint property: %v", err)
		}

		if prop.Val.IsCompound {
			if prop.Val.Class == asn1.ClassContextSpecific {
				req[prop.Key], err = parseReqs(prop.Val.Bytes)
				if err != nil {
					return nil, err
				}
			} else {
				switch prop.Key {
				case "$and-array":
					var andArray []asn1.RawValue
					data, err = asn1.Unmarshal(prop.Val.FullBytes, &andArray)
					if err != nil {
						return nil, fmt.Errorf("failed to ASN.1 parse launch contraint '$and-array' properties: %v", err)
					}
					req[prop.Key] = make([]any, 0, len(andArray))
					for _, and := range andArray {
						r, err := parseReqs(and.FullBytes)
						if err != nil {
							return nil, err
						}
						req[prop.Key] = append(req[prop.Key].([]any), r)
					}
				case "$or-array":
					var orArray []asn1.RawValue
					data, err = asn1.Unmarshal(prop.Val.FullBytes, &orArray)
					if err != nil {
						return nil, fmt.Errorf("failed to ASN.1 parse launch contraint '$or-array' properties: %v", err)
					}
					req[prop.Key] = make([]any, 0, len(orArray))
					for _, or := range orArray {
						r, err := parseReqs(or.FullBytes)
						if err != nil {
							return nil, err
						}
						req[prop.Key] = append(req[prop.Key].([]any), r)
					}
				case "$query":
					var query []Entitlement
					data, err = asn1.Unmarshal(prop.Val.FullBytes, &query)
					if err != nil {
						return nil, fmt.Errorf("failed to ASN.1 parse launch contraint '$query' properties: %v", err)
					}
					req[prop.Key] = make([][]any, 0, len(query))
					for _, q := range query {
						val, err := getValue(q.Value)
						if err != nil {
							return nil, err
						}
						req[prop.Key] = append(req[prop.Key].([][]any), []any{q.Operation, val})
					}
				case "$in":
					var ins []asn1.RawValue
					data, err = asn1.Unmarshal(prop.Val.FullBytes, &ins)
					if err != nil {
						return nil, fmt.Errorf("failed to ASN.1 parse launch contraint '$in' properties: %v", err)
					}
					req[prop.Key] = make([]any, 0, len(ins))
					for _, in := range ins {
						val, err := getValue(in)
						if err != nil {
							return nil, err
						}
						req[prop.Key] = append(req[prop.Key].([]any), val)
					}
				default:
					return nil, fmt.Errorf("unsupported requirement key %s (notify author)", prop.Key)
				}
			}
		} else {
			val, err := getValue(prop.Val)
			if err != nil {
				return nil, err
			}
			req[prop.Key] = val
		}
	}

	return req, nil
}

// ParseLaunchContraints parses the launch constraint bytes
func ParseLaunchContraints(data []byte) (*LaunchContraints, error) {
	var err error

	lc := &LaunchContraints{
		Requirements: make(map[string]any),
	}

	var l launchContraints
	data, err = asn1.UnmarshalWithParams(data, &l, typeAppl)
	if err != nil {
		return nil, fmt.Errorf("failed to ASN.1 parse launch contraint inner data: %v", err)
	}

	lc.Count = l.Count

	for _, field := range l.Fields {
		switch field.Key {
		case "ccat":
			i, err := parseInt64(field.Val.Bytes)
			if err != nil {
				return nil, err
			}
			lc.CCAT = i
		case "comp":
			i, err := parseInt64(field.Val.Bytes)
			if err != nil {
				return nil, err
			}
			lc.COMP = i
		case "reqs":
			reqs, err := parseReqs(field.Val.Bytes)
			if err != nil {
				return nil, err
			}
			for k, v := range reqs {
				lc.Requirements[k] = v
			}
		case "vers":
			i, err := parseInt64(field.Val.Bytes)
			if err != nil {
				return nil, err
			}
			lc.Version = i
		default:
			return nil, fmt.Errorf("unsupported launch contraint field %s (notify author)", field.Key)
		}
	}

	return lc, nil
}
