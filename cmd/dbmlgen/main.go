package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	dbmlgen "github.com/alonelucky/gen"
	"gorm.io/gen"
)

var dbfile = flag.String("dbml", "", "database mark languge")
var output = flag.String("output", "./gen/curd", "curd and model output dir")
var types = flag.String("type", "datetime=string", "specify data type mapping relationship")
var nullable = flag.Bool("null", false, "generate model global configuration")
var cwd, _ = os.Getwd()

func main() {
	flag.Parse()
	if dbfile == nil {
		panic("dbml file will must.")
	}

	f, e := os.Open(*dbfile)
	if e != nil {
		panic(e)
	}

	if !filepath.IsAbs(*output) {
		*output = filepath.Join(cwd, *output)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:          *output,
		FieldNullable:    *nullable,
		FieldWithTypeTag: true,
		Mode:             gen.WithDefaultQuery,
	})

	var opts []dbmlgen.Option
	for _, v := range strings.Split(*types, ",") {
		keys := strings.Split(v, "=")
		if len(keys) == 2 {
			opts = append(opts, dbmlgen.WithType(keys[0], keys[1]))
		}
	}
	dbml := dbmlgen.NewDBML(f, g, opts...)
	if dbml.Error != nil {
		panic(dbml.Error)
	}
	g.ApplyBasic(dbml.All()...)
	g.Execute()
}
