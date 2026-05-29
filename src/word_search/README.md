# Task 1 — Word Search

## 1. Overview

This repository contains the Go implementation for **Task 1: Word Search** from the programming assessment.

The task is to determine whether a given word exists in a two-dimensional character board. A word can be constructed by moving from one cell to another adjacent cell. In this problem, adjacency is limited to **horizontal** and **vertical** neighbours only. Diagonal movement is not allowed. The same board cell must not be reused more than once while constructing a single word.

The implementation is written in Go and focuses on correctness, readability, testability, and a small but useful novelty: in addition to returning whether the word exists, the solution also returns the actual path used to construct the word.

---

## 2. Problem Statement

Given a 2D board of letters and a target word, determine whether the word can be formed by walking through sequentially adjacent cells.

A valid move can go only in one of four directions:

- right
- left
- down
- up

A cell can be used at most once in the same word path.

### Example Board

```text
A B C E
S F C S
A D E E
```

### Expected Results

```text
ABCCED => true
SEE    => true
ABCB   => false
```

`ABCB` is false because the valid path would require reusing the same `B` cell, which is not allowed.

---

## 3. Implemented API

The solution exposes a small interface:

```go
type WordSearch interface {
    Search(string) (bool, []WordSearchResult)
}
```

The constructor is:

```go
func NewWordSearch(matrix [][]byte) (WordSearch, error)
```

The returned result type is:

```go
type WordSearchResult struct {
    RowIndex    int
    ColumnIndex int
    Letter      string
}
```

This means the search operation returns two values:

```go
exists, path := ws.Search("ABCCED")
```

- `exists` tells whether the word exists in the board.
- `path` contains the exact cells used to form the word when a valid path is found.

Example path for `ABCCED`:

```text
(0,0) A -> (0,1) B -> (0,2) C -> (1,2) C -> (2,2) E -> (2,1) D
```

---

## 4. Theoretical Solution

The core algorithm is **Depth-First Search with Backtracking**.

The search starts from every cell in the board. For each starting cell, the algorithm tries to match the first character of the word. If the character matches, it recursively explores the four valid neighbouring directions to match the next character.

A `visited` matrix is used to make sure the same cell is not reused within the current search path. If a recursive path fails, the algorithm backtracks by unmarking the current cell and removing it from the current result path.

The DFS stops successfully when the current word index reaches the length of the target word.

### DFS Decision Rules

At each recursive step, the algorithm rejects the current cell if:

1. the row index is outside the board;
2. the column index is outside the board;
3. the cell has already been visited in the current path;
4. the board character does not match the expected word character.

If none of these conditions apply, the cell is accepted as part of the current path and DFS continues to its neighbours.

---

## 5. Implementation Design

### 5.1 Constructor Validation

The constructor validates that the board has at least one row and at least one column:

```go
func NewWordSearch(matrix [][]byte) (WordSearch, error)
```

Invalid inputs return an error:

```text
rowCount and colCount must be greater than zero
```

Covered invalid cases include:

- `nil` matrix
- empty matrix
- matrix with an empty first row

### 5.2 Search Flow

The public `Search` method applies three early checks before running DFS:

1. Empty word is rejected.
2. A word longer than the total number of board cells is rejected.
3. A word requiring more occurrences of a letter than the board contains is rejected.

After these validations, the algorithm tries every board cell as a potential starting point.

### 5.3 Visited Matrix

For every starting cell, a fresh visited matrix is created:

```go
visited := makeNewVisitedMatrix(ws.rowCount, ws.colCount)
```

This keeps each DFS attempt isolated and avoids leaking visited state from one start position into another.

### 5.4 Path Tracking

The implementation stores the successful route in a map keyed by the word index:

```go
results[index] = WordSearchResult{
    RowIndex:    row,
    ColumnIndex: col,
    Letter:      string(word[index]),
}
```

When backtracking occurs, the corresponding path entry is removed:

```go
delete(*results, index)
```

This allows the final output to explain not only that a word exists, but also how it exists in the board.

---

## 6. Novelty Added

Although the classic solution to Word Search is DFS with backtracking, this implementation adds practical improvements around the core algorithm.

### 6.1 Path-Aware Search Result

Most simple Word Search implementations return only `true` or `false`. This implementation also returns the exact path of the discovered word.

This makes the result more explainable and easier to debug. For assessment purposes, it demonstrates that the algorithm is not treated as a black box; the solution can show the concrete sequence of board cells used to form the word.

Example:

