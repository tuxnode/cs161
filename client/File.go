package client

import (
	"container/list"
	"math"

	userlib "github.com/cs161-staff/project2-userlib"
)

const BlockSize = 512

/* Divide A File Into FileBlocks */
type FileBlock struct {
	BlockUUID userlib.UUID
	block     [BlockSize]byte
}

/*
	 ByteToBlock: convert a File into Blocks
		data: byte stream
*/
func ByteToBlock(data []byte) (block *list.List) {
	size := len(data)
	blockNum := (size + BlockSize - 1) / BlockSize
	list := list.New()

	for i := 0; i < blockNum; i++ {
		start := i * BlockSize
		end := start + BlockSize

		if end > size {
			end = size
		}

		blockUUID := userlib.UUID(userlib.RandomBytes(userlib.UUIDSizeBytes))
		fileBlock := &FileBlock{
			BlockUUID: blockUUID,
		}
		copy(fileBlock.block[:], data[start:end])
		list.PushBack(fileBlock)
	}
	return list
}
