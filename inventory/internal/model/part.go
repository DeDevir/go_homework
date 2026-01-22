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
	Metadata      map[string]*MetadataValue
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type PartsFilter struct {
	Uuid                  []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Dimension struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type MetadataKind string

const (
	MetadataKindString MetadataKind = "string"
	MetadataKindInt64  MetadataKind = "int64"
	MetadataKindDouble MetadataKind = "double"
	MetadataKindBool   MetadataKind = "bool"
)

type MetadataValue struct {
	Kind   MetadataKind
	String *string
	Int64  *int64
	Double *float64
	Bool   *bool
}

//type Tag string

type Category uint8

const (
	CategoryUNKNOWN Category = iota
	CategoryENGINE
	CategoryFuel
	CategoryPorthole
	CategoryWing
)
