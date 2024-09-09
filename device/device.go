package device

import "fmt"

type device struct {
	Username  string `json:"username"`
	ID        int    `json:"id"`
	IP        string `json:"ip"`
	Logged_in bool   `json:"logged_in"`
}

func NewDevice(username string, id int, ip string) *device {
	return &device{
		Username:  username,
		ID:        id,
		IP:        ip,
		Logged_in: false,
	}
}

type DeviceManager struct {
	Devices       map[int]*device      // 主存储，key 是 Device ID
	ipIndex       map[string]int       // IP 到 Device ID 的索引
	usernameIndex map[string][]*device // Username 到 Device 切片的索引
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		Devices:       make(map[int]*device),
		ipIndex:       make(map[string]int),
		usernameIndex: make(map[string][]*device),
	}
}

func (dm *DeviceManager) AddDevice(device *device) error {
	// 检查 ID 是否已存在
	if _, exists := dm.Devices[device.ID]; exists {
		return fmt.Errorf("device with ID %d already exists", device.ID)
	}

	// 检查 IP 是否已存在
	if _, exists := dm.ipIndex[device.IP]; exists {
		return fmt.Errorf("device with IP %s already exists", device.IP)
	}

	// 添加到主存储
	dm.Devices[device.ID] = device

	// 添加到 IP 索引
	dm.ipIndex[device.IP] = device.ID

	// 添加到 Username 索引
	dm.usernameIndex[device.Username] = append(dm.usernameIndex[device.Username], device)

	return nil
}

func (dm *DeviceManager) GetDeviceByID(id int) (*device, bool) {
	device, exists := dm.Devices[id]
	return device, exists
}

func (dm *DeviceManager) GetDeviceByIP(ip string) (*device, bool) {
	if id, ok := dm.ipIndex[ip]; ok {
		return dm.GetDeviceByID(id)
	}
	return nil, false
}

func (dm *DeviceManager) GetDevicesByUsername(username string) []*device {
	return dm.usernameIndex[username]
}

// 新增：通过 Username 和 Device ID 获取特定设备
func (dm *DeviceManager) GetDeviceByUsernameAndID(username string, id int) (*device, bool) {
	devices := dm.usernameIndex[username]
	for _, device := range devices {
		if device.ID == id {
			return device, true
		}
	}
	return nil, false
}

// 新增：删除设备的方法
func (dm *DeviceManager) RemoveDevice(id int) error {
	device, exists := dm.Devices[id]
	if !exists {
		return fmt.Errorf("device with ID %d not found", id)
	}

	// 从主存储中删除
	delete(dm.Devices, id)

	// 从 IP 索引中删除
	delete(dm.ipIndex, device.IP)

	// 从 Username 索引中删除
	devices := dm.usernameIndex[device.Username]
	for i, d := range devices {
		if d.ID == id {
			dm.usernameIndex[device.Username] = append(devices[:i], devices[i+1:]...)
			break
		}
	}

	return nil
}
