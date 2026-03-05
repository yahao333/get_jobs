package storage

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yahao333/get_jobs/internal/config"
)

// DB 全局数据库实例
var DB *gorm.DB

// BossData Boss直聘岗位数据
type BossData struct {
	ID             int64  `gorm:"primaryKey" json:"id"`
	EncryptID      string `gorm:"column:encrypt_id" json:"encrypt_id"`
	EncryptUserID  string `gorm:"column:encrypt_user_id" json:"encrypt_user_id"`
	CompanyName    string `gorm:"column:company_name" json:"company_name"`
	JobName        string `gorm:"column:job_name" json:"job_name"`
	Salary         string `gorm:"column:salary" json:"salary"`
	Location       string `gorm:"column:location" json:"location"`
	Experience     string `gorm:"column:experience" json:"experience"`
	Degree         string `gorm:"column:degree" json:"degree"`
	HRName         string `gorm:"column:hr_name" json:"hr_name"`
	HRPosition     string `gorm:"column:hr_position" json:"hr_position"`
	HRActiveStatus string `gorm:"column:hr_active_status" json:"hr_active_status"`
	DeliveryStatus string `gorm:"column:delivery_status" json:"delivery_status"`
	JobDescription string `gorm:"column:job_description" json:"job_description"`
	JobURL         string `gorm:"column:job_url" json:"job_url"`
	CreatedAt      string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      string `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (BossData) TableName() string {
	return "boss_data"
}

// Blacklist 黑名单
type Blacklist struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	Keyword   string `gorm:"column:keyword" json:"keyword"`
	Type      string `gorm:"column:type" json:"type"`     // company, hr, job
	Source    string `gorm:"column:source" json:"source"` // manual, auto
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}

// TableName 指定表名
func (Blacklist) TableName() string {
	return "blacklist"
}

// DeliveryRecord 投递记录
type DeliveryRecord struct {
	ID          int64  `gorm:"primaryKey" json:"id"`
	JobID       int64  `gorm:"column:job_id" json:"job_id"`
	Platform    string `gorm:"column:platform" json:"platform"`
	Status      string `gorm:"column:status" json:"status"` // success, failed, pending
	Message     string `gorm:"column:message" json:"message"`
	DeliveredAt string `gorm:"column:delivered_at" json:"delivered_at"`
	CreatedAt   string `gorm:"column:created_at" json:"created_at"`
}

// TableName 指定表名
func (DeliveryRecord) TableName() string {
	return "delivery_record"
}

// InitDB 初始化数据库
func InitDB() error {
	dbPath := config.GetString("database.path")

	// 确保数据库目录存在
	dir := dbPath
	for i := len(dir) - 1; i >= 0; i-- {
		if dir[i] == '/' || dir[i] == '\\' {
			dir = dir[:i]
			break
		}
	}
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建数据库目录失败: %w", err)
		}
	}

	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("打开数据库失败: %w", err)
	}

	// 自动迁移表结构
	if err := autoMigrate(db); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	DB = db
	config.Info("数据库初始化成功: ", dbPath)
	return nil
}

// autoMigrate 自动迁移表结构
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&BossData{},
		&Blacklist{},
		&DeliveryRecord{},
	)
}

// Create 创建记录
func Create(data interface{}) error {
	return DB.Create(data).Error
}

// First 根据条件查询第一条记录
func First(result interface{}, conditions ...interface{}) error {
	return DB.First(result, conditions...).Error
}

// Where 根据条件查询
func Where(result interface{}, query interface{}, args ...interface{}) error {
	return DB.Where(query, args...).Find(result).Error
}

// Updates 更新记录
func Updates(data interface{}, conditions ...interface{}) error {
	return DB.Model(data).Updates(conditions).Error
}

// Delete 删除记录
func Delete(data interface{}, conditions ...interface{}) error {
	return DB.Delete(data, conditions...).Error
}

// Count 统计数量
func Count(model interface{}, query interface{}, args ...interface{}) (int64, error) {
	var count int64
	err := DB.Model(model).Where(query, args...).Count(&count).Error
	return count, err
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
