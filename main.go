package main

import (
	"time"
	"bytes"
	"context"
	"flag"
	"log"
	"os"
)

var block_size = int64(512)
var null_block = make([]byte, 2*block_size)

func readAndSplit(ctx context.Context, cancel context.CancelFunc, worker_pool chan func(), block_device *os.File, start, end int64) {
	if end-start < 2*block_size {
		log.Printf("Finished searching at block %+v", start)
		return
	}
	block := make([]byte, 2*block_size)
	location := (end + start) / 2
	location = location - (location % block_size) // align the read on the block
	select {
	case <-ctx.Done():
		return
	default:
	}
	if _, err := block_device.ReadAt(block, location); err != nil {
		log.Printf("Error encountered with start %d and end %d : %+v", start, end, err)
		cancel()
	} else {
		if !bytes.Equal(block, null_block) {
			log.Printf("Got contentful block: %+v", block)
			cancel()
		} else {
			log.Printf("Checking blocks at address %+v with width %+v", location, end-start)
			worker_pool <- func() { readAndSplit(ctx, cancel, worker_pool, block_device, start, location) }
			worker_pool <- func() { readAndSplit(ctx, cancel, worker_pool, block_device, location, end) }
		}
	}
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	worker_pool := make(chan func(), 1000000000)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case f := <-worker_pool:
					f()
				case  <- time.After(10 * time.Second):
				 log.Printf("No more work, assuming we have verified zero-content of disk")
				 cancel()
				}
			}
		}()
	}
	if block_device, err := os.Open(flag.Arg(0)); err != nil {
		log.Fatalf("%s", err)
	} else {
		defer block_device.Close()
		if end, err := block_device.Seek(0, 2); err != nil {
			log.Fatalf("%s", err)
		} else {
			log.Printf("Found end of disk: %+d bytes", end)
			readAndSplit(ctx, cancel, worker_pool, block_device, 0, end - 2* block_size)
		}
	}
	<-ctx.Done()
}
