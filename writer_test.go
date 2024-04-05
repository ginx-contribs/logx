package logx

import "testing"

func TestDefaultWriter(t *testing.T) {
	writer, err := NewWriter(nil)
	if err != nil {
		t.Fatal(err)
	}

	bytes := []byte("hello world")
	n, err := writer.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(bytes) {
		t.Fatalf("expected %d bytes, got %d", len(bytes), n)
	}
}

func TestFileWriter(t *testing.T) {
	writer, err := NewWriter(&WriterOptions{
		Filename: "./testdata/access.log",
		Combine:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Close()

	bytes := []byte("hello world")
	n, err := writer.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(bytes) {
		t.Fatalf("expected %d bytes, got %d", len(bytes), n)
	}
}
