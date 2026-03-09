package metaservice

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	datanodeservice "github.com/ayushchoudhary-3190/Distributed_file_system/internal/datanodes"
	pb "github.com/ayushchoudhary-3190/Distributed_file_system/pb"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetaServer struct {
	pb.UnimplementedMetaServiceServer
	DB                *gorm.DB
	dataNodeClient    pb.DataNodeServiceClient
	nodeControlClient pb.NodeControlServiceClient
}

// ChunkDataWithIndex holds chunk data along with its index for ordered reconstruction
type ChunkDataWithIndex struct {
	Index int32
	Data  []byte
}

// function to add a new file to the metaservice table
func (s *MetaServer) UploadRequest(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	//insert file metadata in metadata table
	tx := s.DB.Begin()

	// Extract chunk IDs from request
	chunkArray := make([]string, len(req.Chunks))
	for i, chunk := range req.Chunks {
		chunkArray[i] = chunk.Chunkid
	}

	file := File_table{
		FileID:     req.Fileid,
		FileName:   req.Filename,
		OwnerID:    req.Ownerid,
		ChunkCount: req.Chunkcount,
		ChunkArray: chunkArray,
		FileSize:   req.Filesize,
	}

	// Create file record in database
	if err := tx.Create(&file).Error; err != nil {
		tx.Rollback()
		response := &pb.UploadFileResponse{
			Path:     req.Filename,
			Response: "Failed to upload file metadata",
		}
		return response, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		response := &pb.UploadFileResponse{
			Path:     req.Filename,
			Response: "Failed to commit transaction",
		}
		return response, err
	}

	// Return success response
	response := &pb.UploadFileResponse{
		Path:     req.Filename,
		Response: "File uploaded successfully",
	}
	return response, nil
}

// function to delete a file from the metaservice table
func (s *MetaServer) DeleteRequest(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	// Start transaction
	tx := s.DB.Begin()

	// Find and delete the file with the given path
	result := tx.Where("file_name = ?", req.Path).Delete(&File_table{})

	if result.Error != nil {
		tx.Rollback()
		response := &pb.DeleteFileResponse{
			Message: "Failed to delete file",
			Error:   result.Error.Error(),
		}
		return response, result.Error
	}

	// Check if any record was actually deleted
	if result.RowsAffected == 0 {
		tx.Rollback()
		response := &pb.DeleteFileResponse{
			Message: "File not found",
			Error:   "file not deleted",
		}
		return response, nil
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		response := &pb.DeleteFileResponse{
			Message: "Failed to commit transaction",
			Error:   err.Error(),
		}
		return response, err
	}

	// Return success response
	response := &pb.DeleteFileResponse{
		Message: "File deleted successfully",
		Error:   "",
	}
	return response, nil
}

// function to list files belonging to a specific owner
func (s *MetaServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	var files []metaservice.File_table

	// Query files by owner_id
	result := s.DB.Where("owner_id = ?", req.OwnerId).Find(&files)

	if result.Error != nil {
		response := &pb.ListFilesResponse{
			Owner:   req.Owner,
			OwnerId: req.OwnerId,
			Count:   0,
			Err:     result.Error.Error(),
		}
		return response, result.Error
	}

	// Create FileListItem array with only file names
	fileList := make([]*pb.FileListItem, len(files))
	for i, file := range files {
		fileList[i] = &pb.FileListItem{
			Filename: file.FileName,
		}
	}

	// Return success response
	response := &pb.ListFilesResponse{
		Filename: fileList,
		Owner:    req.Owner,
		OwnerId:  req.OwnerId,
		Count:    int64(len(fileList)),
		Err:      "",
	}
	return response, nil
}

// function to get file by owner_id and path and reconstruct from chunks using new workflow
func (s *MetaServer) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	var file metaservice.File_table

	// Query file by owner_id and file_name (path)
	result := s.DB.Where("owner_id = ? AND file_name = ?", req.OwnerId, req.Path).First(&file)

	if result.Error != nil {
		response := &pb.GetFileResponse{
			OwnerId: req.OwnerId,
			Path:    req.Path,
			Size:    0,
			File:    nil,
		}
		return response, result.Error
	}

	// Use fileID inside reconstructFileFromId function
	fileData, err := s.reconstructFileFromId(file.FileID)

	if err != nil {
		log.Fatal("failed to reconstruct file from chunk ids")
		response := &pb.GetFileResponse{
			OwnerId: req.OwnerId,
			Path:    req.Path,
			Size:    0,
			File:    nil,
		}
		return response, err
	}

	// Return success response with reconstructed file data
	response := &pb.GetFileResponse{
		OwnerId: file.OwnerID,
		Path:    req.Path,
		Size:    file.FileSize,
		File:    fileData,
	}
	return response, nil
}

// ChunkResult holds the result of reading a chunk with its index for ordering
type ChunkResult struct {
	Index int
	Data  []byte
	Error error
}

