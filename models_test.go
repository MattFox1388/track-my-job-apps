package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDateOnly_Scan(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expected  time.Time
		shouldErr bool
	}{
		{
			name:      "scan string date",
			input:     "2025-09-10",
			expected:  time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
			shouldErr: false,
		},
		{
			name:      "scan time.Time",
			input:     time.Date(2025, 9, 10, 15, 30, 45, 0, time.UTC),
			expected:  time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
			shouldErr: false,
		},
		{
			name:      "scan nil",
			input:     nil,
			expected:  time.Time{},
			shouldErr: false,
		},
		{
			name:      "scan invalid string",
			input:     "invalid-date",
			expected:  time.Time{},
			shouldErr: true,
		},
		{
			name:      "scan invalid type",
			input:     123,
			expected:  time.Time{},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d DateOnly
			err := d.Scan(tt.input)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !d.Time.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, d.Time)
			}
		})
	}
}

func TestDateOnly_Value(t *testing.T) {
	tests := []struct {
		name     string
		input    DateOnly
		expected interface{}
		shouldErr bool
	}{
		{
			name:     "valid date",
			input:    DateOnly{time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC)},
			expected: "2025-09-10",
			shouldErr: false,
		},
		{
			name:     "zero time",
			input:    DateOnly{time.Time{}},
			expected: nil,
			shouldErr: false,
		},
		{
			name:     "date with time components",
			input:    DateOnly{time.Date(2025, 12, 25, 15, 30, 45, 123456789, time.UTC)},
			expected: "2025-12-25",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.input.Value()

			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, value)
			}
		})
	}
}

func TestDateOnly_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    DateOnly
		expected string
		shouldErr bool
	}{
		{
			name:     "valid date",
			input:    DateOnly{time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC)},
			expected: `"2025-09-10"`,
			shouldErr: false,
		},
		{
			name:     "zero time",
			input:    DateOnly{time.Time{}},
			expected: "null",
			shouldErr: false,
		},
		{
			name:     "date with time components ignored",
			input:    DateOnly{time.Date(2025, 12, 25, 15, 30, 45, 123456789, time.UTC)},
			expected: `"2025-12-25"`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.MarshalJSON()

			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

func TestDateOnly_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected time.Time
		shouldErr bool
	}{
		{
			name:     "valid date string",
			input:    []byte(`"2025-09-10"`),
			expected: time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
			shouldErr: false,
		},
		{
			name:     "null value",
			input:    []byte("null"),
			expected: time.Time{},
			shouldErr: false,
		},
		{
			name:     "invalid date format",
			input:    []byte(`"invalid-date"`),
			expected: time.Time{},
			shouldErr: true,
		},
		{
			name:     "valid date different year",
			input:    []byte(`"2024-12-31"`),
			expected: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d DateOnly
			err := d.UnmarshalJSON(tt.input)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !d.Time.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, d.Time)
			}
		})
	}
}

func TestDateOnly_JSONRoundTrip(t *testing.T) {
	type testStruct struct {
		Date DateOnly `json:"date"`
		Name string   `json:"name"`
	}

	original := testStruct{
		Date: DateOnly{time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC)},
		Name: "test",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expectedJSON := `{"date":"2025-09-10","name":"test"}`
	if string(jsonData) != expectedJSON {
		t.Errorf("expected JSON %s, got %s", expectedJSON, string(jsonData))
	}

	// Unmarshal back
	var unmarshaled testStruct
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !unmarshaled.Date.Time.Equal(original.Date.Time) {
		t.Errorf("expected %v, got %v", original.Date.Time, unmarshaled.Date.Time)
	}

	if unmarshaled.Name != original.Name {
		t.Errorf("expected name %s, got %s", original.Name, unmarshaled.Name)
	}
}

func TestDateOnly_NullJSONRoundTrip(t *testing.T) {
	type testStruct struct {
		Date DateOnly `json:"date"`
		Name string   `json:"name"`
	}

	original := testStruct{
		Date: DateOnly{time.Time{}}, // Zero time
		Name: "test",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expectedJSON := `{"date":null,"name":"test"}`
	if string(jsonData) != expectedJSON {
		t.Errorf("expected JSON %s, got %s", expectedJSON, string(jsonData))
	}

	// Unmarshal back
	var unmarshaled testStruct
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !unmarshaled.Date.Time.IsZero() {
		t.Errorf("expected zero time, got %v", unmarshaled.Date.Time)
	}

	if unmarshaled.Name != original.Name {
		t.Errorf("expected name %s, got %s", original.Name, unmarshaled.Name)
	}
}

func TestDateOnly_DatabaseIntegration(t *testing.T) {
	// Test the database driver interfaces work together
	original := DateOnly{time.Date(2025, 9, 10, 15, 30, 45, 0, time.UTC)}
	
	// Test Value() -> Scan() round trip
	value, err := original.Value()
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	var scanned DateOnly
	err = scanned.Scan(value)
	if err != nil {
		t.Fatalf("failed to scan value: %v", err)
	}

	expectedDate := time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC)
	if !scanned.Time.Equal(expectedDate) {
		t.Errorf("expected %v, got %v", expectedDate, scanned.Time)
	}
}
