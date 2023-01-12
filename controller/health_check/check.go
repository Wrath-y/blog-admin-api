package health_check

import (
	"blog-admin-api/core"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// HealthCheck shows `OK` as the ping-pong result.
func HealthCheck(c *core.Context) {
	c.Success(nil)
}

// DiskCheck checks the disk usage.
func DiskCheck(c *core.Context) {
	u, _ := disk.Usage("/")

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	text := "OK"

	if usedPercent >= 95 {
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	c.Success(message)
}

// CPUCheck checks the cpu usage.
func CPUCheck(c *core.Context) {
	cores, _ := cpu.Counts(false)
	usedPercent, _ := cpu.Percent(time.Second, false)

	a, _ := load.Avg()
	l1 := a.Load1
	l5 := a.Load5
	l15 := a.Load15

	text := "OK"

	if l5 >= float64(cores-1) {
		text = "CRITICAL"
	} else if l5 >= float64(cores-2) {
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Load average: %.2f, %.2f, %.2f | Cores: %d | Used %f%%", text, l1, l5, l15, cores, usedPercent)
	c.Success(message)
}

// RAMCheck checks the disk usage.
func RAMCheck(c *core.Context) {
	u, _ := mem.VirtualMemory()

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	text := "OK"

	if usedPercent >= 95 {
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	c.Success(message)
}
