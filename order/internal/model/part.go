package model

import "github.com/google/uuid"

type PartFilter struct {
	UUIDS                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Part struct {
	Uuid          uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity uint64
	Category      Category
	Dimension     *Dimension
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Metadata
}

type Category uint64

const (
	CategoryUnknown Category = iota
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
	Kind  MetadataKind
	Value MetadataValue
}

type MetadataKind string

const (
	MetadataKindString MetadataKind = "string_value"
	MetadataKindInt64  MetadataKind = "int64_value"
	MetadataKindDouble MetadataKind = "double_value"
	MetadataKindBool   MetadataKind = "bool_value"
)

type MetadataValue struct {
	String *string
	Double *float64
	Bool   *bool
	Int    *int64
}
