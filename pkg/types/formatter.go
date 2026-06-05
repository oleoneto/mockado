package types

import (
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
)

type (
	TableFormattable interface{ TableWriter() table.Writer }

	PlainFormatter  struct{}
	TableFormatter  struct{}
	SilentFormatter struct{}

	OutputFormat struct{ *FlagEnum[string] }

	TableOptions struct {
		Title        *string
		Header       table.Row
		Footer       *table.Row
		Style        *table.Style
		OutputMirror io.Writer
		SortOrder    table.SortMode
		SortField    string
		ColumnConfig *[]table.ColumnConfig
	}
)

func (SilentFormatter) Format(data any) ([]byte, error) { return []byte{}, nil }

func (*TableFormatter) Format(data any) ([]byte, error) {
	tw, ok := data.(TableFormattable)
	if !ok {
		return []byte{}, nil
	}

	return []byte(tw.TableWriter().Render()), nil
}

func (PlainFormatter) Format(v any) ([]byte, error) {
	type S interface{ String() string }

	s, ok := v.(S)
	if !ok {
		return []byte(fmt.Sprintf("%+v", v)), nil
	}

	return []byte(s.String()), nil
}

func NewTable(opts TableOptions) table.Writer {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(opts.OutputMirror)

	t.AppendHeader(opts.Header)

	if opts.Style != nil {
		t.SetStyle(*opts.Style)
	}

	if opts.ColumnConfig != nil {
		t.SetColumnConfigs(*opts.ColumnConfig)
	}

	if opts.Title != nil {
		t.SetTitle(*opts.Title)
	}

	if opts.Footer != nil {
		t.AppendFooter(*opts.Footer)
	}

	if opts.SortField != "" {
		t.SortBy([]table.SortBy{{
			Name: opts.SortField,
			Mode: func() table.SortMode {
				if opts.SortOrder < 0 {
					return opts.SortOrder
				}
				return table.Dsc
			}(),
		}})
	}

	return t
}
