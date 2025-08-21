package objc

import (
	"fmt"
	"strings"
)

// ref - https://developer.apple.com/library/archive/documentation/Cocoa/Conceptual/ObjCRuntimeGuide/Articles/ocrtTypeEncodings.html

var typeEncoding = map[string]string{
	"@": "id",
	"#": "Class",
	":": "SEL",
	"c": "char",
	"C": "unsigned char",
	"s": "short",
	"S": "unsigned short",
	"i": "int",
	"I": "unsigned int",
	"l": "long",
	"L": "unsigned long",
	"q": "long long",
	"Q": "unsigned long long",
	"t": "int128",
	"T": "unsigned int128",
	"f": "float",
	"d": "double",
	"D": "long double",
	"b": "bit field",
	"B": "BOOL",
	"v": "void",
	"z": "size_t",
	"Z": "int32",
	"w": "wchar_t",
	"?": "undefined",
	"^": "*",
	"*": "char *",
	"%": "NXAtom",
	// "[":  "", // _C_ARY_B
	// "]":  "", // _C_ARY_E
	// "(":  "", // _C_UNION_B
	// ")":  "", // _C_UNION_E
	// "{":  "", // _C_STRUCT_B
	// "}":  "", // _C_STRUCT_E
	"!":  "vector",
	"Vv": "void",
	"^?": "void *",         // void *
	"@?": "id /* block */", // block type
	// "@?": "void (^)(void)", // block type
}
var typeSpecifiers = map[string]string{
	"A": "atomic",
	"j": "_Complex",
	"!": "vector",
	"r": "const",
	"n": "in",
	"N": "inout",
	"o": "out",
	"O": "by copy",
	"R": "by ref",
	"V": "one way",
	"+": "gnu register",
}

const (
	propertyReadOnly  = "R" // property is read-only.
	propertyBycopy    = "C" // property is a copy of the value last assigned
	propertyByref     = "&" // property is a reference to the value last assigned
	propertyDynamic   = "D" // property is dynamic
	propertyGetter    = "G" // followed by getter selector name
	propertySetter    = "S" // followed by setter selector name
	propertyIVar      = "V" // followed by instance variable  name
	propertyType      = "T" // followed by old-style type encoding.
	propertyWeak      = "W" // 'weak' property
	propertyStrong    = "P" // property GC'able
	propertyAtomic    = "A" // property atomic
	propertyNonAtomic = "N" // property non-atomic
)

type methodEncodedArg struct {
	DecType   string // decoded argument type
	EncType   string // encoded argument type
	StackSize string // variable stack size
}

type varEncodedType struct {
	Head     string // type specifiers
	Variable string // variable name
	DecType  string // decoded variable type
	EncType  string // encoded variable type
	Tail     string // type output tail
}

// decodeMethodTypes decodes the method types and returns a return type and the argument types
func decodeMethodTypes(encodedTypes string) (string, []string) {
	var argTypes []string

	// skip return type
	encArgs := strings.TrimLeft(skipFirstType(encodedTypes), "0123456789")

	for idx, arg := range getArguments(encArgs) {
		switch idx {
		case 0:
			argTypes = append(argTypes, fmt.Sprintf("(%s)self", arg.DecType))
		case 1:
			argTypes = append(argTypes, fmt.Sprintf("(%s)id", arg.DecType))
		default:
			argTypes = append(argTypes, fmt.Sprintf("(%s)arg%d", arg.DecType, idx-1))
		}
	}
	return getReturnType(encodedTypes), argTypes
}

func getMethodWithArgs(method, returnType string, args []string) string {
	if len(args) <= 2 {
		return fmt.Sprintf("(%s)%s;", returnType, method)
	}
	args = args[2:] // skip self and SEL

	parts := strings.Split(method, ":")

	var methodStr string
	if len(parts) > 1 { // method has arguments based on SEL having ':'
		for idx, part := range parts {
			if len(part) == 0 || idx >= len(args) {
				break
			}
			methodStr += fmt.Sprintf("%s:%s ", part, args[idx])
		}
		return fmt.Sprintf("(%s)%s;", returnType, strings.TrimSpace(methodStr))
	}
	// method has no arguments based on SEL not having ':'
	return fmt.Sprintf("(%s)%s;", returnType, method)
}

