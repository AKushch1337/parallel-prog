package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	CrystalSize     = 30
	NumParticles    = 50
	SimulationSteps = 1000
)

type Particle struct {
	id  int
	pos int
}

func main() {
	fmt.Println("Сценарій 1: Без синхронізації")
	runSimulation(false)

	fmt.Println("\nСценарій 2: З використанням Mutex")
	runSimulation(true)
}

func runSimulation(useSync bool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	particles := make([]*Particle, NumParticles)
	crystal := make([]int, CrystalSize)

	for i := 0; i < NumParticles; i++ {
		pos := r.Intn(CrystalSize)
		particles[i] = &Particle{id: i, pos: pos}
		crystal[pos]++
	}

	fmt.Printf("Початковий стан: %d частинок\n", countParticles(crystal))
	printCrystal(crystal)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < NumParticles; i++ {
		wg.Add(1)
		go func(p *Particle) {
			defer wg.Done()

			localRand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(p.id)))

			for step := 0; step < SimulationSteps; step++ {
				dir := -1
				if localRand.Intn(2) == 1 {
					dir = 1
				}

				if useSync {
					mu.Lock()
				}

				oldPos := p.pos
				newPos := oldPos + dir

				if newPos < 0 {
					newPos = 0
				} else if newPos >= CrystalSize {
					newPos = CrystalSize - 1
				}

				crystal[oldPos]--
				if !useSync {
					time.Sleep(10 * time.Microsecond)
				}
				crystal[newPos]++
				p.pos = newPos

				if useSync {
					mu.Unlock()
				}
			}
		}(particles[i])
	}

	wg.Wait()

	finalCount := countParticles(crystal)
	fmt.Printf("Кінцевий стан:   %d частинок\n", finalCount)
	printCrystal(crystal)

	if finalCount != NumParticles {
		fmt.Printf("Помилка: кількість частинок %d замість %d (race condition)\n", finalCount, NumParticles)
	} else {
		fmt.Println("Успіх: кількість частинок не змінилась")
	}
}

func countParticles(crystal []int) int {
	sum := 0
	for _, count := range crystal {
		sum += count
	}
	return sum
}

func printCrystal(crystal []int) {
	var builder strings.Builder
	builder.WriteString("[")
	for i, count := range crystal {
		if count == 0 {
			builder.WriteString(" . ")
		} else {
			builder.WriteString(fmt.Sprintf("%2d ", count))
		}
		if i < len(crystal)-1 {
			builder.WriteString("|")
		}
	}
	builder.WriteString("]")
	fmt.Println(builder.String())
}
