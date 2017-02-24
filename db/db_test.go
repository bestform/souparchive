package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadArchive(t *testing.T) {
	a := NewArchive("fixtures/archive.json")
	a.Read()

	if len(a.Data.Items) != 2 {
		t.Fatalf("Expected 2 items in archive. Got %d", len(a.Data.Items))
	}

	if a.Data.Items[0].Guid != "1" {
		t.Fatalf("Expected first element to be '1', got '%s'", a.Data.Items[0].Guid)
	}
	if a.Data.Items[1].Guid != "2" {
		t.Fatalf("Expected second element to be '2', got '%s'", a.Data.Items[0].Guid)
	}
	if a.Data.Items[0].Timestamp != 100 {
		t.Fatalf("Expected first element to have timestamp '100', got '%s'", a.Data.Items[0].Timestamp)
	}
	if a.Data.Items[1].Timestamp != 200 {
		t.Fatalf("Expected second element to have timestamp '200', got '%s'", a.Data.Items[0].Timestamp)
	}
}

func TestAddToAndPersistArchive(t *testing.T) {
	tempdir := os.TempDir()

	a := NewArchive(filepath.Join(tempdir, "archive.json"))

	a.Add("foo", 100)
	a.Add("bar", 200)

	a.Persist()

	b := NewArchive(filepath.Join(tempdir, "archive.json"))
	b.Read()

	if len(b.Data.Items) != 2 {
		t.Fatalf("Expected 2 items in archive. Got %d", len(b.Data.Items))
	}
	if b.Data.Items[0].Guid != "foo" {
		t.Fatalf("Expected first element to be 'foo', got '%s'", b.Data.Items[0].Guid)
	}
	if b.Data.Items[1].Guid != "bar" {
		t.Fatalf("Expected second element to be 'bar', got '%s'", b.Data.Items[0].Guid)
	}
	if a.Data.Items[0].Timestamp != 100 {
		t.Fatalf("Expected first element to have timestamp '100', got '%s'", a.Data.Items[0].Timestamp)
	}
	if a.Data.Items[1].Timestamp != 200 {
		t.Fatalf("Expected second element to have timestamp '200', got '%s'", a.Data.Items[0].Timestamp)
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
