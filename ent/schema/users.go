package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Users holds the schema definition for the Users entity.
type Users struct {
	ent.Schema
}

// Fields of the Users.
func (Users) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("name"),
		field.String("email"),
	}
}

// Edges of the Users.
func (Users) Edges() []ent.Edge {
	return nil
}