// reconstructFileFromId reconstructs file from chunk IDs using location-based approach
func (s *MetaServer) reconstructFileFromId(fileID string) ([]byte, error) { //// metaservice function
	// Step 1: Call getChunkLocation with fileID to get ChunkLocation array
	locationResp, err := s.GetChunksLocations(context.Background(), &pb.GetChunkLocationRequest{
		FileId: fileID,
	})
	if err != nil {
		return nil, err
	}

	// Channel to collect results from goroutines
	results := make(chan ChunkResult, len(locationResp.Locs))
	var wg sync.WaitGroup

	// Step 2: Read chunks in PARALLEL using goroutines
	for _, chunkLocation := range locationResp.Locs {
		wg.Add(1)
		go func(cl *pb.ChunkLocation) {
			defer wg.Done()

			// Try each replica until find active DataNode
			for _, replica := range cl.Replicas {
				// Check if DataNode is active using Heartbeat with 0.5s timeout
				ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				defer cancel()

				heartbeatResp, err := s.nodeControlClient.Heartbeat(ctx, &pb.HeartbeatRequest{
					NodeId:  replica.NodeId,
					Address: replica.Address,
				})

				if err == nil && heartbeatResp.Status {
					// DataNode is active - read chunk
					req := &pb.ChunkReadRequest{
						ChunkId: cl.ChunkId,
						Offset:  0,
						Length:  33554432, // 32 MB
					}

					resp, err := s.dataNodeClient.ReadChunk(context.Background(), req)
					if err == nil {
						// Send result with Index for ordering
						results <- ChunkResult{
							Index: int(cl.Index),
							Data:  resp.Data,
							Error: nil,
						}
						return
					}
				}
			}

			// If we get here, chunk couldn't be read from any replica
			results <- ChunkResult{
				Index: int(cl.Index),
				Data:  nil,
				Error: fmt.Errorf("failed to read chunk %s from any replica", cl.ChunkId),
			}
		}(chunkLocation)
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var chunkResults []ChunkResult
	for result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
		chunkResults = append(chunkResults, result)
	}

	// Step 3: Sort AFTER reading by Index to ensure correct order
	sort.Slice(chunkResults, func(i, j int) bool {
		return chunkResults[i].Index < chunkResults[j].Index
	})

	// Step 4: Concatenate in correct order
	fileData := []byte{}
	for _, result := range chunkResults {
		fileData = append(fileData, result.Data...)
	}

	return fileData, nil
}

// getChunksLocation is a gRPC function that returns chunk locations for a file
// 1. Takes file_id as string parameter
// 2. Scans metaservice table (File_table) for that file_id
// 3. Gets the ChunkArray containing all chunkIDs
// 4. For each chunkID, calls getChunkAddress to get node information
// 5. Returns structured data using ChunkLocation and DataNodeEndpoint from proto
func (s *MetaServer) GetChunksLocations(ctx context.Context, req *pb.GetChunkLocationRequest) (*pb.GetChunkLocationResponse, error) { //// metaservice gRPC function
	// Step 1: Scan metaservice table for file_id
	var file metaservice.File_table
	result := s.DB.Where("file_id = ?", req.FileId).First(&file)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &pb.GetChunkLocationResponse{
				Locs: nil,
			}, fmt.Errorf("file not found: %s", req.FileId)
		}
		return &pb.GetChunkLocationResponse{
			Locs: nil,
		}, result.Error
	}

	// Step 2: Get the chunk array containing all chunkIDs
	chunkIDs := file.ChunkArray

	// Step 3: For each chunkID, call getChunkAddress to get node information
	var chunkLocations []*pb.ChunkLocation

	for i, chunkID := range chunkIDs {
		// Get chunk addresses (node IDs and addresses) for this chunk
		// Using datanodeservice package to call getChunkAddress helper
		nodeEndpoints := datanodeservice.GetChunkAddress(s.DB, chunkID)

		// Step 4: Structure the data according to ChunkLocation proto
		chunkLocation := &pb.ChunkLocation{
			ChunkId:  chunkID,
			Index:    int32(i),
			Replicas: nodeEndpoints,
		}

		chunkLocations = append(chunkLocations, chunkLocation)
	}

	// Step 5: Return structured data
	return &pb.GetChunkLocationResponse{
		Locs: chunkLocations,
	}, nil
}

// AllocateChunks is a gRPC function that allocates chunks for a file
// It selects DataNodes based on capacity and returns chunk locations
func (s *MetaServer) AllocateChunks(ctx context.Context, req *pb.AllocateChunksRequest) (*pb.AllocateChunksResponse, error) { //// metaservice gRPC function
	// Step 1: Get available DataNodes from Node_table
	// We need to access the datanodes package to query Node_table
	// For now, use placeholder approach - will be implemented with actual Node_table access

	// Get available nodes (placeholder - query from database)
	// This should query the Node_table in datanodes package to get all available nodes
	// Then sort by available capacity

	var availableNodes []string
	var nodeAddresses map[string]string

	// Placeholder: we'll use a simple approach
	// In actual implementation, query Node_table from datanodes database
	availableNodes = []string{"node1", "node2", "node3"}
	nodeAddresses = map[string]string{
		"node1": "localhost:50051",
		"node2": "localhost:50052",
		"node3": "localhost:50053",
	}

	// Step 2: Generate chunk locations for each chunk
	var chunkLocations []*pb.ChunkLocation

	for i := int64(0); i < req.Count; i++ {
		// Generate unique chunk_id
		chunkID := uuid.New().String()

		// Select top 3 nodes with most capacity (round-robin for now as capacity not implemented)
		var replicas []*pb.DataNodeEndpoint

		// Select 3 replicas based on available nodes
		replicaCount := int64(3)
		if req.Count < 3 {
			replicaCount = req.Count
		}

		for j := int64(0); j < replicaCount; j++ {
			nodeIndex := (i + j) % int64(len(availableNodes))
			nodeID := availableNodes[nodeIndex]

			replicas = append(replicas, &pb.DataNodeEndpoint{
				NodeId:  nodeID,
				Address: nodeAddresses[nodeID],
			})
		}

		// Create ChunkLocation
		chunkLocation := &pb.ChunkLocation{
			ChunkId:  chunkID,
			Index:    int32(i),
			Replicas: replicas,
		}

		chunkLocations = append(chunkLocations, chunkLocation)
	}

	// Step 3: Return AllocateChunksResponse (NO database writes yet!)
	return &pb.AllocateChunksResponse{
		FileId: req.FileId,
		Chunks: chunkLocations,
	}, nil
}
