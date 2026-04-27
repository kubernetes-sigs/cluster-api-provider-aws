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

// Package printer provides reusable text/json/yaml output renderers for
// release-tool subcommands. The shape mirrors
// cmd/clusterawsadm/printers/printers.go but is dependency-light: the table
// printer uses text/tabwriter and a small local Table type, so the package
// does not pull in k8s.io/cli-runtime or k8s.io/apimachinery.
package printer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"sigs.k8s.io/yaml"
)

// Type identifies an output format accepted by --output / -o.
type Type string

const (
	// TypeTable selects the human-readable table printer.
	TypeTable Type = "table"
	// TypeJSON selects the indented JSON printer.
	TypeJSON Type = "json"
	// TypeYAML selects the YAML printer.
	TypeYAML Type = "yaml"
)

var (
	// ErrUnknownPrinterType is returned by New for an unrecognized format.
	ErrUnknownPrinterType = errors.New("unknown printer type")
	// ErrTableRequired is returned by the table printer when the input is not *Table.
	ErrTableRequired = errors.New("*printer.Table is required for table output")
)

// Table is the input expected by the table printer. Callers convert their
// domain objects into a Table before printing; non-table printers accept any
// JSON-marshalable value directly.
type Table struct {
	// Columns is the header row.
	Columns []string
	// Rows is the body, one slice of cells per row. Each row should have the
	// same length as Columns.
	Rows [][]string
}

// Printer renders an arbitrary value to a configured io.Writer.
type Printer interface {
	// Print renders in to the printer's underlying writer.
	Print(in interface{}) error
}

// New returns a Printer for the requested output format. The match is
// case-insensitive; an empty value falls back to TypeTable for ergonomics.
func New(printerType string, w io.Writer) (Printer, error) {
	switch Type(strings.ToLower(strings.TrimSpace(printerType))) {
	case TypeTable, "":
		return &tablePrinter{w: w}, nil
	case TypeJSON:
		return &jsonPrinter{w: w}, nil
	case TypeYAML:
		return &yamlPrinter{w: w}, nil
	default:
		return nil, fmt.Errorf("%w %q: must be one of %s, %s, %s",
			ErrUnknownPrinterType, printerType, TypeTable, TypeJSON, TypeYAML)
	}
}

type tablePrinter struct{ w io.Writer }

func (p *tablePrinter) Print(in interface{}) error {
	t, ok := in.(*Table)
	if !ok {
		return ErrTableRequired
	}
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	if len(t.Columns) > 0 {
		fmt.Fprintln(tw, strings.Join(t.Columns, "\t"))
	}
	for _, row := range t.Rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	return tw.Flush()
}

type jsonPrinter struct{ w io.Writer }

func (p *jsonPrinter) Print(in interface{}) error {
	enc := json.NewEncoder(p.w)
	enc.SetIndent("", "  ")
	return enc.Encode(in)
}

type yamlPrinter struct{ w io.Writer }

func (p *yamlPrinter) Print(in interface{}) error {
	b, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("encoding yaml: %w", err)
	}
	_, err = p.w.Write(b)
	return err
}
