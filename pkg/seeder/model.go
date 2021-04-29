package seeder

import (
	"time"

	"github.com/gofrs/uuid"
)

// PgseederSeed is the pgseeder_seeds table model.
type PgseederSeed struct {
	ID          uuid.UUID `pg:",pk, unique, notnull, type:uuid, default:uuid_generate_v4()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	SeedPKName  string `pg:"seed_pk_name"`
	SeedPKValue string `pg:"seed_pk_value"`
	TableName   string
}
