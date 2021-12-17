package tccEstimator

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWrapper(t *testing.T) {
	sequence_number_unwrapper := NewSequenceNumberUnwrapper()
	i := int64(0)
	for i < 65536+65536 {
		assert.Equal(t, i, sequence_number_unwrapper.Unwrap(uint16(i)))
		//fmt.Println( sequence_number_unwrapper.Unwrap(uint16(i)) )
		i++
	}
	//fmt.Println(sequence_number_unwrapper.Unwrap(1))
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "收到信号，监控退出,time=", time.Now().Unix())
			return
		default:
			fmt.Println(name, "goroutine监控中,time=", time.Now().Unix())
			time.Sleep(1 * time.Second)
		}
	}
}

func TestContext(t *testing.T){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	go watch(ctx, "监控1")
	go watch(ctx, "监控2")

	//fmt.Println("现在开始等待8秒,time=", time.Now().Unix())
	time.Sleep(1 * time.Second)

	//fmt.Println("等待8秒结束,准备调用cancel()函数，发现两个子协程已经结束了，time=", time.Now().Unix())
	cancel()

	time.Sleep(1*time.Second)
}
