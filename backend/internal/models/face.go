package models

import (
"time"

"gorm.io/gorm"
)

// FaceDevice represents a face recognition device
type FaceDevice struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceNo     string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"device_no"`
	DeviceName   string         `gorm:"type:varchar(100);not null" json:"device_name"`
	DeviceType   int8           `gorm:"type:tinyint;default:1" json:"device_type"` // 1-门禁设备，2-签到设备，3-通用设备
	Location     string         `gorm:"type:varchar(255)" json:"location"`
	IPAddress    string         `gorm:"type:varchar(50)" json:"ip_address"`
	Port         int            `gorm:"type:int" json:"port"`
	Brand        string         `gorm:"type:varchar(50)" json:"brand"`        // 设备品牌
	Model        string         `gorm:"type:varchar(50)" json:"model"`        // 设备型号
	SerialNumber string         `gorm:"type:varchar(100)" json:"serial_number"` // 序列号
	Status       int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-在线，2-离线，3-故障，4-停用
	LastOnline   *time.Time     `json:"last_online"`
	APIKey       string         `gorm:"type:varchar(255)" json:"api_key"` // 设备API密钥
	APISecret    string         `gorm:"type:varchar(255)" json:"api_secret"` // 设备API密钥
	Config       string         `gorm:"type:text" json:"config"` // JSON格式的设备配置
	Remark       string         `gorm:"type:text" json:"remark"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FaceDevice) TableName() string {
	return "face_devices"
}

// UserFace represents user's face data
type UserFace struct {
ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
UserID       int64          `gorm:"index;not null" json:"user_id"`
FaceID       string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"face_id"` // 人脸库ID
FaceToken    string         `gorm:"type:varchar(255)" json:"face_token"` // 人脸特征值token
FaceImageURL string         `gorm:"type:varchar(500)" json:"face_image_url"` // 人脸照片URL
Quality      float64        `gorm:"type:decimal(5,2)" json:"quality"` // 人脸质量分数
IsMain       bool           `gorm:"default:false" json:"is_main"` // 是否为主照片
Status       int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-正常，2-已删除，3-待审核
RegisteredAt time.Time      `json:"registered_at"` // 录入时间
DeviceID     *int64         `gorm:"index" json:"device_id"` // 录入设备ID
OperatorID   *int64         `gorm:"index" json:"operator_id"` // 操作人ID
Remark       string         `gorm:"type:text" json:"remark"`
CreatedAt    time.Time      `json:"created_at"`
UpdatedAt    time.Time      `json:"updated_at"`
DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserFace) TableName() string {
return "user_faces"
}

// FaceRecognitionLog represents face recognition log
type FaceRecognitionLog struct {
ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
DeviceID       int64          `gorm:"index;not null" json:"device_id"`
UserID         *int64         `gorm:"index" json:"user_id"` // 识别到的用户ID，未识别为null
FaceID         string         `gorm:"type:varchar(100);index" json:"face_id"` // 识别到的人脸ID
Confidence     float64        `gorm:"type:decimal(5,2)" json:"confidence"` // 识别置信度
RecognizedAt   time.Time      `gorm:"index" json:"recognized_at"` // 识别时间
CaptureURL     string         `gorm:"type:varchar(500)" json:"capture_url"` // 抓拍照片URL
IsSuccess      bool           `gorm:"index" json:"is_success"` // 是否识别成功
FailReason     string         `gorm:"type:varchar(255)" json:"fail_reason"` // 失败原因
Temperature    *float64       `gorm:"type:decimal(4,1)" json:"temperature"` // 体温（如果设备支持）
MaskDetection  *bool          `json:"mask_detection"` // 是否佩戴口罩（如果设备支持）
Action         int8           `gorm:"type:tinyint" json:"action"` // 1-签到，2-门禁通过，3-仅识别
CheckInID      *int64         `gorm:"index" json:"check_in_id"` // 关联的签到记录ID
CreatedAt      time.Time      `json:"created_at"`
DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FaceRecognitionLog) TableName() string {
return "face_recognition_logs"
}

// Device status constants
const (
DeviceStatusOnline  = 1 // 在线
DeviceStatusOffline = 2 // 离线
DeviceStatusFault   = 3 // 故障
DeviceStatusDisabled = 4 // 停用
)

// Device type constants
const (
DeviceTypeAccess  = 1 // 门禁设备
DeviceTypeCheckIn = 2 // 签到设备
DeviceTypeGeneral = 3 // 通用设备
)

// Face status constants
const (
FaceStatusNormal  = 1 // 正常
FaceStatusDeleted = 2 // 已删除
FaceStatusPending = 3 // 待审核
)

// Recognition action constants
const (
ActionCheckIn = 1 // 签到
ActionAccess  = 2 // 门禁通过
ActionRecognize = 3 // 仅识别
)

// GetDeviceStatusText returns the text representation of device status
func GetDeviceStatusText(status int8) string {
switch status {
case DeviceStatusOnline:
return "在线"
case DeviceStatusOffline:
return "离线"
case DeviceStatusFault:
return "故障"
case DeviceStatusDisabled:
return "停用"
default:
return "未知"
}
}

// GetDeviceTypeText returns the text representation of device type
func GetDeviceTypeText(deviceType int8) string {
switch deviceType {
case DeviceTypeAccess:
return "门禁设备"
case DeviceTypeCheckIn:
return "签到设备"
case DeviceTypeGeneral:
return "通用设备"
default:
return "未知"
}
}

// GetFaceStatusText returns the text representation of face status
func GetFaceStatusText(status int8) string {
switch status {
case FaceStatusNormal:
return "正常"
case FaceStatusDeleted:
return "已删除"
case FaceStatusPending:
return "待审核"
default:
return "未知"
}
}
