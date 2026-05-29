package main

import (
	"reflect"
	"testing"
)

func TestWordSearch_Search(t *testing.T) {
	tests := []struct {
		name     string
		board    [][]byte
		word     string
		expected bool
		result   []WordSearchResult
	}{
		{
			name: "word exists",
			board: [][]byte{
				{'A', 'B', 'C', 'E'},
				{'S', 'F', 'C', 'S'},
				{'A', 'D', 'E', 'E'},
			},
			word:     "ABCCED",
			expected: true,
			result: []WordSearchResult{
				WordSearchResult{0, 0, "A"},
				WordSearchResult{0, 1, "B"},
				WordSearchResult{0, 2, "C"},
				WordSearchResult{1, 2, "C"},
				WordSearchResult{2, 2, "E"},
				WordSearchResult{2, 1, "D"},
			},
		},
		{
			name: "another word exists",
			board: [][]byte{
				{'A', 'B', 'C', 'E'},
				{'S', 'F', 'C', 'S'},
				{'A', 'D', 'E', 'E'},
			},
			word:     "SEE",
			expected: true,
			result: []WordSearchResult{
				WordSearchResult{1, 3, "S"},
				WordSearchResult{2, 3, "E"},
				WordSearchResult{2, 2, "E"},
			},
		},
		{
			name: "cannot reuse same cell",
			board: [][]byte{
				{'A', 'B', 'C', 'E'},
				{'S', 'F', 'C', 'S'},
				{'A', 'D', 'E', 'E'},
			},
			word:     "ABCB",
			expected: false,
			result:   []WordSearchResult{},
		},
		{
			name: "empty word",
			board: [][]byte{
				{'A', 'B', 'C', 'E'},
				{'S', 'F', 'C', 'S'},
				{'A', 'D', 'E', 'E'},
			},
			word:     "",
			expected: false,
			result:   []WordSearchResult{},
		},
		{
			name: "single cell match",
			board: [][]byte{
				{'A'},
			},
			word:     "A",
			expected: true,
			result: []WordSearchResult{
				WordSearchResult{0, 0, "A"},
			},
		},
		{
			name: "single cell unmatch",
			board: [][]byte{
				{'A'},
			},
			word:     "B",
			expected: false,
			result:   []WordSearchResult{},
		},
		{
			name: "single cell cannot be reused",
			board: [][]byte{
				{'A'},
			},
			word:     "AA",
			expected: false,
			result:   []WordSearchResult{},
		},
		{
			name: "diagonal is not allowed",
			board: [][]byte{
				{'A', 'X'},
				{'X', 'B'},
			},
			word:     "AB",
			expected: false,
			result:   []WordSearchResult{},
		},
		{
			name: "word longer than board cells",
			board: [][]byte{
				{'A', 'B'},
				{'C', 'D'},
			},
			word:     "ABCDE",
			expected: false,
			result:   []WordSearchResult{},
		},
		{
			name: "passed repeated letters",
			board: [][]byte{
				{'A', 'A'},
				{'A', 'A'},
			},
			word:     "AAA",
			expected: true,
			result: []WordSearchResult{
				WordSearchResult{0, 0, "A"},
				WordSearchResult{0, 1, "A"},
				WordSearchResult{1, 1, "A"},
			},
		},
		{
			name: "pass whole matrix repeated letters",
			board: [][]byte{
				{'A', 'A'},
				{'A', 'A'},
			},
			word:     "AAAA",
			expected: true,
			result: []WordSearchResult{
				WordSearchResult{0, 0, "A"},
				WordSearchResult{0, 1, "A"},
				WordSearchResult{1, 1, "A"},
				WordSearchResult{1, 0, "A"},
			},
		},
		{
			name: "failed repeated letters",
			board: [][]byte{
				{'A', 'A'},
				{'A', 'A'},
			},
			word:     "AAAAA",
			expected: false,
			result:   []WordSearchResult{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ws, err := NewWordSearch(tt.board)
				if err != nil {
					t.Errorf("NewWordSearch() error = %v", err)
				}
				exist, result := ws.Search(tt.word)

				if exist != tt.expected {
					t.Fatalf(
						"Search(%q) = %v, expected %v",
						tt.word,
						exist,
						tt.expected,
					)
				}
				if !reflect.DeepEqual(result, tt.result) {
					t.Fatalf(
						"Search(%q) = %v, expected %v",
						tt.word,
						result,
						tt.result,
					)
				}
			},
		)
	}
}

func TestNewWordSearch(t *testing.T) {
	tests := []struct {
		name        string
		matrix      [][]byte
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "nil matrix returns error",
			matrix:      nil,
			wantErr:     true,
			expectedErr: "rowCount and colCount must be greater than zero",
		},
		{
			name:        "empty matrix returns error",
			matrix:      [][]byte{},
			wantErr:     true,
			expectedErr: "rowCount and colCount must be greater than zero",
		},
		{
			name: "matrix with empty first row returns error",
			matrix: [][]byte{
				{},
			},
			wantErr:     true,
			expectedErr: "rowCount and colCount must be greater than zero",
		},
		{
			name: "valid single cell matrix",
			matrix: [][]byte{
				{'A'},
			},
			wantErr: false,
		},
		{
			name: "valid multiple rows and columns matrix",
			matrix: [][]byte{
				{'A', 'B', 'C', 'E'},
				{'S', 'F', 'C', 'S'},
				{'A', 'D', 'E', 'E'},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := NewWordSearch(tt.matrix)

				if tt.wantErr {
					if err == nil {
						t.Fatalf("expected error, got nil")
					}

					if err.Error() != tt.expectedErr {
						t.Fatalf(
							"expected error %q, got %q",
							tt.expectedErr,
							err.Error(),
						)
					}

					if got != nil {
						t.Fatalf("expected WordSearch to be nil, got %#v", got)
					}

					return
				}

				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if got == nil {
					t.Fatalf("expected WordSearch instance, got nil")
				}
			},
		)
	}
}
