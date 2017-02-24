package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadArchive(t *testing.T) {
	a := NewArchive("fixtures/archive.json")
	a.Read()

	if len(a.Data.Guid) != 2 {
		t.Fatalf("Expected 2 items in archive. Got %d", len(a.Data.Guid))
	}

	if a.Data.Guid[0] != "1" {
		t.Fatalf("Expected first element to be '1', got '%s'", a.Data.Guid[0])
	}
	if a.Data.Guid[1] != "2" {
		t.Fatalf("Expected second element to be '2', got '%s'", a.Data.Guid[0])
	}
}

func TestAddToAndPersistArchive(t *testing.T) {
	tempdir := os.TempDir()

	a := NewArchive(filepath.Join(tempdir, "archive.json"))

	a.Add("foo")
	a.Add("bar")

	a.Persist()

	b := NewArchive(filepath.Join(tempdir, "archive.json"))
	b.Read()

	if len(b.Data.Guid) != 2 {
		t.Fatalf("Expected 2 items in archive. Got %d", len(b.Data.Guid))
	}
	if b.Data.Guid[0] != "foo" {
		t.Fatalf("Expected first element to be 'foo', got '%s'", b.Data.Guid[0])
	}
	if b.Data.Guid[1] != "bar" {
		t.Fatalf("Expected second element to be 'bar', got '%s'", b.Data.Guid[0])
	}
}

func TestArchiveContains(t *testing.T) {
	a := NewArchive("fixtures/archive.json")
	a.Read()

	if a.Contains("1") == false {
		t.Fatal("Expected archive to contain '1', but it didn't")
	}
	if a.Contains("2") == false {
		t.Fatal("Expected archive to contain '2', but it didn't")
	}
	if a.Contains("3") == true {
		t.Fatal("Expected archive to not contain '3', but it did")
	}

}
