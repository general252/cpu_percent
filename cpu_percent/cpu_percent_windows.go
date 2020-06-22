package cpu_percent

import (
	"golang.org/x/sys/windows"
	"log"
	"time"
	"unsafe"
)

/*
//#include <pdh.h>
typedef unsigned long DWORD;
// Union specialization for double values
typedef struct _PDH_FMT_COUNTERVALUE_DOUBLE {
	DWORD  CStatus;
	double DoubleValue;
} PDH_FMT_COUNTERVALUE_DOUBLE;
*/
import "C"

var (
	modPdh = windows.NewLazySystemDLL("pdh.dll")

	pdhOpenQuery                = modPdh.NewProc("PdhOpenQuery")
	pdhAddCounter               = modPdh.NewProc("PdhAddCounterW")
	pdhCollectQueryData         = modPdh.NewProc("PdhCollectQueryData")
	pdhGetFormattedCounterValue = modPdh.NewProc("PdhGetFormattedCounterValue")
	pdhCloseQuery               = modPdh.NewProc("PdhCloseQuery")
)

const (
	PDH_FMT_DOUBLE   = 0x00000200
	PDH_INVALID_DATA = 0xc0000bc6
	PDH_NO_DATA      = 0x800007d5
)

// createQuery XXX
// copied from https://github.com/mackerelio/mackerel-agent/
func createQuery() (windows.Handle, error) {
	var query windows.Handle
	r, _, err := pdhOpenQuery.Call(0, 0, uintptr(unsafe.Pointer(&query)))
	if r != 0 {
		return 0, err
	}
	return query, nil
}

// createCounter XXX
func createCounter(query windows.Handle, cname string) (windows.Handle, error) {
	var counter windows.Handle
	r, _, err := pdhAddCounter.Call(
		uintptr(query),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(cname))),
		0,
		uintptr(unsafe.Pointer(&counter)))
	if r != 0 {
		return windows.InvalidHandle, err
	}
	return counter, nil
}

// getCounterValue get counter value from handle
func getCounterValue(counter windows.Handle) (float64, error) {
	var value C.PDH_FMT_COUNTERVALUE_DOUBLE
	r, _, err := pdhGetFormattedCounterValue.Call(uintptr(counter), PDH_FMT_DOUBLE, uintptr(0), uintptr(unsafe.Pointer(&value)))
	if r != 0 && r != PDH_INVALID_DATA {
		return 0.0, err
	}
	return float64(value.DoubleValue), nil
}

// collectQueryData
func collectQueryData(query windows.Handle) {
	r, _, err := pdhCollectQueryData.Call(uintptr(query))
	if r != 0 && err != nil {
		if r == PDH_NO_DATA {
			log.Printf("this metric has not data. %v", err)
			return
		} else {
			log.Println(err)
		}
		return
	}
}

// closeQuery
func closeQuery(query windows.Handle) {
	_, _, _ = pdhCloseQuery.Call(uintptr(query))
}

// GetCpuPercent
func GetCpuPercent(interval time.Duration) (float64, error) {
	query, err := createQuery()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	counter, err := createCounter(query, "\\Processor Information(_Total)\\% Processor Utility")
	if err != nil {
		log.Println(err)
		return 0, err
	}

	collectQueryData(query)
	time.Sleep(interval)
	collectQueryData(query)

	data, err := getCounterValue(counter)

	closeQuery(query)

	return data, err
}

func test() {
	for {
		percent, _ := GetCpuPercent(time.Second)
		log.Println(percent)
	}
}
