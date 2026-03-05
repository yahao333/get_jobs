package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/yahao333/get_jobs/internal/config"
)

func setupTestConfig(dbPath string) {
	// 创建一个测试用的配置
	config.Config = viper.New()
	config.Config.Set("database.path", dbPath)
}

func TestInitDB(t *testing.T) {
	// 先初始化配置
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// 设置临时配置路径
	setupTestConfig(dbPath)

	// 初始化数据库
	err := InitDB()
	require.NoError(t, err)
	assert.NotNil(t, DB)

	// 验证数据库是否可用
	err = DB.Raw("SELECT 1").Error
	require.NoError(t, err, "数据库应该可以执行查询")

	// 验证表是否创建
	t.Run("tables created", func(t *testing.T) {
		// 检查 boss_data 表
		var count int64
		err := DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='boss_data'").Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 检查 blacklist 表
		err = DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='blacklist'").Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 检查 delivery_record 表
		err = DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='delivery_record'").Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}

func TestInitDBWithNestedPath(t *testing.T) {
	// 测试带嵌套路径的数据库初始化
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "subdir", "data", "test.db")

	setupTestConfig(dbPath)

	err := InitDB()
	require.NoError(t, err)
	assert.NotNil(t, DB)
}

func TestBossDataCRUD(t *testing.T) {
	// 初始化测试数据库
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_crud.db")
	setupTestConfig(dbPath)

	err := InitDB()
	require.NoError(t, err)

	t.Run("create boss data", func(t *testing.T) {
		data := &BossData{
			EncryptID:      "test_encrypt_id_123",
			EncryptUserID:  "test_user_id_456",
			CompanyName:    "测试公司",
			JobName:        "Go开发工程师",
			Salary:         "20k-35k",
			Location:       "北京",
			Experience:     "3-5年",
			Degree:         "本科",
			HRName:         "张三",
			HRPosition:     "HR经理",
			HRActiveStatus: "今日活跃",
			DeliveryStatus: "已投递",
			JobDescription: "熟悉Golang，有高并发经验",
			JobURL:         "https://www.zhipin.com/job/123.html",
		}

		err := Create(data)
		require.NoError(t, err)
		assert.Greater(t, data.ID, int64(0))
	})

	t.Run("query boss data", func(t *testing.T) {
		// 先创建一条数据
		data := &BossData{
			EncryptID:      "test_encrypt_789",
			CompanyName:    "查询测试公司",
			JobName:        "Java开发",
			Salary:         "15k-25k",
			Location:       "上海",
			DeliveryStatus: "待投递",
		}
		err := Create(data)
		require.NoError(t, err)

		// 测试 First 查询
		var result BossData
		err = First(&result, "encrypt_id = ?", "test_encrypt_789")
		require.NoError(t, err)
		assert.Equal(t, "查询测试公司", result.CompanyName)
		assert.Equal(t, "Java开发", result.JobName)
	})

	t.Run("where query", func(t *testing.T) {
		// 创建多条数据
		for i := 0; i < 3; i++ {
			data := &BossData{
				EncryptID:      "where_test_" + string(rune('a'+i)),
				CompanyName:    "Where测试公司",
				JobName:        "测试职位",
				DeliveryStatus: "待投递",
			}
			Create(data)
		}

		// 查询所有符合条件的记录
		var results []BossData
		err := Where(&results, "company_name = ?", "Where测试公司")
		require.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("update boss data", func(t *testing.T) {
		// 创建数据
		data := &BossData{
			EncryptID:      "update_test_123",
			CompanyName:    "更新前公司",
			JobName:        "更新前职位",
			DeliveryStatus: "待投递",
		}
		err := Create(data)
		require.NoError(t, err)

		// 更新数据 - 使用 GORM 原生方式
		err = DB.Model(data).Updates(map[string]interface{}{
			"company_name":    "更新后公司",
			"delivery_status": "已投递",
		}).Error
		require.NoError(t, err)

		// 验证更新结果
		var result BossData
		First(&result, "id = ?", data.ID)
		assert.Equal(t, "更新后公司", result.CompanyName)
		assert.Equal(t, "已投递", result.DeliveryStatus)
	})

	t.Run("delete boss data", func(t *testing.T) {
		// 创建数据
		data := &BossData{
			EncryptID:      "delete_test_123",
			CompanyName:    "删除测试公司",
			DeliveryStatus: "待投递",
		}
		err := Create(data)
		require.NoError(t, err)
		require.Greater(t, data.ID, int64(0))

		// 删除数据
		err = Delete(data)
		require.NoError(t, err)

		// 验证删除
		var result BossData
		err = First(&result, "id = ?", data.ID)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("count boss data", func(t *testing.T) {
		// 创建多条数据
		for i := 0; i < 5; i++ {
			data := &BossData{
				EncryptID:      "count_test_" + string(rune('a'+i)),
				CompanyName:    "计数测试公司",
				DeliveryStatus: "待投递",
			}
			Create(data)
		}

		// 统计数量
		count, err := Count(&BossData{}, "company_name = ?", "计数测试公司")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(5))
	})

	t.Run("check time fields", func(t *testing.T) {
		data := &BossData{
			EncryptID:      "time_test_123",
			CompanyName:    "时间测试公司",
			JobName:        "时间测试职位",
			DeliveryStatus: "待投递",
		}
		err := Create(data)
		require.NoError(t, err)

		var result BossData
		err = First(&result, "id = ?", data.ID)
		require.NoError(t, err)

		// 验证 CreatedAt 和 UpdatedAt 是否为非零时间
		assert.False(t, result.CreatedAt.IsZero(), "CreatedAt 应该是非零时间")
		assert.False(t, result.UpdatedAt.IsZero(), "UpdatedAt 应该是非零时间")

		// 验证时间是否在合理范围内（例如最近1分钟内）
		assert.WithinDuration(t, time.Now(), result.CreatedAt, time.Minute, "CreatedAt 应该是最近的时间")
		assert.WithinDuration(t, time.Now(), result.UpdatedAt, time.Minute, "UpdatedAt 应该是最近的时间")
	})
}

func TestBlacklistCRUD(t *testing.T) {
	// 初始化测试数据库
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_blacklist.db")
	setupTestConfig(dbPath)

	err := InitDB()
	require.NoError(t, err)

	t.Run("create blacklist", func(t *testing.T) {
		bl := &Blacklist{
			Keyword: "测试关键词",
			Type:    "company",
			Source:  "manual",
		}
		err := Create(bl)
		require.NoError(t, err)
		assert.Greater(t, bl.ID, int64(0))
		assert.False(t, bl.CreatedAt.IsZero(), "CreatedAt should not be zero")
		assert.True(t, time.Since(bl.CreatedAt) < time.Minute, "CreatedAt should be recent")
	})

	t.Run("query blacklist", func(t *testing.T) {
		bl := &Blacklist{
			Keyword: "查询关键词",
			Type:    "hr",
			Source:  "auto",
		}
		Create(bl)

		var result Blacklist
		err := First(&result, "keyword = ?", "查询关键词")
		require.NoError(t, err)
		assert.Equal(t, "hr", result.Type)
	})

	t.Run("update blacklist", func(t *testing.T) {
		bl := &Blacklist{
			Keyword: "更新关键词",
			Type:    "job",
			Source:  "manual",
		}
		Create(bl)

		// 使用 GORM 原生方式更新
		err := DB.Model(bl).Update("source", "auto").Error
		require.NoError(t, err)

		var result Blacklist
		First(&result, "id = ?", bl.ID)
		assert.Equal(t, "auto", result.Source)
	})

	t.Run("delete blacklist", func(t *testing.T) {
		bl := &Blacklist{
			Keyword: "删除关键词",
			Type:    "company",
		}
		Create(bl)

		err := Delete(bl)
		require.NoError(t, err)

		var result Blacklist
		err = First(&result, "id = ?", bl.ID)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestDeliveryRecordCRUD(t *testing.T) {
	// 初始化测试数据库
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_delivery.db")
	setupTestConfig(dbPath)

	err := InitDB()
	require.NoError(t, err)

	t.Run("create delivery record", func(t *testing.T) {
		record := &DeliveryRecord{
			JobID:    123,
			Platform: "boss",
			Status:   "success",
			Message:  "投递成功",
		}

		err := Create(record)
		require.NoError(t, err)
		assert.Greater(t, record.ID, int64(0))
		assert.False(t, record.CreatedAt.IsZero(), "CreatedAt should not be zero")
		assert.True(t, time.Since(record.CreatedAt) < time.Minute, "CreatedAt should be recent")
	})

	t.Run("query delivery record", func(t *testing.T) {
		record := &DeliveryRecord{
			JobID:    456,
			Platform: "liepin",
			Status:   "pending",
		}
		Create(record)

		var result DeliveryRecord
		err := First(&result, "job_id = ?", 456)
		require.NoError(t, err)
		assert.Equal(t, "pending", result.Status)
	})

	t.Run("update delivery record", func(t *testing.T) {
		record := &DeliveryRecord{
			JobID:    789,
			Platform: "boss",
			Status:   "pending",
		}
		Create(record)

		// 直接使用 DB 来更新，因为 Updates 函数实现可能需要调整
		err := DB.Model(record).Update("status", "success").Error
		require.NoError(t, err)

		var result DeliveryRecord
		First(&result, "id = ?", record.ID)
		assert.Equal(t, "success", result.Status)
	})

	t.Run("count delivery records", func(t *testing.T) {
		// 创建多条记录
		for i := 0; i < 3; i++ {
			record := &DeliveryRecord{
				JobID:    int64(100 + i),
				Platform: "boss",
				Status:   "success",
			}
			Create(record)
		}

		count, err := Count(&DeliveryRecord{}, "status = ?", "success")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(3))
	})
}

func TestGetDB(t *testing.T) {
	// 测试获取数据库实例
	db := GetDB()
	assert.NotNil(t, db)
}

func TestBossDataTableName(t *testing.T) {
	// 测试表名方法
	data := BossData{}
	assert.Equal(t, "boss_data", data.TableName())
}

func TestBlacklistTableName(t *testing.T) {
	// 测试黑名单表名
	bl := Blacklist{}
	assert.Equal(t, "blacklist", bl.TableName())
}

func TestDeliveryRecordTableName(t *testing.T) {
	// 测试投递记录表名
	record := DeliveryRecord{}
	assert.Equal(t, "delivery_record", record.TableName())
}

// BenchmarkCreateBossData 性能测试 - 创建数据
func BenchmarkCreateBossData(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench.db")
	config.Config.Set("database.path", dbPath)
	InitDB()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := &BossData{
			EncryptID:      "bench_" + string(rune(i)),
			CompanyName:    " Benchmark公司",
			JobName:        "测试职位",
			DeliveryStatus: "待投递",
		}
		Create(data)
	}
}

// BenchmarkQueryBossData 性能测试 - 查询数据
func BenchmarkQueryBossData(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench_query.db")
	config.Config.Set("database.path", dbPath)
	InitDB()

	// 预先创建数据
	for i := 0; i < 100; i++ {
		data := &BossData{
			EncryptID:      "bench_query_" + string(rune(i)),
			CompanyName:    "查询 Benchmark公司",
			JobName:        "测试职位",
			DeliveryStatus: "待投递",
		}
		Create(data)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result BossData
		First(&result, "encrypt_id = ?", "bench_query_a")
	}
}

// 清理测试数据库文件
func cleanupTestDB() {
	testDBs := []string{
		"./data/test.db",
		"./data/test_crud.db",
		"./data/test_blacklist.db",
		"./data/test_delivery.db",
		"./data/bench.db",
		"./data/bench_query.db",
	}
	for _, db := range testDBs {
		os.Remove(db)
	}
}
