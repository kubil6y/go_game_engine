package type_registry

import (
	"testing"
)

func TestNewTypeRegistry(t *testing.T) {
	reg := NewTypeRegistry(5)
	if reg.maxItems != 5 {
		t.Errorf("Expected maxItems to be 5, got %d", reg.maxItems)
	}
	if reg.nextID != 0 {
		t.Errorf("Expected nextID to be 0, got %d", reg.nextID)
	}
	if len(reg.typeIDs) != 0 {
		t.Errorf("Expected typeIDs to be empty, got %d items", len(reg.typeIDs))
	}
}

func TestSizeTypeRegistry(t *testing.T) {
	type foo struct{}
	type bar struct{}
	reg := NewTypeRegistry(5)
	reg.Register(foo{})
	reg.Register(bar{})
	reg.Register(foo{})
	reg.Register(foo{})
	if reg.Size() != 2 {
		t.Errorf("Expected typeIDs to be empty, got %d items", reg.Size())
	}
}

func TestRegisterNilItem(t *testing.T) {
	reg := NewTypeRegistry(5)
	_, err := reg.Register(nil)
	if err != ErrNilItem {
		t.Errorf("Expected ErrNilItem, got %v", err)
	}
}

func TestRegisterExceedingMaxItems(t *testing.T) {
	reg := NewTypeRegistry(1) // Only allow one item
	_, err := reg.Register("first")
	if err != nil {
		t.Fatalf("Unexpected error during first registration: %v", err)
	}
	_, err = reg.Register("second")
	if err != ErrMaxItemsExceeded {
		t.Errorf("Expected ErrMaxItemsExceeded, got %v", err)
	}
}

func TestRegisterSameItem(t *testing.T) {
	reg := NewTypeRegistry(5)
	id1, err := reg.Register("item")
	if err != nil {
		t.Fatalf("Unexpected error during registration: %v", err)
	}
	id2, err := reg.Register("item")
	if err != nil {
		t.Fatalf("Unexpected error during second registration: %v", err)
	}
	if id1 != id2 {
		t.Errorf("Expected the same ID for the same item, got id1: %d, id2: %d", id1, id2)
	}
}

func TestGetRegisteredItem(t *testing.T) {
	reg := NewTypeRegistry(5)
	_, err := reg.Register("item")
	if err != nil {
		t.Fatalf("Unexpected error during registration: %v", err)
	}

	id, err := reg.Get("item")
	if err != nil {
		t.Fatalf("Unexpected error during get: %v", err)
	}
	if id != 0 {
		t.Errorf("Expected ID to be 0 for first registered item, got %d", id)
	}
}

func TestGetNilItem(t *testing.T) {
	reg := NewTypeRegistry(5)
	_, err := reg.Get(nil)
	if err != ErrNilItem {
		t.Errorf("Expected ErrNilItem, got %v", err)
	}
}

func TestGetUnregisteredItem(t *testing.T) {
	reg := NewTypeRegistry(5)
	_, err := reg.Get("item")
	if err != ErrTypeNotFound {
		t.Errorf("Expected ErrTypeNotFound, got %v", err)
	}
}

func TestConcurrentRegister(t *testing.T) {
	reg := NewTypeRegistry(10)
	done := make(chan struct{})

	for i := 0; i < 10; i++ {
		go func(i int) {
			defer func() { done <- struct{}{} }()
			_, err := reg.Register(i) // Registering an int
			if err != nil && err != ErrMaxItemsExceeded {
				t.Errorf("Unexpected error during concurrent registration: %v", err)
			}
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestConcurrentGet(t *testing.T) {
	reg := NewTypeRegistry(10)
	_, _ = reg.Register("item") // Registering an item

	done := make(chan struct{})

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			id, err := reg.Get("item")
			if err != nil {
				t.Errorf("Unexpected error during concurrent get: %v", err)
			}
			if id != 0 {
				t.Errorf("Expected ID to be 0 for registered item, got %d", id)
			}
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}
}
