package metaservice


type File_table struct{
	FileID	  int64 `gorm:"primaryKey;autoIncrement"`
	FileName  string `gorm:"not null"`
	OwnerID	  string `gorm:"unique;not null"`
	ChunkCount int64  `gorm:"not null"`
	ChunkArray []string `gorm:"type:text[];not null"`
	FileSize    int64   `gorm:"not null"`
}


