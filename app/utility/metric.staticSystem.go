package utility

import (
	"github.com/mackerelio/go-osstat/cpu"
	//"github.com/mackerelio/go-osstat/disk"
	"github.com/mackerelio/go-osstat/memory"
)

// statisResultStruct là cấu trúc chứa kết quả thống kê
type statisResultStruct struct {
	Percent float64     `json:"percent" bson:"percent"` // Phần trăm sử dụng
	Static  interface{} `json:"static" bson:"static"`   // Thông tin chi tiết
}

// GetMemoryStatic lấy thông tin thống kê bộ nhớ
func GetMemoryStatic() interface{} {
	memory, err := memory.Get()
	if err != nil {
		// Nếu có lỗi, trả về nil
		//fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}

	// Tạo kết quả thống kê
	result := new(statisResultStruct)
	result.Static = memory
	result.Percent = float64(memory.Used) / float64(memory.Total) * 100
	return result
}

// GetCpuStatic lấy thông tin thống kê CPU
func GetCpuStatic() interface{} {
	cpu, err := cpu.Get()
	if err != nil {
		// Nếu có lỗi, trả về nil
		//fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}

	// Tạo kết quả thống kê
	result := new(statisResultStruct)
	result.Static = cpu
	result.Percent = float64(cpu.Nice+cpu.System+cpu.User) / float64(cpu.Total) * 100
	return result

	return cpu
}

// GetDiskStatic lấy thông tin thống kê đĩa (đã bị comment)
// func GetDiskStatic() interface{} {
// 	disk, err := disk.Get()
// 	if err != nil {
// 		// Nếu có lỗi, trả về nil
// 		//fmt.Fprintf(os.Stderr, "%s\n", err)
// 		return nil
// 	}

// 	return disk
// }
