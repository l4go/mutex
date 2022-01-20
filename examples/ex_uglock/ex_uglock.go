package main

import (
	"log"
	"time"
	"sync"

	"github.com/l4go/mutex"
)

const (
	WORKERS = 5
)

func main() {
	data := 10
	wg := new(sync.WaitGroup)
	w := WORKERS

	m := mutex.NewUgMutex()
	for w > 0 {
		wg.Add(1)
		go writer(m, wg, &data)

		wg.Add(1)
		go reader(m, wg, &data)
		wg.Add(1)
		go reader(m, wg, &data)
		wg.Add(1)
		go reader(m, wg, &data)

		wg.Add(1)
		go upgradeWriter(m, wg, &data)

		w--
	}
	wg.Wait()
}

func writer(m sync.Locker, wg *sync.WaitGroup, data *int) {
	defer wg.Done()

	m.Lock()
	defer m.Unlock()

	log.Println("writer data : ", *data)
	*data += 1
	log.Println("writer added : ", *data)
}

func reader(m *mutex.UgMutex, wg *sync.WaitGroup, data *int) {
	defer wg.Done()

	m.RLock()
	defer m.RUnlock()

	log.Println("reader data : ", *data)
}

func upgradeWriter(m *mutex.UgMutex, wg *sync.WaitGroup, data *int) {
	defer wg.Done()

	m.UgLock()
	defer m.UgUnlock()

	log.Println("upgrade writer data : ", *data)

	time.Sleep(1000)

	m.Upgrade()
	log.Println("upgraded !!!  : ", *data)
	*data += 1
	log.Println("upgrade writer added : ", *data)
}
