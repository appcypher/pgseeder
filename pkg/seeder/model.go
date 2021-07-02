package seeder

import (
	"time"

	"github.com/gofrs/uuid"
)

// PgseederSeed is the pgseeder_seeds table model.
type PgseederSeed struct {
	ID         uuid.UUID `pg:",pk, unique, notnull, type:uuid, default:uuid_generate_v4()"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	SeedKeys   []SeedKey `pg:"seed_keys"`
	TableName  string
	TableOrder int
}

// SeedKey represnts primary/composite key info of a seed.
type SeedKey struct {
	Name  string `pg:"name"`
	Value string `pg:"value"`
	Type  string `pg:"type"`
}

// SeedKeyType is the type of a key.
type SeedKeyType = string

// ...
const (
	SeedKeyTypeString  SeedKeyType = "SeedKeyTypeString"
	SeedKeyTypeInteger SeedKeyType = "SeedKeyTypeInteger"
)
