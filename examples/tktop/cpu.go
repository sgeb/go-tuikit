package main

import (
	"time"

	"fmt"
	"os"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/sgeb/go-tuikit/tuikit/binding"
)

type Cpu struct {
	User binding.Float32Property
	Sys  binding.Float32Property
	Idle binding.Float32Property

	cpu sigar.Cpu

	lastUser  uint64
	lastSys   uint64
	lastIdle  uint64
	lastTotal uint64

	stop chan struct{}
}

func NewCpu() *Cpu {
	return &Cpu{
		User: binding.NewFloat32Property(),
		Sys:  binding.NewFloat32Property(),
		Idle: binding.NewFloat32Property(),
		stop: make(chan struct{}),
	}
}

func (c *Cpu) get() error {
	if err := c.cpu.Get(); err != nil {
		return err
	}

	user := c.cpu.User
	sys := c.cpu.Sys
	idle := c.cpu.Idle
	total := c.cpu.Total()

	diffUser := float32(user - c.lastUser)
	diffSys := float32(sys - c.lastSys)
	diffIdle := float32(idle - c.lastIdle)
	diffTotal := float32(total - c.lastTotal)

	c.User.Set(100.0 * diffUser / diffTotal)
	c.Sys.Set(100.0 * diffSys / diffTotal)
	c.Idle.Set(100.0 * diffIdle / diffTotal)

	c.lastUser = user
	c.lastSys = sys
	c.lastIdle = idle
	c.lastTotal = total

	return nil
}

func (c *Cpu) Start(interval time.Duration) {
	c.get()
	go func() {
		ticker := time.Tick(interval)
		for {
			select {
			case <-ticker:
				fmt.Fprintln(os.Stderr, "tick")
				c.get()
			case <-c.stop:
				return
			}
		}
	}()
}

func (c *Cpu) Stop() {
	c.stop <- struct{}{}
}
