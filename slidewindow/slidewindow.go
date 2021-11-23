package slidewindow

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type SlideWindow struct {
	limit      int32
	intervalMs int
	count      int32
	index      int
	start      time.Time
	mu         sync.Mutex
	windows    []window
}

type window struct {
	count int32
}

func NewSlideWindwo(limit int32, intervalMs int) *SlideWindow {
	if intervalMs == 0 {
		intervalMs = 10
	}

	num := (1000 / intervalMs)
	return &SlideWindow{limit: limit, intervalMs: intervalMs, count: 0, index: 0, start: time.Now(), windows: make([]window, num)}
}

func (slidewindow *SlideWindow) CanVisit() bool {
	return atomic.LoadInt32(&slidewindow.count) < slidewindow.limit
}

func (slidewindow *SlideWindow) Visit() bool {
	now := time.Now()

	slidewindow.mu.Lock()

	diff := (int)(now.UnixMilli() - slidewindow.start.UnixMilli())
	if diff < 1000 {
		//一秒内，直接计算窗口索引
		slidewindow.index = diff / slidewindow.intervalMs
	} else if diff < 2000 {
		//超过一秒，但是未超过两秒，窗口开始时间向前推移一秒
		//同时，将过期的窗口的请求数清除
		slidewindow.start = slidewindow.start.Add(time.Second)
		slidewindow.index = (diff % 1000) / slidewindow.intervalMs
		for i := 0; i <= slidewindow.index; i++ {
			slidewindow.count -= slidewindow.windows[i].count
			slidewindow.windows[i].count = 0
		}
	} else {
		//超过2秒，则之前的所有请求都已经过期，全部清除，并且将窗口的开始时间向前推进相应的秒数
		slidewindow.index = (diff % 1000) / slidewindow.intervalMs
		seconds := (diff / 1000)
		slidewindow.start = slidewindow.start.Add(time.Duration(seconds) * time.Second)
		for i := range slidewindow.windows {
			slidewindow.windows[i].count = 0
		}
		slidewindow.count = 0
	}

	slidewindow.windows[slidewindow.index].count++
	slidewindow.count++

	slidewindow.mu.Unlock()

	if atomic.LoadInt32(&slidewindow.count) > slidewindow.limit {
		atomic.AddInt32(&slidewindow.count, -1)
		atomic.AddInt32(&slidewindow.windows[slidewindow.index].count, -1)
		fmt.Printf("reject request\n")
		return false
	}

	fmt.Printf("accept request\n")
	return true
}
