package main

import (
	"fmt"
	// "os"
	"runtime"
	"runtime/debug"
)

var global *int

func print_of_gcstats(g *debug.GCStats) string {
	return fmt.Sprintf("last_gc_time: %v, num_gc: %v, pause_total: %v, pause_history: %v, pause_end: %v, pause_quantities: %v", g.LastGC, g.NumGC, g.PauseTotal, g.Pause, g.PauseEnd, g.PauseQuantiles)
}

func print_memstats(m *runtime.MemStats) {
	fmt.Printf("alloc_heap: %v, total_alloc_cumu: %v, sys_from_os: %v, lookups: %v, mallocs: %v, frees: %v, heap_objects: %v\n", m.Alloc, m.TotalAlloc, m.Sys, m.Lookups, m.Mallocs, m.Frees, m.HeapObjects)
}

func f() {
	memstats := new(runtime.MemStats)
	// gcstats := new(debug.GCStats)
	// debug.ReadGCStats(gcstats)
	// fmt.Printf("Inside f() and at the beginning: %s\n", print_of_gcstats(gcstats))
	fmt.Printf("Inside f() and at the beginning:\n")
	// debug.WriteHeapDump(os.Stdout.Fd())
	runtime.ReadMemStats(memstats)
	print_memstats(memstats)
	var x int
	x = 1
	global = &x
	fmt.Printf("Inside f() and at the end:\n")
	runtime.ReadMemStats(memstats)
	print_memstats(memstats)
	// debug.WriteHeapDump(os.Stdout.Fd())
	// debug.ReadGCStats(gcstats)
	// fmt.Printf("Inside f() and at the end: %s\n", print_of_gcstats(gcstats))
}

func g() {
	memstats := new(runtime.MemStats)
	// gcstats := new(debug.GCStats)
	// debug.ReadGCStats(gcstats)
	// fmt.Printf("Inside g() and at the beginning: %s\n", print_of_gcstats(gcstats))
	fmt.Printf("Inside g() and at the beginning:\n")
	// debug.WriteHeapDump(os.Stdout.Fd())
	runtime.ReadMemStats(memstats)
	print_memstats(memstats)
	y := new(int)
	*y = 1
	fmt.Printf("Inside g() and at the end:\n")
	runtime.ReadMemStats(memstats)
	print_memstats(memstats)
	// debug.WriteHeapDump(os.Stdout.Fd())
	// debug.ReadGCStats(gcstats)
	// fmt.Printf("Inside g() and at the end: %s\n", print_of_gcstats(gcstats))
}

func main() {
	fmt.Printf("global in main (before f()) is: %p\n", global)
	fmt.Printf("Before we call f()\n")
	memstats := new(runtime.MemStats)
	runtime.ReadMemStats(memstats)
	print_memstats(memstats)
	// gcstats := new(debug.GCStats)
	// debug.ReadGCStats(gcstats)
	// fmt.Printf("Before we call f(): %s\n", print_of_gcstats(gcstats))
	// debug.WriteHeapDump(os.Stdout.Fd())
	f()
	fmt.Printf("after call f()\n")
	// debug.WriteHeapDump(os.Stdout.Fd())
	// fmt.Printf("after we call f(): %s\n", print_of_gcstats(gcstats))
	runtime.ReadMemStats(memstats)
	print_memstats(memstats)
	g()
	fmt.Printf("after call g()\n")
	runtime.ReadMemStats(memstats)
	print_memstats(memstats)
	// debug.WriteHeapDump(os.Stdout.Fd())
	// fmt.Printf("after we call f() and g(): %s\n", print_of_gcstats(gcstats))
	fmt.Printf("global in main (after f()) is: %p\n", global)
}
