package seeder

// Flags is commandline flags
type Flags struct {
	AddAll             bool
	RemoveAll          bool
	SeedToAdd          string
	SeedToRemove       string
	FileToCreate       string
	Directory          string
	ConnectionURI      string
	SeedKeyNames       string
	SeedKeyIntegerType bool
	Verbose            bool
}
