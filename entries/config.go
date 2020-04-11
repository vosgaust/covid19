package entries

// Config defines configuration fields for mysql storage
type Config struct {
	Connection string `default:"user:password@/database"`
	Table      string `default:"entries"`
}
