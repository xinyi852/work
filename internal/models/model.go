package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"
)

const (
	// 开启状态
	STATUS_ON = 1
	// 关闭状态
	STATUS_OFF = 0
)

// BaseModel 模型基类
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
	CommonTimestampsField
}

type CommonTimestampsField struct {
	CreatedAt *NullTime `gorm:"column:created_at;" json:"created_at,omitempty"`
	UpdatedAt *NullTime `gorm:"column:updated_at;" json:"updated_at,omitempty"`
}

type NullTime struct {
	sql.NullTime
}

// 通过接口实现序列化
func (s NullTime) MarshalJSON() ([]byte, error) {
	if s.Valid {
		//res, err := s.Time.MarshalJSON() // 如果想修改时间格式，在这里修改
		return []byte(fmt.Sprintf("\"%v\"", s.Time.Format("2006-01-02 15:04:05"))), nil
		//return res, err
	} else {
		return json.Marshal("")
	}
}

// 反序列化
func (s *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "\"\"" {
		return nil
	}
	err := s.Time.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	if !s.Time.IsZero() {
		s.Valid = true
	}
	return nil
}

// GetStringID 获取 ID 的字符串格式
func (m BaseModel) GetStringID() string {
	return cast.ToString(m.ID)
}
