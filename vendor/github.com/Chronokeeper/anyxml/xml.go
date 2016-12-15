// Copyright 2012-2016 xiaolipeng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

// xml.go - basically the core of X2j for map[string]interface{} values.
//          NewMapXml, NewMapXmlReader, mv.Xml, mv.XmlWriter
// see x2j and j2x for wrappers to provide end-to-end transformation of XML and JSON messages.

package anyxml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"time"
)

// --------------------------------- Xml, XmlIndent - from mxj -------------------------------

const (
	DefaultRootTag          = "doc"
	UseGoEmptyElementSyntax = false // if 'true' encode empty element as "<tag></tag>" instead of "<tag/>
)

// From: github.com/clbanning/mxj/xml.go with functions relabled: Xml() --> anyxml().
// Encode a Map as XML.  The companion of NewMapXml().
// The following rules apply.
//    - The key label "#text" is treated as the value for a simple element with attributes.
//    - Map keys that begin with a hyphen, '-', are interpreted as attributes.
//      It is an error if the attribute doesn't have a []byte, string, number, or boolean value.
//    - Map value type encoding:
//          > string, bool, float64, int, int32, int64, float32: per "%v" formating
//          > []bool, []uint8: by casting to string
//          > structures, etc.: handed to xml.Marshal() - if there is an error, the element
//            value is "UNKNOWN"
//    - Elements with only attribute values or are null are terminated using "/>".
//    - If len(mv) == 1 and no rootTag is provided, then the map key is used as the root tag, possible.
//      Thus, `{ "key":"value" }` encodes as "<key>value</key>".
//    - To encode empty elements in a syntax consistent with encoding/xml call UseGoXmlEmptyElementSyntax().
func anyxml(m map[string]interface{}, rootTag ...string) ([]byte, error) {
	var err error
	s := new(string)
	p := new(pretty) // just a stub

	if len(m) == 1 && len(rootTag) == 0 {
		for key, value := range m {
			// if it an array, see if all values are map[string]interface{}
			// we force a new root tag if we'll end up with no key:value in the list
			// so: key:[string_val, bool:true] --> <doc><key>string_val</key><bool>true</bool></key></doc>
			switch value.(type) {
			case []interface{}:
				for _, v := range value.([]interface{}) {
					switch v.(type) {
					case map[string]interface{}: // noop
					default: // anything else
						err = mapToXmlIndent(false, s, DefaultRootTag, m, p)
						goto done
					}
				}
			}
			err = mapToXmlIndent(false, s, key, value, p)
		}
	} else if len(rootTag) == 1 {
		err = mapToXmlIndent(false, s, rootTag[0], m, p)
	} else {
		err = mapToXmlIndent(false, s, DefaultRootTag, m, p)
	}
done:
	return []byte(*s), err
}

func anyxmlWithDateFormat(dateFormat string, m map[string]interface{}, rootTag ...string) ([]byte, error) {
	var err error
	s := new(string)
	p := new(pretty) // just a stub

	if len(m) == 1 && len(rootTag) == 0 {
		for key, value := range m {
			// if it an array, see if all values are map[string]interface{}
			// we force a new root tag if we'll end up with no key:value in the list
			// so: key:[string_val, bool:true] --> <doc><key>string_val</key><bool>true</bool></key></doc>
			switch value.(type) {
			case []interface{}:
				for _, v := range value.([]interface{}) {
					switch v.(type) {
					case map[string]interface{}: // noop
					default: // anything else
						err = mapToXmlIndentWithDateFormat(dateFormat, false, s, DefaultRootTag, m, p)
						goto done
					}
				}
			}
			err = mapToXmlIndentWithDateFormat(dateFormat, false, s, key, value, p)
		}
	} else if len(rootTag) == 1 {
		err = mapToXmlIndentWithDateFormat(dateFormat, false, s, rootTag[0], m, p)
	} else {
		err = mapToXmlIndentWithDateFormat(dateFormat, false, s, DefaultRootTag, m, p)
	}
done:
	return []byte(*s), err
}

func anyxmlIndentWithDateFormat(dateFormat string, m map[string]interface{}, prefix string, indent string, rootTag ...string) ([]byte, error) {
	var err error
	s := new(string)
	p := new(pretty) // just a stub

	if len(m) == 1 && len(rootTag) == 0 {
		for key, value := range m {
			// if it an array, see if all values are map[string]interface{}
			// we force a new root tag if we'll end up with no key:value in the list
			// so: key:[string_val, bool:true] --> <doc><key>string_val</key><bool>true</bool></key></doc>
			switch value.(type) {
			case []interface{}:
				for _, v := range value.([]interface{}) {
					switch v.(type) {
					case map[string]interface{}: // noop
					default: // anything else
						err = mapToXmlIndentWithDateFormat(dateFormat, true, s, DefaultRootTag, m, p)
						goto done
					}
				}
			}
			err = mapToXmlIndentWithDateFormat(dateFormat, true, s, key, value, p)
		}
	} else if len(rootTag) == 1 {
		err = mapToXmlIndentWithDateFormat(dateFormat, true, s, rootTag[0], m, p)
	} else {
		err = mapToXmlIndentWithDateFormat(dateFormat, true, s, DefaultRootTag, m, p)
	}
done:
	return []byte(*s), err
}

