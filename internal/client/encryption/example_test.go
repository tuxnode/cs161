package encryption

import (
	"testing"
)

func TestServiceAPI(t *testing.T) {
	storage := NewMemoryStorage()
	keyStore := NewMemoryKeyStore()

	userService := NewUserService(storage, keyStore)
	fileService := NewFileService(storage, keyStore)
	invitationService := NewInvitationService(storage, keyStore)

	// Create users
	alice, err := userService.InitUser("alice", "password123")
	if err != nil {
		t.Fatalf("failed to create alice: %v", err)
	}

	bob, err := userService.InitUser("bob", "password456")
	if err != nil {
		t.Fatalf("failed to create bob: %v", err)
	}

	// Store file
	err = fileService.StoreFile(alice, "test.txt", []byte("Hello, World!"))
	if err != nil {
		t.Fatalf("failed to store file: %v", err)
	}

	// Load file
	content, err := fileService.LoadFile(alice, "test.txt")
	if err != nil {
		t.Fatalf("failed to load file: %v", err)
	}
	if string(content) != "Hello, World!" {
		t.Fatalf("file content mismatch: got %s", string(content))
	}

	// Share file
	invPtr, err := invitationService.CreateInvitation(alice, "test.txt", "bob")
	if err != nil {
		t.Fatalf("failed to create invitation: %v", err)
	}

	// Accept invitation
	err = invitationService.AcceptInvitation(bob, "alice", invPtr, "test.txt")
	if err != nil {
		t.Fatalf("failed to accept invitation: %v", err)
	}

	// Bob loads file
	content, err = fileService.LoadFile(bob, "test.txt")
	if err != nil {
		t.Fatalf("failed to load file as bob: %v", err)
	}
	if string(content) != "Hello, World!" {
		t.Fatalf("file content mismatch for bob: got %s", string(content))
	}

	// Revoke access
	err = invitationService.RevokeAccess(alice, "test.txt", "bob")
	if err != nil {
		t.Fatalf("failed to revoke access: %v", err)
	}

	// Bob should not be able to load file anymore
	_, err = fileService.LoadFile(bob, "test.txt")
	if err == nil {
		t.Fatal("bob should not be able to load file after revocation")
	}
}
