package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		// id UUID PRIMARY KEY DEFAULT gen_random_uuid()
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id"),

		// title TEXT NOT NULL
		field.Text("title").
			NotEmpty(),

		// description TEXT
		field.Text("description").
			Optional(),

		// completed BOOLEAN NOT NULL DEFAULT FALSE
		field.Bool("completed").
			Default(false),

		// category_id UUID REFERENCES categories(id) ON DELETE SET NULL
		field.UUID("category_id", uuid.UUID{}).
			Optional().
			Nillable(),

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

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return []ent.Edge{
		// category_id UUID REFERENCES categories(id) ON DELETE SET NULL
		edge.To("category", Category.Type).
			Unique().
			Field("category_id"),
	}
}
