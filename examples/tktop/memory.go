package main

import "github.com/sgeb/go-tuikit/tuikit/binding"

type Memory struct {
	Total      binding.Uint64Property
	Used       binding.Uint64Property
	Free       binding.Uint64Property
	ActualFree binding.Uint64Property
	ActualUsed binding.Uint64Property
}
