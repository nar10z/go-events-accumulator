package collector_example

import (
	"context"
	"fmt"
	"testing"
	"time"

	acc "github.com/lrweck/accumulator"
	gocoll "github.com/nar10z/go-collector"
)

const (
	flushSize     = 1000
	flushInterval = time.Second
)

type Data struct {
	i int
}

func Benchmark_accum(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	b.ResetTimer()
	b.Run("#1.1 go-collector, channel", func(b *testing.B) {
		summary := 0

		collector, _ := gocoll.New[*Data](flushSize, flushInterval, func(events []*Data) error {
			summary += len(events)
			time.Sleep(time.Microsecond)
			return nil
		})

		for i := 0; i < b.N; i++ {
			_ = collector.AddAsync(ctx, &Data{i: i})
		}

		collector.Stop()

		fmt.Printf("#1. summary=%d\n", summary)
		if summary != b.N {
			b.Fail()
		}
	})
	b.Run("#1.2 go-collector, list", func(b *testing.B) {
		summary := 0

		collector, _ := gocoll.NewWithStorage[*Data](flushSize, flushInterval, func(events []*Data) error {
			summary += len(events)
			time.Sleep(time.Microsecond)
			return nil
		}, gocoll.List)

		for i := 0; i < b.N; i++ {
			_ = collector.AddAsync(ctx, &Data{i: i})
		}

		collector.Stop()

		fmt.Printf("#1. summary=%d\n", summary)
		if summary != b.N {
			b.Fail()
		}
	})
	b.Run("#1.3 go-collector, slice", func(b *testing.B) {
		summary := 0

		collector, _ := gocoll.NewWithStorage[*Data](flushSize, flushInterval, func(events []*Data) error {
			summary += len(events)
			time.Sleep(time.Microsecond)
			return nil
		}, gocoll.Slice)

		for i := 0; i < b.N; i++ {
			_ = collector.AddAsync(ctx, &Data{i: i})
		}

		collector.Stop()

		fmt.Printf("#1. summary=%d\n", summary)
		if summary != b.N {
			b.Fail()
		}
	})
	b.Run("#1.4 go-collector, stdList", func(b *testing.B) {
		summary := 0

		collector, _ := gocoll.NewWithStorage[*Data](flushSize, flushInterval, func(events []*Data) error {
			summary += len(events)
			time.Sleep(time.Microsecond)
			return nil
		}, gocoll.StdList)

		for i := 0; i < b.N; i++ {
			_ = collector.AddAsync(ctx, &Data{i: i})
		}

		collector.Stop()

		fmt.Printf("#1. summary=%d\n", summary)
		if summary != b.N {
			b.Fail()
		}
	})

	b.Run("#2. lrweck/accumulator", func(b *testing.B) {
		summary := 0

		inputChan := make(chan *Data, flushSize)
		batch := acc.New(inputChan, flushSize, flushInterval)

		go func() {
			for i := 0; i < b.N; i++ {
				inputChan <- &Data{i: i}
			}
			close(inputChan)
		}()

		_ = batch.Accumulate(ctx, func(o acc.CallOrigin, items []*Data) {
			summary += len(items)
			time.Sleep(time.Microsecond)
		})

		fmt.Printf("#2. summary=%d\n", summary)
		if summary != b.N {
			b.Fail()
		}
	})
}
