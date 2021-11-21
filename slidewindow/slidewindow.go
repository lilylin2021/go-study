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
	mu         sync.Mutex
	windows    []window
}

type window struct {
	start int64
	count int32
}

func NewSlideWindwo(limit int32, intervalMs int) *SlideWindow {
	if intervalMs == 0 {
		intervalMs = 10
	}

	num := (1000 / intervalMs)
	return &SlideWindow{limit: limit, intervalMs: intervalMs, count: 0, index: 0, windows: make([]window, num)}
}

func (slidewindow *SlideWindow) CanVisit() bool {
	return atomic.LoadInt32(&slidewindow.count) < slidewindow.limit
}

func (slidewindow *SlideWindow) Visit() bool {
	now := time.Now().UnixMilli()

	slidewindow.mu.Lock()

	index := slidewindow.index % cap(slidewindow.windows)

	if slidewindow.count == 0 {
		slidewindow.windows[index].start = now
	} else {
		if (int)(now-slidewindow.windows[index].start) >= slidewindow.intervalMs {
			slidewindow.index = (int)(now-slidewindow.windows[index].start) / slidewindow.intervalMs
			index = slidewindow.index % cap(slidewindow.windows)
			if index == 0 {
				slidewindow.count = slidewindow.count - slidewindow.windows[index].count
				slidewindow.windows[index].count = 0
			}
			slidewindow.windows[index].start = now
		}
	}

	slidewindow.count++
	slidewindow.mu.Unlock()

	if atomic.LoadInt32(&slidewindow.count) > slidewindow.limit {
		atomic.AddInt32(&slidewindow.count, -1)
		fmt.Printf("reject request\n")
		return false
	}
	atomic.AddInt32(&slidewindow.windows[index].count, 1)
	fmt.Printf("accept request\n")
	return true
}
