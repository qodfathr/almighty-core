package models

import (
	"testing"

	"github.com/almighty/almighty-core/app"
	"github.com/almighty/almighty-core/resource"
	"github.com/stretchr/testify/assert"
)

func TestCompatibleFields(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)
	a := FieldDefinition{
		Required: true,
		Type: ListType{
			SimpleType:    SimpleType{Kind: KindList},
			ComponentType: SimpleType{Kind: KindString},
		},
	}
	b := FieldDefinition{
		Required: true,
		Type: ListType{
			SimpleType:    SimpleType{Kind: KindList},
			ComponentType: SimpleType{Kind: KindString},
		},
	}
	assert.True(t, compatibleFields(a, b))
}

func TestConvertTypeFromModels(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)

	//------------------------------
	// Work item type in model space
	//------------------------------

	a := WorkItemType{
		Name:       "foo",
		Version:    42,
		ParentPath: "something",
		Fields: map[string]FieldDefinition{
			"aListType": FieldDefinition{
				Type: EnumType{
					BaseType:   SimpleType{KindString},
					SimpleType: SimpleType{KindEnum},
					Values:     []interface{}{"open", "done", "closed"},
				},
				Required: true,
			},
		},
	}

	//----------------------------
	// Work item type in app space
	//----------------------------

	// Create an enumeration of animal names
	typeStrings := []string{"open", "done", "closed"}

	// Convert string slice to slice of interface{} in O(n) time.
	typeEnum := make([]interface{}, len(typeStrings))
	for i := range typeStrings {
		typeEnum[i] = typeStrings[i]
	}

	// Create the type for "animal-type" field based on the enum above
	stString := "string"
	expected := app.WorkItemType{
		Name:    "foo",
		Version: 42,
		Fields: map[string]*app.FieldDefinition{
			"aListType": &app.FieldDefinition{
				Required: true,
				Type: &app.FieldType{
					BaseType: &stString,
					Kind:     "enum",
					Values:   typeEnum,
				},
			},
		},
	}

	result := convertTypeFromModels(&a)

	assert.Equal(t, expected.Name, result.Name)
	assert.Equal(t, expected.Version, result.Version)
	assert.Len(t, result.Fields, len(expected.Fields))
	assert.Equal(t, expected.Fields, result.Fields)
}

func TestConvertAnyToKind(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)

	_, err := convertAnyToKind(1234)
	assert.NotNil(t, err)
}

func TestConvertStringToKind(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)

	_, err := convertStringToKind("DefinitivelyNotAKind")
	assert.NotNil(t, err)
}

func TestConvertFieldTypeToModels(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)

	// Create an enumeration of animal names
	typeStrings := []string{"open", "done", "closed"}

	// Convert string slice to slice of interface{} in O(n) time.
	typeEnum := make([]interface{}, len(typeStrings))
	for i := range typeStrings {
		typeEnum[i] = typeStrings[i]
	}

	// Create the type for "animal-type" field based on the enum above
	stString := "string"

	_ = &app.FieldType{
		BaseType: &stString,
		Kind:     "DefinitivelyNotAType",
		Values:   typeEnum,
	}
	_, err := convertFieldTypeToModels(app.FieldType{Kind: "DefinitivelyNotAType"})
	assert.NotNil(t, err)
}
