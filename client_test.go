//go:build integration
// +build integration

package simpleturso

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestFullIntegration(t *testing.T) {
	godotenv.Load(".env")

	dbUrl := os.Getenv("TURSO_URL")
	dbKey := os.Getenv("TURSO_KEY")

	if dbUrl == "" || dbKey == "" {
		t.Fatal("TURSO_URL or TURSO_KEY environment variables are missing")
	}

	// 2. Test Initialization
	// This will internally run 'SELECT 1' to verify the connection
	t.Run("Initialize Package", func(t *testing.T) {
		conf := &Config{
			DbUrl: dbUrl,
			DbKey: dbKey,
		}
		err := Init(conf)
		if err != nil {
			t.Fatalf("Init failed: %v. Check if your URL starts with libsql:// or https://", err)
		}
	})

	// 3. Test Client Functionality (LogToTurso)
	t.Run("Log Entry", func(t *testing.T) {
		// Ensure logs table exists first
		setupQuery := `CREATE TABLE IF NOT EXISTS logs (
			timestamp TEXT, 
			application TEXT, 
			level TEXT, 
			message TEXT
		)`
		err := Execute(setupQuery, nil)
		if err != nil {
			t.Fatalf("Failed to ensure logs table: %v", err)
		}

		mockErr := fmt.Errorf("test error")

		err = LogToTurso("INFO", mockErr, "Integration test log", "simple_turso_sdk", "UTC")
		if err != nil {
			t.Errorf("LogToTurso failed: %v", err)
		}
	})

	// 4. Test Generic Select
	t.Run("Generic Select", func(t *testing.T) {
		type Result struct {
			Val float64
		}

		query := "SELECT 100 as val"
		res, err := Select(query, nil, func(row map[string]interface{}) (Result, error) {
			return Result{Val: row["val"].(float64)}, nil
		})

		if err != nil {
			t.Fatalf("Select failed: %v", err)
		}

		if len(res) == 0 || res[0].Val != 100 {
			t.Errorf("Expected 100, got %v", res)
		}
	})
}

func TestInitValidationErrors(t *testing.T) {
	// Test Invalid Key
	confBadKey := &Config{
		DbUrl: "https://valid.turso.io",
		DbKey: "not-a-jwt",
	}
	if err := Init(confBadKey); err == nil || err.Error() != "invalid db key" {
		t.Errorf("Expected 'invalid db key' error, got: %v", err)
	}

	// Test Invalid URL
	confBadUrl := &Config{
		DbUrl: "https://wrong-domain.com",
		DbKey: "hdr.pay.sig",
	}
	// Note: If Init was already called successfully in another test,
	// sync.Once might skip the config assignment, but the validation
	// check at the top of Init happens EVERY time.
	if err := Init(confBadUrl); err == nil || err.Error() != "invalid db url" {
		t.Errorf("Expected 'invalid db url' error, got: %v", err)
	}
}
