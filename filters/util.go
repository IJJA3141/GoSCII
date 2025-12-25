package filters

import (
	"cmp"
	"runtime"
	"sync"
)

func split(lines int, lambda func(start, end int)) *sync.WaitGroup {
	var wg sync.WaitGroup
	var split int
	var cpus = runtime.GOMAXPROCS(0)

	if cpus >= lines {
		for line := range lines {
			wg.Go(func() { lambda(line, line+1) })
		}
	} else {
		if cpus == 1 {
			split = lines
		} else {
			split = int(lines / (cpus - 1))
		}

		for cpu := range cpus {
			wg.Go(func() { lambda(cpu*split, (cpu+1)*split) })
		}
	}

	return &wg
}

// clamp constrains a value to lie within the inclusive range [lower, upper].
//
// If value is less than lower, clamp returns lower.
// If value is greater than upper, clamp returns upper.
// Otherwise, it returns value unchanged.
func clamp[T cmp.Ordered](value, lower, upper T) T { return max(lower, min(upper, value)) }
