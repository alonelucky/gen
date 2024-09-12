// Package gen as dbml
package gen

import (
	"bytes"
	"io"
	"strings"
	"unicode"

	"github.com/duythinht/dbml-go/core"
	"github.com/duythinht/dbml-go/parser"
	"github.com/duythinht/dbml-go/scanner"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/helper"
)

var types = make(map[string]string)
var enums = make(map[string]core.Enum)

type Option func(*dbml)

func WithType(coltype, gotyp string) Option {
	return func(d *dbml) {
		types[coltype] = gotyp
	}
}

type dbml struct {
	g     *gen.Generator
	dbml  *core.DBML
	Error error
}

// NewDBML as
func NewDBML(r io.Reader, g *gen.Generator, opts ...Option) *dbml {
	var sml dbml
	s := scanner.NewScanner(r)
	parser := parser.NewParser(s)
	sml.dbml, sml.Error = parser.Parse()
	sml.g = g
	for _, v := range opts {
		v(&sml)
	}
	return &sml
}

func (s *dbml) All() (lst []interface{}) {
	for _, v := range s.dbml.Enums {
		enums[v.Name] = v
	}
	for _, v := range s.dbml.Tables {
		mf := s.g.GenerateModelFrom(NewDBMLObject(v))
		lst = append(lst, mf)
	}
	return
}

type dbmlobject struct {
	core.Table
}

// NewDBMLObject as
func NewDBMLObject(model core.Table) helper.Object {
	var s dbmlobject
	s.Table = model
	return &s
}

// TableName return table name
func (s *dbmlobject) TableName() string {
	return s.Table.Name
}

// StructName return struct name
func (s *dbmlobject) StructName() string {
	if s.Table.As != "" {
		return s.Table.As
	}
	return CamelName(s.Table.Name)
}

// FileName return field name
func (s *dbmlobject) FileName() string {
	return UnderscoreName(s.Table.Name)
}

// ImportPkgPaths return need import package path
func (s *dbmlobject) ImportPkgPaths() []string {
	return nil
}

// Fields return field array
func (s *dbmlobject) Fields() (lst []helper.Field) {
	for _, v := range s.Columns {
		lst = append(lst, NewDBMLField(v))
	}
	return
}

type dbmlfield struct {
	core.Column
}

// NewDBMLField as
func NewDBMLField(model core.Column) helper.Field {
	var s dbmlfield
	s.Column = model
	return &s
}

// Name return field name
func (s *dbmlfield) Name() string {
	return CamelName(s.Column.Name)
}

// Type return field type
func (s *dbmlfield) Type() (typ string) {
	if v := enums[s.Column.Type]; v.Name != "" {
		return "string"
	}

	if typ = types[s.Column.Type]; typ != "" {
		return
	}
	if s.Column.Type == "bigint" {
		return "int64"
	}

	if strings.Contains(s.Column.Type, "int") || strings.Contains(s.Column.Type, "year") {
		return "int32"
	}

	if strings.Contains(s.Column.Type, "bool") {
		return "bool"
	}

	if strings.Contains(s.Column.Type, "decimal") || strings.Contains(s.Column.Type, "float") {
		return "float64"
	}

	if strings.Contains(s.Column.Type, "text") || strings.Contains(s.Column.Type, "char") || strings.Contains(s.Column.Type, "json") {
		return "string"
	}

	if strings.Contains(s.Column.Type, "blob") {
		return "[]byte"
	}

	if strings.Contains(s.Column.Type, "time") || strings.Contains(s.Column.Type, "date") {
		return "time.Time"
	}

	return "interface{}"
}

// ColumnName return column name
func (s *dbmlfield) ColumnName() string {
	return s.Column.Name
}

// GORMTag return gorm tag
func (s *dbmlfield) GORMTag() string {
	var typ = []string{`type:` + s.Column.Type}
	if v := enums[s.Column.Type]; v.Name != "" {
		typ = []string{`type:varchar`}
	}
	if s.Column.Settings.PK {
		typ = append(typ, "primaryKey")
	}

	if s.Column.Settings.Default != "" {
		typ = append(typ, "default:"+s.Column.Settings.Default)
	}

	if s.Column.Settings.Note != "" {
		typ = append(typ, "comment:"+s.Column.Settings.Note)
	}

	if s.Column.Settings.Increment {
		typ = append(typ, "autoincrement")
	}

	if !s.Column.Settings.Null {
		typ = append(typ, "not null")
	}

	if s.Column.Settings.Unique {
		typ = append(typ, "unique")
	}
	return strings.Join(typ, ";")
}

// JSONTag return json tag
func (s *dbmlfield) JSONTag() string {
	return s.Column.Name
}

// Tag return field tag
func (s *dbmlfield) Tag() field.Tag {
	return make(field.Tag)
}

// Comment return comment
func (s *dbmlfield) Comment() string {
	return s.Column.Settings.Note
}

const space = ' '

// CamelName 下划线写法转为驼峰写法
func CamelName(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	prev := ' '
	names := []rune(name)
	for k := range names {
		if prev == space {
			prev = names[k]
			names[k] -= 32
			continue
		}
		prev = names[k]
	}
	return strings.Replace(string(names), " ", "", -1)
}

// UnderscoreName as
func UnderscoreName(name string) string {
	var buf bytes.Buffer
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buf.WriteByte('_')
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
