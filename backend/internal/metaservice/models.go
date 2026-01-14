package metaservice


type File_table struct{
	FileID	  int64 `gorm:"primaryKey;autoIncrement"`
	FileName  string `gorm:"not null"`
	OwnerID	  string `gorm:"unique;not null"`
	ChunkCount int64  `gorm:"not null"`
	ChunkArray []string `gorm:"type:text[];not null"`
	FileSize    int64   `gorm:"not null"`
}

type Chunk_table struct{
	ChunkID string `gorm:"primaryKey"`
	NodeID  []string `gorm:"type:text[];not null"`
	Index	int64  
}

type Node_table struct{
	NodeID string
	BaseDir string
}