// Encode a map[string]interface{} as a pretty XML string.
// See Xml for encoding rules.
func anyxmlIndent(m map[string]interface{}, prefix string, indent string, rootTag ...string) ([]byte, error) {
	var err error
	s := new(string)
	p := new(pretty)
	p.indent = indent
	p.padding = prefix

	if len(m) == 1 && len(rootTag) == 0 {
		// this can extract the key for the single map element
		// use it if it isn't a key for a list
		for key, value := range m {
			if _, ok := value.([]interface{}); ok {
				err = mapToXmlIndent(true, s, DefaultRootTag, m, p)
			} else {
				err = mapToXmlIndent(true, s, key, value, p)
			}
		}
	} else if len(rootTag) == 1 {
		err = mapToXmlIndent(true, s, rootTag[0], m, p)
	} else {
		err = mapToXmlIndent(true, s, DefaultRootTag, m, p)
	}
	return []byte(*s), err
}

type pretty struct {
	indent   string
	cnt      int
	padding  string
	mapDepth int
	start    int
}

func (p *pretty) Indent() {
	p.padding += p.indent
	p.cnt++
}

func (p *pretty) Outdent() {
	if p.cnt > 0 {
		p.padding = p.padding[:len(p.padding)-len(p.indent)]
		p.cnt--
	}
}

// where the work actually happens
// returns an error if an attribute is not atomic
func mapToXmlIndent(doIndent bool, s *string, key string, value interface{}, pp *pretty) error {
	var endTag bool
	var isSimple bool
	p := &pretty{pp.indent, pp.cnt, pp.padding, pp.mapDepth, pp.start}
	key = template.HTMLEscapeString(key)

	switch value.(type) {
	case map[string]interface{}, []byte, string, float64, bool, int, int32, int64, float32:
		if doIndent {
			*s += p.padding
		}
		*s += `<` + key
	}
	switch value.(type) {
	case map[string]interface{}:
		vv := value.(map[string]interface{})
		lenvv := len(vv)
		// scan out attributes - keys have prepended hyphen, '-'
		var cntAttr int
		for k, v := range vv {
			if k[:1] == "-" {
				switch v.(type) {
				case string:
					*s += ` ` + k[1:] + `="` + fmt.Sprintf("%v", template.HTMLEscapeString(v.(string))) + `"`
					cntAttr++
				case float64, bool, int, int32, int64, float32:
					*s += ` ` + k[1:] + `="` + fmt.Sprintf("%v", v) + `"`
					cntAttr++
				case []byte: // allow standard xml pkg []byte transform, as below
					*s += ` ` + k[1:] + `="` + fmt.Sprintf("%v", string(v.([]byte))) + `"`
					cntAttr++
				default:
					return fmt.Errorf("invalid attribute value for: %s", k)
				}
			}
		}
		// only attributes?
		if cntAttr == lenvv {
			break
		}
		// simple element? Note: '#text" is an invalid XML tag.
		if v, ok := vv["#text"]; ok {
			if cntAttr+1 < lenvv {
				return errors.New("#text key occurs with other non-attribute keys")
			}
			*s += ">" + fmt.Sprintf("%v", v)
			endTag = true
			break
		}
		// close tag with possible attributes
		*s += ">"
		if doIndent {
			*s += "\n"
		}
		// something more complex
		p.mapDepth++
		var i int
		for k, v := range vv {
			if k[:1] == "-" {
				continue
			}
			switch v.(type) {
			case []interface{}:
			default:
				if i == 0 && doIndent {
					p.Indent()
				}
			}
			i++
			mapToXmlIndent(doIndent, s, k, v, p)
			switch v.(type) {
			case []interface{}: // handled in []interface{} case
			default:
				if doIndent {
					p.Outdent()
				}
			}
			i--
		}
		p.mapDepth--
		endTag = true
	case []interface{}:
		for _, v := range value.([]interface{}) {
			if doIndent {
				p.Indent()
			}
			mapToXmlIndent(doIndent, s, key, v, p)
			if doIndent {
				p.Outdent()
			}
		}
		return nil
	case nil:
		// terminate the tag
		*s += "<" + key
		break
	default: // handle anything - even goofy stuff
		switch value.(type) {
		case string:
			*s += ">" + fmt.Sprintf("%v", template.HTMLEscapeString(value.(string)))
		case float64, bool, int, int32, int64, float32:
			*s += ">" + fmt.Sprintf("%v", value)
		case []byte: // NOTE: byte is just an alias for uint8
			// similar to how xml.Marshal handles []byte structure members
			*s += ">" + string(value.([]byte))
		default:
			var v []byte
			var err error
			if doIndent {
				v, err = xml.MarshalIndent(value, p.padding, p.indent)
			} else {
				v, err = xml.Marshal(value)
			}
			if err != nil {
				*s += ">UNKNOWN"
			} else {
				*s += string(v)
			}
		}
		isSimple = true
		endTag = true
	}

	if endTag {
		if doIndent {
			if !isSimple {
				//				if p.mapDepth == 0 {
				//					p.Outdent()
				//				}
				*s += p.padding
			}
		}
		switch value.(type) {
		case map[string]interface{}, []byte, string, float64, bool, int, int32, int64, float32:
			*s += `</` + key + ">"
		}
	} else if UseGoEmptyElementSyntax {
		*s += "></" + key + ">"
	} else {
		*s += "/>"
	}
	if doIndent {
		if p.cnt > p.start {
			*s += "\n"
		}
		p.Outdent()
	}

	return nil
}

