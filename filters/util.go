package filters

import (
	"runtime"
	"sync"
)

type lambda func(_start, _end int)

func Split(_Δ int, _fc lambda) *sync.WaitGroup {
	var wg sync.WaitGroup
	var split int

	cpus := runtime.GOMAXPROCS(0)
	wg.Add(cpus)

	if cpus != 0 {
		split = int(_Δ / (cpus - 1))
	} else {
		split = _Δ
	}

	for cpu := range cpus {
		go func() { defer wg.Done(); _fc(cpu*split, (cpu+1)*split) }()
	}

	return &wg
}

func assert(condition bool) {
	if !condition {
		panic(1)
	}
}
