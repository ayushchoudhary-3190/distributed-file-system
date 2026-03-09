package datanodeservice

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ayushchoudhary-3190/Distributed_file_system/pb"
	"gorm.io/gorm"
)

type datanodeserver struct {
	DB          *gorm.DB
	NodeID      string
	NodeAddress string
	pb.UnimplementedDataNodeServiceServer
}

type nodeserver struct {
	DB *gorm.DB
	pb.UnimplementedNodeControlServiceServer
}

func NewDataNodeServer(db *gorm.DB) (*datanodeserver, error) {
	nodeID := os.Getenv("NODE_ID")
	nodeAddress := os.Getenv("NODE_ADDRESS")

	if nodeID == "" {
		return nil, fmt.Errorf("NODE_ID environment variable not set")
	}
	if nodeAddress == "" {
		return nil, fmt.Errorf("NODE_ADDRESS environment variable not set")
	}

	return &datanodeserver{
		DB:          db,
		NodeID:      nodeID,
		NodeAddress: nodeAddress,
	}, nil
}

func (dns *datanodeserver) WriteChunk(ctx context.Context, req *pb.ChunkWriteRequest) (*pb.ChunkWriteResponse, error) {
	if req.ChunkId == "" {
		return &pb.ChunkWriteResponse{Ok: false}, nil
	}

	if len(req.Data) == 0 {
		return &pb.ChunkWriteResponse{Ok: false}, nil
	}

	filePath := filepath.Join(dns.NodeAddress, req.ChunkId+".bin")

	err := os.WriteFile(filePath, req.Data, 0600)
	if err != nil {
		return &pb.ChunkWriteResponse{Ok: false}, nil
	}

	return &pb.ChunkWriteResponse{Ok: true}, nil
}

func (dns *datanodeserver) ReadChunk(ctx context.Context, req *pb.ChunkReadRequest) (*pb.ChunkReadResponse, error) {
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

// readChunkFromStorage reads chunk binary data from local storage
// Path format: {NodeAddress}/{chunkID}.bin
func (dns *datanodeserver) readChunkFromStorage(chunkID string) ([]byte, error) {
	if dns.NodeAddress == "" {
		return nil, fmt.Errorf("node address not configured")
	}

	filePath := filepath.Join(dns.NodeAddress, chunkID+".bin")

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("chunk %s not found at %s", chunkID, filePath)
		}
		return nil, fmt.Errorf("failed to read chunk %s: %w", chunkID, err)
	}

	return data, nil
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

// Heartbeat is the RPC implementation for NodeControlService
// It checks if a node is active by verifying node_id and address, and updates heartbeat timestamp
func (ns *nodeserver) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	// Step 1: Validate request
	if req.NodeId == "" {
		return &pb.HeartbeatResponse{
			Status:  false,
			Message: "node_id cannot be empty",
		}, nil
	}

	// Step 2: Check if node exists in Node_table
	var node Node_table
	result := ns.DB.Where("node_id = ?", req.NodeId).First(&node)
	if result.Error != nil {
		return &pb.HeartbeatResponse{
			Status:  false,
			Message: "node not found",
		}, nil
	}

	// Step 3: Verify address matches
	if node.BaseDir != req.Address {
		return &pb.HeartbeatResponse{
			Status:  false,
			Message: "address mismatch",
		}, nil
	}

	// Step 4: Update last heartbeat timestamp
	node.LastHeartbeat = time.Now().Unix()
	ns.DB.Save(&node)

	// Step 5: Return success
	return &pb.HeartbeatResponse{
		Status:  true,
		Message: "node is active",
	}, nil
}
