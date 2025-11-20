package service

import (
"errors"
"fmt"
"gym-admin/internal/models"
"gym-admin/internal/repository"
"time"
)

type FaceDeviceService struct {
	repo *repository.FaceDeviceRepository
}

func NewFaceDeviceService() *FaceDeviceService {
	return &FaceDeviceService{
		repo: repository.NewFaceDeviceRepository(),
	}
}

func (s *FaceDeviceService) CreateDevice(device *models.FaceDevice) error {
	// Generate device number
	device.DeviceNo = s.generateDeviceNo()
	device.Status = models.DeviceStatusOffline // Default status: offline
	
	return s.repo.Create(device)
}

func (s *FaceDeviceService) GetDevice(id int64) (*models.FaceDevice, error) {
	return s.repo.GetByID(id)
}

func (s *FaceDeviceService) GetDeviceByNo(deviceNo string) (*models.FaceDevice, error) {
	return s.repo.GetByDeviceNo(deviceNo)
}

func (s *FaceDeviceService) ListDevices(page, pageSize int, status *int8, deviceType *int8) ([]models.FaceDevice, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize, status, deviceType)
}

func (s *FaceDeviceService) UpdateDevice(id int64, updates map[string]interface{}) error {
	device, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("device not found")
	}

	// Update fields
	if deviceName, ok := updates["device_name"].(string); ok {
		device.DeviceName = deviceName
	}
	if deviceType, ok := updates["device_type"].(float64); ok {
		device.DeviceType = int8(deviceType)
	}
	if location, ok := updates["location"].(string); ok {
		device.Location = location
	}
	if ipAddress, ok := updates["ip_address"].(string); ok {
		device.IPAddress = ipAddress
	}
	if port, ok := updates["port"].(float64); ok {
		device.Port = int(port)
	}
	if brand, ok := updates["brand"].(string); ok {
		device.Brand = brand
	}
	if model, ok := updates["model"].(string); ok {
		device.Model = model
	}
	if serialNumber, ok := updates["serial_number"].(string); ok {
		device.SerialNumber = serialNumber
	}
	if apiKey, ok := updates["api_key"].(string); ok {
		device.APIKey = apiKey
	}
	if apiSecret, ok := updates["api_secret"].(string); ok {
		device.APISecret = apiSecret
	}
	if config, ok := updates["config"].(string); ok {
		device.Config = config
	}
	if remark, ok := updates["remark"].(string); ok {
		device.Remark = remark
	}

	return s.repo.Update(device)
}

func (s *FaceDeviceService) DeleteDevice(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("device not found")
	}
	return s.repo.Delete(id)
}

func (s *FaceDeviceService) EnableDevice(id int64) error {
	return s.repo.UpdateStatus(id, models.DeviceStatusOnline)
}

func (s *FaceDeviceService) DisableDevice(id int64) error {
	return s.repo.UpdateStatus(id, models.DeviceStatusDisabled)
}

func (s *FaceDeviceService) SetDeviceOffline(id int64) error {
	return s.repo.UpdateStatus(id, models.DeviceStatusOffline)
}

func (s *FaceDeviceService) SetDeviceFault(id int64) error {
	return s.repo.UpdateStatus(id, models.DeviceStatusFault)
}

func (s *FaceDeviceService) UpdateLastOnline(id int64) error {
	return s.repo.UpdateLastOnline(id)
}

func (s *FaceDeviceService) GetDeviceStatusSummary() (map[string]interface{}, error) {
	onlineCount, err := s.repo.CountByStatus(models.DeviceStatusOnline)
	if err != nil {
		return nil, err
	}

	offlineCount, err := s.repo.CountByStatus(models.DeviceStatusOffline)
	if err != nil {
		return nil, err
	}

	faultCount, err := s.repo.CountByStatus(models.DeviceStatusFault)
	if err != nil {
		return nil, err
	}

	disabledCount, err := s.repo.CountByStatus(models.DeviceStatusDisabled)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"online":   onlineCount,
		"offline":  offlineCount,
		"fault":    faultCount,
		"disabled": disabledCount,
		"total":    onlineCount + offlineCount + faultCount + disabledCount,
	}, nil
}

func (s *FaceDeviceService) GetOnlineDevices() ([]models.FaceDevice, error) {
	return s.repo.GetOnlineDevices()
}

// generateDeviceNo generates a unique device number
func (s *FaceDeviceService) generateDeviceNo() string {
	return fmt.Sprintf("FD%s%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)
}
