package backend

import(
	
)

func UploadChunkHandler(c *fiber.Ctx) error {
    // 1️⃣ Read metadata from headers
    filePath := c.Get("X-File-Path")
    ownerID := c.Get("X-Owner-Id")
    offsetStr := c.Get("X-Offset")

    if filePath == "" || ownerID == "" || offsetStr == "" {
        return fiber.NewError(fiber.StatusBadRequest, "missing headers")
    }

    offset, err := strconv.ParseInt(offsetStr, 10, 64)
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "invalid offset")
    }

    // 2️⃣ Read raw chunk bytes
    chunkData := c.Body()
    if len(chunkData) == 0 {
        return fiber.NewError(fiber.StatusBadRequest, "empty chunk")
    }

    // 3️⃣ Call backend DFS logic (NO gRPC HERE)
    err = .WriteChunk(
        filePath,
        ownerID,
        offset,
        chunkData,
    )
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }

    return c.JSON(fiber.Map{
        "status": "chunk received",
        "offset": offset,
        "size":   len(chunkData),
    })
}