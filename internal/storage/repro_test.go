package storage

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// BossDataWithTime is the current struct
type BossDataWithTime struct {
	ID        int64     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (BossDataWithTime) TableName() string {
	return "boss_data"
}

func TestReproduction(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "repro.db")

	// 1. Initialize DB manually to simulate old state
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err)

	// Create table with TEXT columns (SQLite default for string)
	// We use a struct with string fields to create the table and insert bad data
	type BossDataString struct {
		ID        int64  `gorm:"primaryKey"`
		CreatedAt string `gorm:"column:created_at"`
		UpdatedAt string `gorm:"column:updated_at"`
	}

	err = db.Table("boss_data").AutoMigrate(&BossDataString{})
	require.NoError(t, err)

	// Insert data with format that might be problematic
	// GORM default format is roughly "2006-01-02 15:04:05.999999999-07:00"
	// Let's insert "2023/01/01 10:00:00" which is not standard GORM format
	badTimeStr := "2023/01/01 10:00:00"
	err = db.Table("boss_data").Create(&BossDataString{
		CreatedAt: badTimeStr,
		UpdatedAt: badTimeStr,
	}).Error
	require.NoError(t, err)

	// 2. Try to read using the new struct with time.Time
	var result BossDataWithTime
	err = db.First(&result).Error
	
	// This might fail or result might have zero time
	if err != nil {
		t.Logf("Read failed as expected: %v", err)
	} else {
		t.Logf("Read succeeded, CreatedAt: %v", result.CreatedAt)
		// Check if it was parsed correctly
		// If GORM/SQLite driver can't parse, it typically returns zero time
		if result.CreatedAt.IsZero() {
			t.Log("CreatedAt is zero, meaning parsing failed silently")
		} else {
			// Check if it matches expected time
			expected, _ := time.Parse("2006/01/02 15:04:05", badTimeStr)
			if !result.CreatedAt.Equal(expected) {
				t.Logf("CreatedAt mismatch. Got %v, want %v", result.CreatedAt, expected)
			} else {
				t.Log("Surprisingly, it worked!")
			}
		}
	}
}
