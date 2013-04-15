package types

func init() {
	Register(&singleValue{
		name:        "integer",
		label:       "Integer",
		description: "A single integer",
		template:    "textbox",
		fieldType:   FieldInteger,
	})

	Register(&singleValue{
		name:        "string",
		label:       "Text",
		description: "A single string value",
		fieldType:   FieldString,
		template:    "textbox",
	})

	// Register(&singleValue{
	// 	name:        "float",
	// 	label:       "Floating Point",
	// 	description: "A single floating Point entry",
	// 	template:    "textbox",
	// })
}