func getPropertyAttributeTypes(attrs string) string {
	// var ivarStr string
	var typeStr string
	var attrsStr string
	var attrsList []string

	for _, attr := range strings.Split(attrs, ",") {
		if strings.HasPrefix(attr, propertyType) {
			attr = strings.TrimPrefix(attr, propertyType)
			if strings.HasPrefix(attr, "@\"") {
				typeStr = strings.Trim(attr, "@\"") + " *"
			} else {
				if val, ok := typeEncoding[attr]; ok {
					typeStr = val + " "
				}
			}
		} else if strings.HasPrefix(attr, propertyIVar) {
			// found ivar name
			// ivarStr = strings.TrimPrefix(attr, propertyIVar)
			continue
		} else {
			// TODO: handle the following cases
			// @property struct YorkshireTeaStruct structDefault; ==> T{YorkshireTeaStruct="pot"i"lady"c},VstructDefault
			// @property int (*functionPointerDefault)(char *);   ==> T^?,VfunctionPointerDefault
			switch attr {
			case propertyGetter:
				attr = strings.TrimPrefix(attr, propertyGetter)
				attrsList = append(attrsList, fmt.Sprintf("getter=%s", attr))
			case propertySetter:
				attr = strings.TrimPrefix(attr, propertySetter)
				attrsList = append(attrsList, fmt.Sprintf("setter=%s", attr))
			case propertyReadOnly:
				attrsList = append(attrsList, "readonly")
			case propertyNonAtomic:
				attrsList = append(attrsList, "nonatomic")
			case propertyAtomic:
				attrsList = append(attrsList, "atomic")
			case propertyBycopy:
				attrsList = append(attrsList, "copy")
			case propertyByref:
				attrsList = append(attrsList, "retain")
			case propertyWeak:
				attrsList = append(attrsList, "weak")
			case propertyDynamic:
				attrsList = append(attrsList, "dynamic")
			case propertyStrong:
				attrsList = append(attrsList, "collectable")
			}
		}
	}

	if len(attrsList) > 0 {
		attrsStr = fmt.Sprintf("(%s) ", strings.Join(attrsList, ", "))
	}

	return fmt.Sprintf("%s%s", attrsStr, typeStr)
}

func getIVarType(ivType string) string {
	if strings.HasPrefix(ivType, "@\"") && len(ivType) > 1 {
		return strings.Trim(ivType, "@\"") + " *"
	}
	return decodeType(ivType) + " "
}

func getReturnType(types string) string {
	if len(types) == 0 {
		return ""
	}
	return decodeType(strings.TrimSuffix(types, skipFirstType(types)))
}

func decodeType(encType string) string {
	var s string

	if len(encType) == 0 {
		return ""
	}

	if strings.HasPrefix(encType, "^") {
		return decodeType(encType[1:]) + " *" // pointer
	}

	if strings.HasPrefix(encType, "@?") {
		if len(encType) > 2 { // TODO: remove this??
			pointerType := decodeType(encType[2:])
			return "id  ^" + pointerType
		} else {
			return "id /* block */"
		}
	}

	if strings.HasPrefix(encType, "^?") {
		if len(encType) > 2 {
			pointerType := decodeType(encType[2:])
			return "void * " + pointerType
		} else {
			return "void *"
		}
	}

	if typ, ok := typeEncoding[encType]; ok {
		return typ
	}

	if spec, ok := typeSpecifiers[string(encType[0])]; ok {
		return spec + " " + decodeType(encType[1:])
	}

	if strings.HasPrefix(encType, "@\"") && len(encType) > 1 {
		return strings.Trim(encType, "@\"") + " *"
	}

	if strings.HasPrefix(encType, "b") {
		return decodeBitfield(encType)
	}

	if len(encType) > 2 {
		if strings.HasPrefix(encType, "[") { // ARRAY
			inner := encType[strings.IndexByte(encType, '[')+1 : strings.LastIndexByte(encType, ']')]
			s += decodeArray(inner)
		} else if strings.HasPrefix(encType, "{") { // STRUCT
			inner := encType[strings.IndexByte(encType, '{')+1 : strings.LastIndexByte(encType, '}')]
			s += decodeStructure(inner)
		} else if strings.HasPrefix(encType, "(") { // UNION
			inner := encType[strings.IndexByte(encType, '(')+1 : strings.LastIndexByte(encType, ')')]
			s += decodeUnion(inner)
		}
	}

	if len(s) == 0 {
		return encType
	}

	return s
}

