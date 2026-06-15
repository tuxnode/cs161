package client

import (
	userlib "github.com/cs161-staff/project2-userlib"
)

type Access struct {
	FileKey   []byte
	InodeUUID userlib.UUID
}

type Invitation struct {
	FileKey   []byte
	InodeUUID userlib.UUID
}
