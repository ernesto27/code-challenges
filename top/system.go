package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type SystemInfo struct {
	Uptime string
}

func NewSystemInfo() *SystemInfo {
	return &SystemInfo{}
}

func (s *SystemInfo) GetUptime() error {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return fmt.Errorf("unexpected format in /proc/uptime")
	}

	uptimeSeconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return err
	}

	totalSeconds := int(uptimeSeconds)
	days := totalSeconds / (24 * 3600)
	hours := (totalSeconds % (24 * 3600)) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if days > 0 {
		s.Uptime = fmt.Sprintf("%d days, %02d:%02d:%02d", days, hours, minutes, seconds)
	} else {
		s.Uptime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return nil
}
