package datanodeservice

type Config struct {
	NodeID      string `gorm:"-"`
	NodeAddress string `gorm:"-"`
}

type Chunk_table struct {
	ChunkID string   `gorm:"primaryKey"`
	NodeID  []string `gorm:"type:text[];not null"`
	Index   int64
}

type Node_table struct {
	NodeID        string
	LastHeartbeat int64
}
