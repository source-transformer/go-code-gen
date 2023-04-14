// Copyright 2021-present Drop Fake Inc. All rights reserved.

package csharp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CSharpTestSuite struct {
	suite.Suite
	//VariableThatShouldStartAtFive int
}

type FooBar struct {
	text string
}

func (s *CSharpTestSuite) TestSimpleCSharpCodegen() {
	fooBar := &FooBar{}
	csharpStr, err := GenerateCSharpFromInstance(fooBar)
	fmt.Println("csharpStr:", csharpStr, "err:", err)
	expected := `
public class FooBar
{
	public string text { get; set; }
}
`
	s.Assertions.Equal(expected, csharpStr)
}

type ItemPresentable struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty" jsonschema:"example=icon-ability-melee-1"`
}

type ItemWithUserDefinedType struct {
	Presentable *ItemPresentable `json:"presentable,omitempty"`
}

type ItemWithUserDefinedTypeWithSkipField struct {
	ID          string           `json:"id" jsonschema:"title=ID"`
	Variant     string           `json:"variant"`
	Category    string           `json:"category"`
	Name        string           `json:"name"`
	Presentable *ItemPresentable `json:"presentable,omitempty"`
}

func (s *CSharpTestSuite) TestGenerateCSharp() {
	testInstance2 := ItemWithUserDefinedType{}
	testInstance2.Presentable = &ItemPresentable{}
	testInstance2.Presentable.Title = "asdf"
	testInstance2.Presentable.Description = "foo"
	csharpStr, err := GenerateCSharpFromInstance(testInstance2)
	s.Assertions.NoError(err)
	expectedCSharpStr := `
	public class ItemPresentable
	{
		public string Title { get; set; }
		public string Description { get; set; }
		public string Icon { get; set; }
	}

	public class ItemWithUserDefinedType
	{
		public ItemPresentable Presentable { get; set; }
	}
	`
	expectedCSharpStr = strings.ReplaceAll(expectedCSharpStr, "\t", "")
	csharpStr = strings.ReplaceAll(csharpStr, "\t", "")
	s.Assertions.Equal(expectedCSharpStr, csharpStr)
	s.Assertions.NoError(err)
}

func (s *CSharpTestSuite) TestGenerateCSharpWithSkipField() {
	testInstance := ItemWithUserDefinedTypeWithSkipField{}
	schemaForCSharp, err := GenerateCSharpFromInstance(&testInstance)
	s.Assertions.NoError(err)
	expectedCSharpStr := `
	public class ItemWithUserDefinedTypeWithSkipField
	{
		public string ID { get; set; }
		public string Variant { get; set; }
		public string Category { get; set; }
		public string Name { get; set; }
		public ItemPresentable Presentable { get; set; }
	}
	`
	expectedCSharpStr = strings.ReplaceAll(expectedCSharpStr, "\t", "")
	schemaForCSharp = strings.ReplaceAll(schemaForCSharp, "\t", "")
	s.Assertions.Equal(expectedCSharpStr, schemaForCSharp)
	s.Assertions.NoError(err)
}

func (s *CSharpTestSuite) TestGenerateCSharpWithNilField() {
	testInstance2 := ItemWithUserDefinedType{}
	testInstance2.Presentable = nil
	csharpStr, err := GenerateCSharpFromInstance(testInstance2)
	s.Assertions.NoError(err)
	expectedCSharpStr := `
	public class ItemWithUserDefinedType
	{
		public ItemPresentable Presentable { get; set; }
	}
	`
	expectedCSharpStr = strings.ReplaceAll(expectedCSharpStr, "\t", "")
	csharpStr = strings.ReplaceAll(csharpStr, "\t", "")
	s.Assertions.Equal(expectedCSharpStr, csharpStr)
	s.Assertions.NoError(err)
}

type WeightedVariation struct {
	ID     string `json:"id"`
	Weight int    `json:"weight"`
}

type StructWithArray struct {
	Variations []*WeightedVariation `json:"variations,omitempty"`
}