// where the work actually happens
// returns an error if an attribute is not atomic
func mapToXmlIndentWithDateFormat(dateFormat string, doIndent bool, s *string, key string, value interface{}, pp *pretty) error {
	var endTag bool
	var isSimple bool
	p := &pretty{pp.indent, pp.cnt, pp.padding, pp.mapDepth, pp.start}
	key = template.HTMLEscapeString(key)

	//start tag
	switch value.(type) {
	case map[string]interface{}, []byte, string, float64, bool, int, int32, int64, float32, time.Time:
		if doIndent {
			*s += p.padding
		}
		*s += `<` + key
	}

	switch value.(type) {
	case map[string]interface{}:
		vv := value.(map[string]interface{})
		lenvv := len(vv)
		var cntAttr int
		for k, v := range vv {
			if k[:1] == "-" {
				switch v.(type) {
				case string:
					*s += ` ` + k[1:] + `="` + fmt.Sprintf("%v", template.HTMLEscapeString(v.(string))) + `"`
					cntAttr++
				case float64, bool, int, int32, int64, float32:
					*s += ` ` + k[1:] + `="` + fmt.Sprintf("%v", v) + `"`
					cntAttr++
				case []byte: // allow standard xml pkg []byte transform, as below
					*s += ` ` + k[1:] + `="` + fmt.Sprintf("%v", string(v.([]byte))) + `"`
					cntAttr++
				default:
					return fmt.Errorf("invalid attribute value for: %s", k)
				}
			}
		}
		// only attributes?
		if cntAttr == lenvv {
			break
		}
		// simple element? Note: '#text" is an invalid XML tag.
		if v, ok := vv["#text"]; ok {
			if cntAttr+1 < lenvv {
				return errors.New("#text key occurs with other non-attribute keys")
			}
			*s += ">" + fmt.Sprintf("%v", v)
			endTag = true
			break
		}
		// close tag with possible attributes
		*s += ">"
		if doIndent {
			*s += "\n"
		}
		// something more complex
		p.mapDepth++
		var i int
		for k, v := range vv {
			if k[:1] == "-" {
				continue
			}
			switch v.(type) {
			case []interface{}:
			default:
				if i == 0 && doIndent {
					p.Indent()
				}
			}
			i++
			mapToXmlIndentWithDateFormat(dateFormat, doIndent, s, k, v, p)
			switch v.(type) {
			case []interface{}: // handled in []interface{} case
			default:
				if doIndent {
					p.Outdent()
				}
			}
			i--
		}
		p.mapDepth--
		endTag = true
	case []interface{}:
		for _, v := range value.([]interface{}) {
			if doIndent {
				p.Indent()
			}
			mapToXmlIndentWithDateFormat(dateFormat, doIndent, s, key, v, p)
			if doIndent {
				p.Outdent()
			}
		}
		return nil
	case nil:
		*s += "<" + key
		break
	default: // handle anything - even goofy stuff
		switch value.(type) {
		case string:
			*s += ">" + fmt.Sprintf("%v", template.HTMLEscapeString(value.(string)))
		case float64, bool, int, int32, int64, float32:
			*s += ">" + fmt.Sprintf("%v", value)
		case []byte: // NOTE: byte is just an alias for uint8
			// similar to how xml.Marshal handles []byte structure members
			*s += ">" + string(value.([]byte))
		case time.Time:
			*s += ">" + (value.(time.Time)).Format(dateFormat)
		default:
			var v []byte
			var err error
			if doIndent {
				v, err = xml.MarshalIndent(value, p.padding, p.indent)
			} else {
				v, err = xml.Marshal(value)
			}
			if err != nil {
				*s += ">UNKNOWN"
			} else {
				*s += string(v)
			}
		}
		isSimple = true
		endTag = true
	}

	if endTag {
		if doIndent {
			if !isSimple {
				//				if p.mapDepth == 0 {
				//					p.Outdent()
				//				}
				*s += p.padding
			}
		}
		switch value.(type) {
		case map[string]interface{}, []byte, string, float64, bool, int, int32, int64, float32, time.Time:
			*s += `</` + key + ">"
		}
	} else if UseGoEmptyElementSyntax {
		*s += "></" + key + ">"
	} else {
		*s += "/>"
	}
	if doIndent {
		if p.cnt > p.start {
			*s += "\n"
		}
		p.Outdent()
	}

	return nil
}
