package model

import "time"

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity uint64
	Category      Category
	Dimension     *Dimension
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Metadata
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type Filter struct {
	UUIDs                 map[string]struct{}
	Names                 map[string]struct{}
	Tags                  map[string]struct{}
	Categories            map[Category]struct{}
	ManufacturerCountries map[string]struct{}
}

type Category uint8

const (
	CategoryUndefined Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Dimension struct {
	Width  float64
	Height float64
	Length float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Metadata struct {
	Kind   MetadataKind
	String *string
	Double *float64
	Bool   *bool
	Int64  *int64
}

type MetadataKind string

const (
	MetadataKindString MetadataKind = "string"
	MetadataKindInt64  MetadataKind = "int64"
	MetadataKindDouble MetadataKind = "double"
	MetadataKindBool   MetadataKind = "bool"
)