func (s *CSharpTestSuite) TestGenerateCSharpWithArray() {
	testInstance := StructWithArray{}
	weightedVariation := &WeightedVariation{ID: "asdf", Weight: 1}
	testInstance.Variations = []*WeightedVariation{weightedVariation}
	csharpStr, err := GenerateCSharpFromInstance(testInstance)
	s.Assertions.NoError(err)
	expectedCSharpStr := `
	public class WeightedVariation
	{
		public string ID { get; set; }
		public int Weight { get; set; }
	}

	public class StructWithArray
	{
		public WeightedVariation[] Variations { get; set; }
	}
	`
	expectedCSharpStr = strings.ReplaceAll(expectedCSharpStr, "\t", "")
	csharpStr = strings.ReplaceAll(csharpStr, "\t", "")
	s.Assertions.Equal(expectedCSharpStr, csharpStr)
	s.Assertions.NoError(err)
}

type AnotherDataItemWithUserDefinedType struct {
	Group string `json:"group"`
	Name  string `json:"name"`

	Enabled      bool                 `json:"enabled,omitempty"`
	NewUsersOnly bool                 `json:"newUsersOnly,omitempty"`
	Variations   []*WeightedVariation `json:"variations,omitempty"`
}

func (s *CSharpTestSuite) TestGenerateCSharpWithArray2() {
	testInstance := AnotherDataItemWithUserDefinedType{}
	weightedVariation := &WeightedVariation{ID: "asdf", Weight: 1}
	testInstance.Variations = []*WeightedVariation{weightedVariation}
	csharpStr, err := GenerateCSharpFromInstance(testInstance)
	s.Assertions.NoError(err)
	expectedCSharpStr := `
	public class WeightedVariation
	{
		public string ID { get; set; }
		public int Weight { get; set; }
	}

	public class AnotherDataItemWithUserDefinedType
	{
		public string Group { get; set; }
		public string Name { get; set; }
		public bool Enabled { get; set; }
		public bool NewUsersOnly { get; set; }
		public WeightedVariation[] Variations { get; set; }
	}
	`
	expectedCSharpStr = strings.ReplaceAll(expectedCSharpStr, "\t", "")
	csharpStr = strings.ReplaceAll(csharpStr, "\t", "")
	s.Assertions.Equal(expectedCSharpStr, csharpStr)
	s.Assertions.NoError(err)
}

type ManyBasicTypes struct {
	Scale32       float32 `json:"scale32"`
	Scale64       float64 `json:"scale64"`
	SignedInt32   int32   `json:"signedInt32"`
	SignedInt64   int64   `json:"signedInt64"`
	UnsignedInt32 uint32  `json:"unsignedInt32"`
	UnsignedInt64 uint64  `json:"unsignedInt64"`
	Boolean       bool    `json:"boolean"`
	Text          string  `json:"text"`
	Text2         string  `json:"fromJsonText2"`
}

type AnotherType struct {
	SomeField  string `json:"someField"`
	SomeField2 string `json:"fromJsonSomeField2"`
}
type UserDefinedType struct {
	Field1        ManyBasicTypes   `json:"field1"`
	Array1        []ManyBasicTypes `json:"array1"`
	AnotherField  AnotherType      `json:"anotherField"`
	AnotherField2 AnotherType      `json:"fromJsonAnotherField2"`
}

func (s *CSharpTestSuite) TestCustomCSharpCodegen() {
	testInstance := UserDefinedType{}
	csharpStr, err := GenerateCSharpFromInstance(testInstance)
	s.Assertions.NoError(err)
	expectedCSharpStr := `
	public class ManyBasicTypes
	{
		public float Scale32 { get; set; }
		public float Scale64 { get; set; }
		public int SignedInt32 { get; set; }
		public int SignedInt64 { get; set; }
		public uint UnsignedInt32 { get; set; }
		public uint UnsignedInt64 { get; set; }
		public bool Boolean { get; set; }
		public string Text { get; set; }
		public string Text2 { get; set; }
	}

	public class AnotherType
	{
		public string SomeField { get; set; }
		public string SomeField2 { get; set; }
	}

	public class UserDefinedType
	{
		public ManyBasicTypes Field1 { get; set; }
		public ManyBasicTypes[] Array1 { get; set; }
		public AnotherType AnotherField { get; set; }
		public AnotherType AnotherField2 { get; set; }
	}
	`
	expectedCSharpStr = strings.ReplaceAll(expectedCSharpStr, "\t", "")
	csharpStr = strings.ReplaceAll(csharpStr, "\t", "")
	s.Assertions.Equal(expectedCSharpStr, csharpStr)
}

func TestCSharpTestSuite(t *testing.T) {
	suite.Run(t, new(CSharpTestSuite))
}
