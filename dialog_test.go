package main

import (
	"math/rand"
	"testing"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomStringArray(count, length int) []string {
	arr := make([]string, count)
	for i := range count {
		arr[i] = RandomString(length)
	}
	return arr
}

func NewComboBoxItems_range(names []string) []*ComboBoxUintStruct {
	items := make([]*ComboBoxUintStruct, len(names))
	for i, n := range names {
		items[i] = &ComboBoxUintStruct{Enums: uint32(i), Name: n}
	}
	return items
}

func NewComboBoxItems_forindex(names []string) []*ComboBoxUintStruct {
	out := make([]*ComboBoxUintStruct, 0, len(names))
	for i := uint32(0); i < uint32(len(names)); i++ {
		out = append(out, &ComboBoxUintStruct{
			Enums: i,
			Name:  names[i],
		})
	}
	return out
}

var data []string = RandomStringArray(10000, 16)

func Benchmark_NewComboBoxItems1(b *testing.B) {
	for b.Loop() {
		NewComboBoxItems_range(data)
	}
}

func Benchmark_NewComboBoxItems2(b *testing.B) {
	for b.Loop() {
		NewComboBoxItems_forindex(data)
	}
}
