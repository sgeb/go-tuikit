package main

import (
	"time"

	"fmt"
	"os"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/sgeb/go-tuikit/tuikit/binding"
)

//----------------------------------------------------------------------------
// Cpu
//----------------------------------------------------------------------------

type Cpu struct {
	User    binding.Uint64Property
	Nice    binding.Uint64Property
	Sys     binding.Uint64Property
	Idle    binding.Uint64Property
	Wait    binding.Uint64Property
	Irq     binding.Uint64Property
	SoftIrq binding.Uint64Property
	Stolen  binding.Uint64Property

	cpu  sigar.Cpu
	stop chan struct{}
}

func NewCpu() *Cpu {
	return &Cpu{
		User:    binding.NewUint64Property(),
		Nice:    binding.NewUint64Property(),
		Sys:     binding.NewUint64Property(),
		Idle:    binding.NewUint64Property(),
		Wait:    binding.NewUint64Property(),
		Irq:     binding.NewUint64Property(),
		SoftIrq: binding.NewUint64Property(),
		Stolen:  binding.NewUint64Property(),
		stop:    make(chan struct{}),
	}
}

func (c *Cpu) get() error {
	if err := c.cpu.Get(); err != nil {
		return err
	}

	c.User.Set(c.cpu.User)
	c.Nice.Set(c.cpu.Nice)
	c.Sys.Set(c.cpu.Sys)
	c.Idle.Set(c.cpu.Idle)
	c.Wait.Set(c.cpu.Wait)
	c.Irq.Set(c.cpu.Irq)
	c.SoftIrq.Set(c.cpu.SoftIrq)
	c.Stolen.Set(c.cpu.Stolen)

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

//----------------------------------------------------------------------------
// Memory
//----------------------------------------------------------------------------

type Memory struct {
	Total      binding.Uint64Property
	Used       binding.Uint64Property
	Free       binding.Uint64Property
	ActualFree binding.Uint64Property
	ActualUsed binding.Uint64Property
}
