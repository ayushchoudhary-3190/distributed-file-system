package datanodeservice

import (
	"context"
	"fmt"
	"github.com/ayushchoudhary-3190/Distributed_file_system/pb"
	"gorm.io/gorm"
)

type datanodeserver struct {
	DB *gorm.DB
	pb.UnimplementedDataNodeServiceServer
}

func (dns *datanodeserver) ReadChunks(ctx context.Context, req *pb.ChunkReadRequest) (*pb.ChunkReadResponse, error) {
	// Handle edge case: empty chunk_id
	if req.ChunkId == "" {
		response := &pb.ChunkReadResponse{
			Data: nil,
			Eof:  false,
		}
		return response, fmt.Errorf("chunk_id cannot be empty")
	}

	// Read chunk data from storage (placeholder - implement actual storage logic)
	chunkData, err := dns.readChunkFromStorage(req.ChunkId)
	if err != nil {
		response := &pb.ChunkReadResponse{
			Data: nil,
			Eof:  false,
		}
		return response, fmt.Errorf("failed to read chunk: %v", err)
	}

	// Handle edge case: empty file/chunk
	if len(chunkData) == 0 {
		response := &pb.ChunkReadResponse{
			Data: nil,
			Eof:  true,
		}
		return response, nil
	}

	// Calculate actual data to read
	chunkSize := uint64(len(chunkData))
	var actualOffset, actualLength uint64

	// Ensure offset is within bounds
	if req.Offset >= chunkSize {
		response := &pb.ChunkReadResponse{
			Data: nil,
			Eof:  true,
		}
		return response, nil
	}
	actualOffset = req.Offset

	// Check if this is the last chunk read (should go to EOF regardless of length)
	if actualOffset+req.Length >= chunkSize {
		// Last increment: read from offset to EOF, ignore length
		actualLength = chunkSize - actualOffset
		return &pb.ChunkReadResponse{
			Data: chunkData[actualOffset : actualOffset+actualLength],
			Eof:  true,
		}, nil
	} else {
		// Normal chunk read: use specified length
		actualLength = req.Length
		return &pb.ChunkReadResponse{
			Data: chunkData[actualOffset : actualOffset+actualLength],
			Eof:  false,
		}, nil
	}
}

// readChunkFromStorage is a helper function to read chunk data from storage (placeholder)
func (dns *datanodeserver) readChunkFromStorage(chunkID string) ([]byte, error) {
	// TODO: Implement actual storage reading logic
	// This should read chunk data from disk, database, or other storage
	// For now, return sample data
	return []byte("sample chunk data for " + chunkID), nil
}