func decodeArray(arrayType string) string {
	numIdx := strings.LastIndexAny(arrayType, "0123456789")
	if len(arrayType) == 1 {
		return fmt.Sprintf("x[%s]", arrayType)
	}
	return fmt.Sprintf("%s x[%s]", decodeType(arrayType[numIdx+1:]), arrayType[:numIdx+1])
}

func decodeStructure(structure string) string {
	return decodeStructOrUnion(structure, "struct")
}

func decodeUnion(unionType string) string {
	return decodeStructOrUnion(unionType, "union")
}

func decodeBitfield(bitfield string) string {
	span := encodingGetSizeOfArguments(bitfield)
	return fmt.Sprintf("unsigned int x :%d", span)
}

func getFieldName(field string) (string, string) {
	if strings.HasPrefix(field, "\"") {
		name, rest, ok := strings.Cut(strings.TrimPrefix(field, "\""), "\"")
		if !ok {
			return "", field
		}
		return name, rest
	}
	return "", field
}

func decodeStructOrUnion(typ, kind string) string {
	name, rest, _ := strings.Cut(typ, "=")
	if name == "?" {
		name = ""
	} else {
		name += " "
	}

	if len(rest) == 0 {
		return fmt.Sprintf("%s %s", kind, strings.TrimSpace(name))
	}

	var idx int
	var fields []string
	fieldName, rest := getFieldName(rest)
	field, rest, ok := CutType(rest)

	for ok {
		if fieldName != "" {
			dtype := decodeType(field)
			if strings.HasSuffix(dtype, " *") {
				fields = append(fields, fmt.Sprintf("%s%s;", dtype, fieldName))
			} else {
				fields = append(fields, fmt.Sprintf("%s %s;", dtype, fieldName))
			}
		} else {
			if strings.HasPrefix(field, "b") {
				span := encodingGetSizeOfArguments(field)
				fields = append(fields, fmt.Sprintf("unsigned int x%d :%d;", idx, span))
			} else if strings.HasPrefix(field, "[") {
				array := decodeType(field)
				array = strings.TrimSpace(strings.Replace(array, "x", fmt.Sprintf("x%d", idx), 1))
				fields = append(fields, array)
			} else {
				fields = append(fields, fmt.Sprintf("%s x%d;", decodeType(field), idx))
			}
			idx++
		}
		fieldName, rest = getFieldName(rest)
		field, rest, ok = CutType(rest)
	}

	return fmt.Sprintf("%s %s{ %s }", kind, name, strings.Join(fields, " "))
}

func skipFirstType(typStr string) string {
	i := 0
	typ := []byte(typStr)
	for {
		switch typ[i] {
		case 'O': /* bycopy */
			fallthrough
		case 'n': /* in */
			fallthrough
		case 'o': /* out */
			fallthrough
		case 'N': /* inout */
			fallthrough
		case 'r': /* const */
			fallthrough
		case 'V': /* oneway */
			fallthrough
		case '^': /* pointers */
			i++
		case '@': /* objects */
			if i+1 < len(typ) && typ[i+1] == '?' {
				i++ /* Blocks */
			} else if i+1 < len(typ) && typ[i+1] == '"' {
				i++
				for typ[i+1] != '"' {
					i++ /* Class */
				}
				i++
			}
			return string(typ[i+1:])
		case '[': /* arrays */
			i++
			for typ[i] >= '0' && typ[i] <= '9' {
				i++
			}
			return string(typ[i+subtypeUntil(string(typ[i:]), ']')+1:])
		case '{': /* structures */
			i++
			return string(typ[i+subtypeUntil(string(typ[i:]), '}')+1:])
		case '(': /* unions */
			i++
			return string(typ[i+subtypeUntil(string(typ[i:]), ')')+1:])
		default: /* basic types */
			i++
			return string(typ[i:])
		}
	}
}

