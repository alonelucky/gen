package main

import (
	"os"
	"testing"

	dbmlgen "github.com/alonelucky/gen"
	"gorm.io/gen"
)

func TestXXX(t *testing.T) {
	f, e := os.Open("./tables.dbml")
	if e != nil {
		t.Fatal(e)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./curd",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	var opts []dbmlgen.Option
	opts = append(opts, dbmlgen.WithType("datetime", "string"))
	dbml := dbmlgen.NewDBML(f, g, opts...)
	if dbml.Error != nil {
		t.Fatal(dbml.Error)
	}
	g.ApplyBasic(dbml.All()...)
	g.Execute()
}
