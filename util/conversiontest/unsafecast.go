/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conversiontest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

// UnsafeCastInput contains configuration for CheckUnsafeStructCasts.
type UnsafeCastInput struct {
	// Scheme is the runtime scheme containing registered types.
	Scheme *runtime.Scheme
	// PackagePath is the full package import path (e.g., "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1").
	PackagePath string
	// GeneratedFile is the absolute path to zz_generated.conversion.go.
	GeneratedFile string
	// ManualConversionDirs contains absolute paths to directories with hand-written Convert_* functions.
	ManualConversionDirs []string
	// ExtraTypes is a slice of zero values of types not reachable from Scheme kinds.
	ExtraTypes []any
}

// CheckUnsafeStructCasts verifies that unsafe whole-struct casts in generated conversions
// are layout-safe and don't bypass hand-written field converters.
func CheckUnsafeStructCasts(tb testing.TB, in UnsafeCastInput) {
	tb.Helper()
	// Build type registry from scheme
	reg := buildTypeRegistry(in.Scheme, in.ExtraTypes)

	// Parse generated file and extract imports
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, in.GeneratedFile, nil, 0)
	if err != nil {
		tb.Fatalf("failed to parse %s: %v", in.GeneratedFile, err)
	}

	imports := extractImports(f)
	local := in.PackagePath

	// Collect manual conversions from hand-written files
	manual := collectManualConversions(in.ManualConversionDirs, local)

	// Check each autoConvert_* function
	casts, unresolved := 0, 0
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || !strings.HasPrefix(fd.Name.Name, "autoConvert_") {
			continue
		}
		if len(fd.Body.List) == 0 {
			continue
		}

		// Check if first statement is whole-struct unsafe.Pointer cast
		as, ok := fd.Body.List[0].(*ast.AssignStmt)
		if !ok {
			continue
		}
		if len(as.Lhs) == 0 {
			continue
		}
		st, ok := as.Lhs[0].(*ast.StarExpr)
		if !ok {
			continue
		}
		id, ok := st.X.(*ast.Ident)
		if !ok || id.Name != "out" {
			continue
		}

		casts++

		// Extract in and out type keys
		ps := fd.Type.Params.List
		inK := typeKey(ps[0].Type, local, imports)
		outK := inK
		if len(ps[0].Names) > 0 && len(ps) > 1 && len(ps[1].Names) > 0 {
			outK = typeKey(ps[1].Type, local, imports)
		}

		in, ok1 := reg[inK]
		out, ok2 := reg[outK]
		if !ok1 || !ok2 {
			unresolved++
			tb.Logf("unresolved %s (%s -> %s)", fd.Name.Name, inK, outK)
			if !ok1 {
				tb.Errorf("type %s not in scheme or ExtraTypes; add a zero value to ExtraTypes", inK)
			}
			if !ok2 {
				tb.Errorf("type %s not in scheme or ExtraTypes; add a zero value to ExtraTypes", outK)
			}
			continue
		}

		// Check layout safety
		if !sameLayout(in, out, map[typePair]bool{}) {
			tb.Errorf("%s: unsafe whole-struct cast between layout-different types %s and %s; run make generate", fd.Name.Name, inK, outK)
		}

		// Check for bypassed manual conversions
		if in.Kind() == reflect.Struct && out.Kind() == reflect.Struct && in.NumField() == out.NumField() {
			for i := range in.NumField() {
				fi, fo := in.Field(i), out.Field(i)
				fiKey := fi.Type.PkgPath() + "." + fi.Type.Name()
				foKey := fo.Type.PkgPath() + "." + fo.Type.Name()
				if m, ok := manual[[2]string{fiKey, foKey}]; ok {
					tb.Errorf("%s: unsafe cast bypasses hand-written %s for field %s; run make generate", fd.Name.Name, m, fi.Name)
				}
			}
		}
	}

	tb.Logf("whole-struct casts: %d, unresolved: %d", casts, unresolved)
}

type typePair struct{ a, b reflect.Type }

func buildTypeRegistry(scheme *runtime.Scheme, extra []any) map[string]reflect.Type {
	reg := make(map[string]reflect.Type)
	for _, typ := range scheme.AllKnownTypes() {
		walkTypes(typ, reg)
	}
	for _, e := range extra {
		walkTypes(reflect.TypeOf(e), reg)
	}
	return reg
}