func subtypeUntil(typ string, end byte) int {
	var level int
	head := typ
	for len(typ) > 0 {
		if level == 0 && typ[0] == end {
			return len(head) - len(typ)
		}
		switch typ[0] {
		case ']', '}', ')':
			level -= 1
		case '[', '{', '(':
			level += 1
		}
		typ = typ[1:]
	}
	return 0
}

func CutType(typStr string) (string, string, bool) {
	var i int

	if len(typStr) == 0 {
		return "", "", false
	}
	if len(typStr) == 1 {
		return typStr, "", true
	}

	typ := []byte(typStr)
	for {
		switch typ[i] {
		case 'O': /* bycopy */
			fallthrough
		case 'n': /* in */
			fallthrough
		case 'o': /* out */
			fallthrough
		case 'N': /* inout */
			fallthrough
		case 'r': /* const */
			fallthrough
		case 'V': /* oneway */
			fallthrough
		case '^': /* pointers */
			i++
		case '@': /* objects */
			if i+1 < len(typ) && typ[i+1] == '?' {
				i++ /* Blocks */
			} else if i+1 < len(typ) && typ[i+1] == '"' {
				i++
				for typ[i+1] != '"' {
					i++ /* Class */
				}
				i++
			}
			return string(typ[:i+1]), string(typ[i+1:]), true
		case 'b': /* bitfields */
			i++
			for i < len(typ) && typ[i] >= '0' && typ[i] <= '9' {
				i++
			}
			return string(typ[:i]), string(typ[i:]), true
		case '[': /* arrays */
			i++
			for typ[i] >= '0' && typ[i] <= '9' {
				i++
			}
			return string(typ[:i+subtypeUntil(string(typ[i:]), ']')+1]), string(typ[i+subtypeUntil(string(typ[i:]), ']')+1:]), true
		case '{': /* structures */
			i++
			return string(typ[:i+subtypeUntil(string(typ[i:]), '}')+1]), string(typ[i+subtypeUntil(string(typ[i:]), '}')+1:]), true
		case '(': /* unions */
			i++
			return string(typ[:i+subtypeUntil(string(typ[i:]), ')')+1]), string(typ[i+subtypeUntil(string(typ[i:]), ')')+1:]), true
		default: /* basic types */
			i++
			return string(typ[:i]), string(typ[i:]), true
		}
	}
}

func getNumberOfArguments(types string) int {
	var nargs int
	// First, skip the return type
	types = skipFirstType(types)
	// Next, skip stack size
	types = strings.TrimLeft(types, "0123456789")
	// Now, we have the arguments - count how many
	for len(types) > 0 {
		// Traverse argument type
		types = skipFirstType(types)
		// Skip GNU runtime's register parameter hint
		types = strings.TrimPrefix(types, "+")
		// Traverse (possibly negative) argument offset
		types = strings.TrimPrefix(types, "-")
		types = strings.TrimLeft(types, "0123456789")
		// Made it past an argument
		nargs++
	}
	return nargs
}

func getArguments(encArgs string) []methodEncodedArg {
	var args []methodEncodedArg
	t, rest, ok := CutType(encArgs)
	for ok {
		i := 0
		for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
			i++
		}
		args = append(args, methodEncodedArg{
			EncType:   t,
			DecType:   decodeType(t),
			StackSize: rest[:i],
		})
		t, rest, ok = CutType(rest[i:])
	}
	return args
}

func encodingGetSizeOfArguments(typedesc string) uint {
	var stackSize uint

	stackSize = 0
	typedesc = skipFirstType(typedesc)

	for i := 0; i < len(typedesc); i++ {
		if typedesc[i] >= '0' && typedesc[i] <= '9' {
			stackSize = stackSize*10 + uint(typedesc[i]-'0')
		} else {
			break
		}
	}

	return stackSize
}
