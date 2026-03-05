package storage

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// parseTime 尝试解析多种格式的时间字符串
func parseTime(s string) time.Time {
	if s == "" {
		return time.Now()
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05-07:00",
		"2006/01/02 15:04:05",
		"2006-01-02",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t
		}
	}
	// 如果都解析失败，返回当前时间
	return time.Now()
}

// migrateTable 迁移单个表的数据
func migrateTable(db *gorm.DB, tableName string, model interface{}, timeFields []string) error {
	// 检查表是否存在
	if !db.Migrator().HasTable(tableName) {
		return nil
	}

	// 检查字段类型
	needsMigration := false
	for _, field := range timeFields {
		var columnType string
		// SQLite 获取列类型
		err := db.Raw("SELECT type FROM pragma_table_info(?) WHERE name = ?", tableName, field).Scan(&columnType).Error
		if err != nil {
			return err
		}
		// 如果是 TEXT 类型，说明需要迁移 (GORM 默认 time.Time 是 datetime)
		if strings.EqualFold(columnType, "text") {
			needsMigration = true
			break
		}
	}

	if !needsMigration {
		return nil
	}

	fmt.Printf("Migrating table %s from TEXT time fields to DATETIME...\n", tableName)

	// 开始事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 重命名旧表
		backupTableName := tableName + "_backup_" + time.Now().Format("20060102150405")
		if err := tx.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tableName, backupTableName)).Error; err != nil {
			return err
		}

		// 2. 创建新表 (通过 AutoMigrate)
		if err := tx.AutoMigrate(model); err != nil {
			return err
		}

		// 3. 读取旧数据并插入新表
		// 使用 map[string]interface{} 读取，以获取原始字符串
		var oldRows []map[string]interface{}
		if err := tx.Table(backupTableName).Find(&oldRows).Error; err != nil {
			return err
		}

		for _, row := range oldRows {
			// 处理时间字段
			for _, field := range timeFields {
				if val, ok := row[field]; ok {
					if strVal, ok := val.(string); ok {
						row[field] = parseTime(strVal)
					} else if val == nil {
						row[field] = time.Now() // 或者保持 nil，如果字段允许 null
					}
				}
			}
			// 插入新表
			if err := tx.Table(tableName).Create(row).Error; err != nil {
				return err
			}
		}

		// 4. 删除旧表 (可选，为了安全可以先保留，或者确认无误后删除)
		// 这里选择保留备份表，或者可以记录日志提示用户删除
		// if err := tx.Exec(fmt.Sprintf("DROP TABLE %s", backupTableName)).Error; err != nil {
		// 	return err
		// }
		fmt.Printf("Migration completed for %s. Backup table: %s\n", tableName, backupTableName)
		return nil
	})
}

// CheckAndMigrate 检查并执行数据迁移
func CheckAndMigrate(db *gorm.DB) error {
	// 迁移 BossData
	if err := migrateTable(db, "boss_data", &BossData{}, []string{"created_at", "updated_at"}); err != nil {
		return fmt.Errorf("migrating boss_data failed: %w", err)
	}

	// 迁移 Blacklist
	if err := migrateTable(db, "blacklist", &Blacklist{}, []string{"created_at"}); err != nil {
		return fmt.Errorf("migrating blacklist failed: %w", err)
	}

	// 迁移 DeliveryRecord
	if err := migrateTable(db, "delivery_record", &DeliveryRecord{}, []string{"created_at", "delivered_at"}); err != nil {
		return fmt.Errorf("migrating delivery_record failed: %w", err)
	}

	return nil
}
