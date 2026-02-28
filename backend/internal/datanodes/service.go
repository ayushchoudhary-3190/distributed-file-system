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

func (dns *datanodeserver) WriteChunk(ctx context.Context, req *pb.ChunkWriteRequest) (*pb.ChunkWriteResponse, error) {

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

// getChunkAddress is a package-level helper function that returns array of node endpoints for a chunk
// Returns []*pb.DataNodeEndpoint containing node_id and address
// Accepts *gorm.DB as parameter to query the database
func GetChunkAddress(db *gorm.DB, chunkID string) []*pb.DataNodeEndpoint { //// helper function
	var endpoints []*pb.DataNodeEndpoint

	// Step 1: Query Chunk_table to find which nodes have this chunk
	var chunk Chunk_table
	result := db.Where("chunk_id = ?", chunkID).First(&chunk)
	if result.Error != nil {
		// If chunk not found, return empty endpoints
		return endpoints
	}

	// Step 2: For each nodeID that has this chunk, query Node_table to get address
	for _, nodeID := range chunk.NodeID {
		var node Node_table
		nodeResult := db.Where("node_id = ?", nodeID).First(&node)
		if nodeResult.Error != nil {
			continue // Skip if node not found
		}

		// Step 3: Add to endpoints array
		endpoints = append(endpoints, &pb.DataNodeEndpoint{
			NodeId:  node.NodeID,
			Address: node.BaseDir,
		})
	}

	return endpoints
}