```text
Word: ABCCED
Result: true
Path:
0: row=0, col=0, letter=A
1: row=0, col=1, letter=B
2: row=0, col=2, letter=C
3: row=1, col=2, letter=C
4: row=2, col=2, letter=E
5: row=2, col=1, letter=D
```

### 6.2 Frequency-Based Pre-Search Pruning

Before running DFS, the implementation counts the letters in the board and the letters required by the word.

If the board does not contain enough occurrences of any required letter, the method returns `false` immediately.

Example:

```text
Board:
A A
A A

Word: AAAAA
```

The board contains only four `A` cells, but the word requires five `A` characters. Therefore, the algorithm rejects the word before starting DFS.

This avoids unnecessary recursive search and improves performance for impossible inputs.

### 6.3 Explicit Test Coverage for Problem Constraints

The test suite does not only test the happy path. It also verifies the important constraints from the problem statement:

- a cell cannot be reused;
- diagonal movement is not allowed;
- empty words are rejected;
- single-cell boards work correctly;
- words longer than the board capacity are rejected;
- repeated-letter boards are handled correctly;
- invalid boards return constructor errors.

---

## 7. Complexity Analysis

Let:

```text
R = number of rows
C = number of columns
L = length of the target word
```

In the worst case, DFS may start from every cell. From each cell, it can explore up to four directions for each character of the word.

Therefore, the worst-case time complexity is:

```text
O(R × C × 4^L)
```

In practice, the effective branching factor is usually lower because:

- visited cells cannot be reused;
- character mismatches stop recursion early;
- frequency pruning rejects impossible words before DFS starts.

The space complexity is:

```text
O(R × C + L)
```

The `R × C` term comes from the visited matrix. The `L` term comes from the recursive call stack and the stored result path.

---

## 8. Test Cases

The project includes table-driven Go tests in `word_search_test.go`.

### 8.1 Happy Cases

| Test Case | Board Type | Word | Expected |
|---|---|---:|---:|
| Word exists | Standard example board | `ABCCED` | `true` |
| Another word exists | Standard example board | `SEE` | `true` |
| Single cell match | 1×1 board | `A` | `true` |
| Repeated letters | 2×2 board of `A`s | `AAA` | `true` |
| Whole matrix repeated letters | 2×2 board of `A`s | `AAAA` | `true` |

### 8.2 Corner Cases

| Test Case | Purpose | Expected |
|---|---|---:|
| Empty word | Reject invalid empty search input | `false` |
| Single cell mismatch | Verify one-cell failure | `false` |
| Single cell cannot be reused | Enforce no cell reuse | `false` |
| Diagonal is not allowed | Enforce horizontal/vertical adjacency only | `false` |
| Word longer than board cells | Reject impossible word length | `false` |
| Failed repeated letters | Reject word requiring more letters than available | `false` |
| Nil matrix | Constructor validation | error |
| Empty matrix | Constructor validation | error |
| Empty first row | Constructor validation | error |

---

## 9. How to Run

### 9.1 Prerequisites

Install Go:

```bash
go version
```

The implementation uses only the Go standard library.

### 9.2 Run Tests

From the project directory:

```bash
go test ./...
```

Expected result:

```text
ok
```

### 9.3 Run Tests Verbosely

```bash
go test -v ./...
```

---

## 10. Assumptions

This implementation follows the assessment problem statement and uses the following assumptions:

1. The board is represented as `[][]byte`.
2. The target word is represented as a Go `string`.
3. Movement is allowed only horizontally and vertically.
4. Diagonal movement is not valid.
5. A cell cannot be reused in the same search path.
6. Empty words are treated as invalid and return `false`.
7. The constructor rejects boards with zero rows or zero columns.
8. The implementation assumes a rectangular matrix shape based on the first row length.

---

## 11. Notes for Further Improvement

The current solution is intentionally focused and readable. Possible future improvements include:

1. Validate that all rows have the same column length.
2. Reverse the search word when the last character is rarer than the first character, reducing the number of starting positions.
3. Replace the visited matrix with in-place marking to reduce auxiliary memory, while keeping the current approach for clarity and safety.
4. Add benchmark tests for large boards and repeated-character worst-case scenarios.
5. Add a small CLI wrapper so users can provide a board and word from the command line.

---

## 12. Summary

This solution implements the Word Search task using DFS with backtracking, supported by input validation, early pruning, path tracking, and comprehensive table-driven tests.

The main strength of the implementation is that it is not limited to a boolean answer. It can also explain the successful path, making the result transparent, debuggable, and easier to evaluate.
