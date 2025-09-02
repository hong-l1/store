package apiserver

import (
	"github.com/IBM/sarama"
	"math/rand"
	"sync"
	"time"
)

type Fn struct {
}

func (f *Fn) Handle(msg *sarama.ConsumerMessage, t string) error {
	f.ListenHeartbeat(t)
	return nil
}
func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func (fn *Fn) ListenHeartbeat(server string) {
	time := time.Now()
	mutex.Lock()
	dataServers[server] = time
	mutex.Unlock()
}
func RemoveExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}
func NewFn() *Fn {
	return &Fn{}
}