func walkTypes(t reflect.Type, reg map[string]reflect.Type) {
	//nolint:govet // False positive on reflect constants
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice || t.Kind() == reflect.Array || t.Kind() == reflect.Map {
		if t.Kind() == reflect.Map {
			walkTypes(t.Key(), reg)
		}
		t = t.Elem()
	}
	if t.Name() != "" && t.PkgPath() != "" {
		k := t.PkgPath() + "." + t.Name()
		if _, ok := reg[k]; ok {
			return
		}
		reg[k] = t
	}
	if t.Kind() == reflect.Struct {
		for i := range t.NumField() {
			walkTypes(t.Field(i).Type, reg)
		}
	}
}

func extractImports(f *ast.File) map[string]string {
	imports := make(map[string]string)
	for _, im := range f.Imports {
		p, _ := strconv.Unquote(im.Path.Value)
		n := filepath.Base(p)
		if im.Name != nil {
			n = im.Name.Name
		}
		imports[n] = p
	}
	return imports
}

func typeKey(e ast.Expr, local string, imports map[string]string) string {
	if s, ok := e.(*ast.StarExpr); ok {
		e = s.X
	}
	switch x := e.(type) {
	case *ast.Ident:
		return local + "." + x.Name
	case *ast.SelectorExpr:
		if id, ok := x.X.(*ast.Ident); ok {
			return imports[id.Name] + "." + x.Sel.Name
		}
	}
	return ""
}

func sameLayout(a, b reflect.Type, seen map[typePair]bool) bool {
	if a == b {
		return true
	}
	p := typePair{a, b}
	if seen[p] {
		return true
	}
	seen[p] = true
	if a.Kind() != b.Kind() || a.Size() != b.Size() {
		return false
	}
	switch a.Kind() {
	case reflect.Struct:
		if a.NumField() != b.NumField() {
			return false
		}
		for i := range a.NumField() {
			fa, fb := a.Field(i), b.Field(i)
			if fa.Name != fb.Name || fa.Offset != fb.Offset || !sameLayout(fa.Type, fb.Type, seen) {
				return false
			}
		}
	case reflect.Ptr, reflect.Slice: //nolint:govet // False positive on reflect constants
		return sameLayout(a.Elem(), b.Elem(), seen)
	case reflect.Array:
		return a.Len() == b.Len() && sameLayout(a.Elem(), b.Elem(), seen)
	case reflect.Map:
		return sameLayout(a.Key(), b.Key(), seen) && sameLayout(a.Elem(), b.Elem(), seen)
	case reflect.Interface, reflect.Func, reflect.Chan:
		return false
	}
	return true
}

func collectManualConversions(dirs []string, local string) map[[2]string]string {
	manual := make(map[[2]string]string)
	fset := token.NewFileSet()

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			fn := entry.Name()
			if strings.HasPrefix(fn, "zz_generated") || strings.HasSuffix(fn, "_test.go") {
				continue
			}
			if !strings.HasSuffix(fn, ".go") {
				continue
			}

			src, err := os.ReadFile(filepath.Join(dir, fn)) //nolint:gosec // dir is validated, fn is from os.ReadDir
			if err != nil {
				continue
			}

			mf, err := parser.ParseFile(fset, fn, src, parser.ParseComments)
			if err != nil {
				continue
			}

			mi := extractImports(mf)

			for _, d := range mf.Decls {
				fd, ok := d.(*ast.FuncDecl)
				if !ok || fd.Recv != nil {
					continue
				}
				if !strings.HasPrefix(fd.Name.Name, "Convert_") {
					continue
				}
				if fd.Type.Params.NumFields() != 3 {
					continue
				}

				// Check for copy-only comment
				if fd.Doc != nil {
					for _, comment := range fd.Doc.List {
						if strings.Contains(comment.Text, "+k8s:conversion-fn=copy-only") {
							goto next_func
						}
					}
				}

				{
					ps := fd.Type.Params.List
					var exprs []ast.Expr
					for _, p := range ps {
						for range p.Names {
							exprs = append(exprs, p.Type)
						}
					}
					if len(exprs) >= 2 {
						inK := typeKey(exprs[0], local, mi)
						outK := typeKey(exprs[1], local, mi)
						if inK != "" && outK != "" {
							manual[[2]string{inK, outK}] = fd.Name.Name
						}
					}
				}

			next_func:
			}
		}
	}

	return manual
}
