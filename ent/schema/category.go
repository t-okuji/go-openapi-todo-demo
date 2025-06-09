package schema

import (
	"regexp"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		// id UUID PRIMARY KEY DEFAULT gen_random_uuid()
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id"),

		// name VARCHAR(50) NOT NULL CHECK (LENGTH(name) >= 1)
		field.String("name").
			MaxLen(50).
			NotEmpty(),

		// description VARCHAR(255)
		field.String("description").
			MaxLen(255).
			Optional(),

		// color VARCHAR(7) CHECK (color ~ '^#[0-9A-Fa-f]{6}$') DEFAULT '#6c757d'
		field.String("color").
			MaxLen(7).
			Default("#6c757d").
			Match(regexp.MustCompile("^#[0-9A-Fa-f]{6}$")),

		// created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		field.Time("created_at").
			Default(time.Now).
			Immutable(),

		// updated_at TIMESTAMP WITH TIME ZONE
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		// One-to-many relationship with todos
		edge.From("todos", Todo.Type).
			Ref("category"),
	}
}