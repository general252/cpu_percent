// +build linux darwin freebsd netbsd openbsd

package cpu_percent

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

func GetCpuPercent(interval time.Duration) (float64, error) {
	if data, err := cpu.Percent(interval, false); err != nil {
		return 0, err
	} else if len(data) > 0 {
		return data[0], nil
	} else {
		return 0, fmt.Errorf("cpu count is zero")
	}
}
