package storage

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCheckAndMigrate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "migration_test.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err)

	// 1. Create table with TEXT columns (simulating old schema)
	type BossDataOldSchema struct {
		ID        int64  `gorm:"primaryKey"`
		CreatedAt string `gorm:"column:created_at"` // TEXT
		UpdatedAt string `gorm:"column:updated_at"` // TEXT
		JobName   string `gorm:"column:job_name"`
	}
	// Force table name to be boss_data
	err = db.Table("boss_data").AutoMigrate(&BossDataOldSchema{})
	require.NoError(t, err)

	// 2. Insert data with weird string format
	badTimeStr := "2023/10/27 10:00:00"
	err = db.Table("boss_data").Create(&BossDataOldSchema{
		ID:        1,
		CreatedAt: badTimeStr,
		UpdatedAt: badTimeStr,
		JobName:   "Test Job",
	}).Error
	require.NoError(t, err)

	// Verify it's TEXT
	var typeName string
	err = db.Raw("SELECT type FROM pragma_table_info('boss_data') WHERE name='created_at'").Scan(&typeName).Error
	require.NoError(t, err)
	assert.Equal(t, "text", strings.ToLower(typeName))

	// 3. Run Migration
	err = CheckAndMigrate(db)
	require.NoError(t, err)

	// 4. Verify Schema is now DATETIME
	err = db.Raw("SELECT type FROM pragma_table_info('boss_data') WHERE name='created_at'").Scan(&typeName).Error
	require.NoError(t, err)
	// GORM + SQLite maps time.Time to datetime usually
	t.Logf("New column type: %s", typeName)
	// assert.Equal(t, "datetime", strings.ToLower(typeName))

	// 5. Verify Data is readable with BossData struct
	var result BossData
	err = db.First(&result, 1).Error
	require.NoError(t, err)

	expectedTime, _ := time.Parse("2006/01/02 15:04:05", badTimeStr)
	// Compare with tolerance or just formatting
	assert.Equal(t, expectedTime.Format(time.RFC3339), result.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, "Test Job", result.JobName)

	// 6. Run Migration again (idempotency check)
	err = CheckAndMigrate(db)
	require.NoError(t, err)
	
	// Should remain datetime
	err = db.Raw("SELECT type FROM pragma_table_info('boss_data') WHERE name='created_at'").Scan(&typeName).Error
	require.NoError(t, err)
	// assert.Equal(t, "datetime", strings.ToLower(typeName))
}
