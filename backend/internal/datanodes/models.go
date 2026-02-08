package datanodeservice


type Chunk_table struct{
	ChunkID string `gorm:"primaryKey"`
	NodeID  []string `gorm:"type:text[];not null"`
	Index	int64  
}

type Node_table struct{
	NodeID string
	BaseDir string
}
