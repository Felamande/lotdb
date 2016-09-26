package query

import (
	"sync"
	"testing"

	"github.com/jinzhu/gorm"
)

func BenchmarkQueryNew(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q, err := Connect("sqlite3", "/home/kirigiri/GoPath/src/github.com/Felamande/lotdb/resource/db/lottery.2.sqlite3")
		if err != nil {
			b.Fatal(err)
		}
		q.Sum(45).Include(1, 2, 3).Exclude(14).Result()
	}
}

func BenchmarkQueryDefault(b *testing.B) {
	b.ResetTimer()
	db, err := gorm.Open("sqlite3", "/home/kirigiri/GoPath/src/github.com/Felamande/lotdb/resource/db/lottery.2.sqlite3")
	if err != nil {
		b.Fatal(err)
	}
	q := NewQuery(db)
	for i := 0; i < b.N; i++ {
		_, _ = q.Sum(45).Include(1, 2, 3).Exclude(14).Result()
	}
}

func BenchmarkQueryParallel(b *testing.B) {
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q, err := Connect("sqlite3", "/home/kirigiri/GoPath/src/github.com/Felamande/lotdb/resource/db/lottery.2.sqlite3")
			if err != nil {
				b.Fatal(err)
			}
			_, _ = q.Sum(45).Include(1, 2, 3).Exclude(14).Result()
		}
	})
}

func BenchmarkQueryGo40(b *testing.B) {
	b.ResetTimer()

	var wg sync.WaitGroup
	wg.Add(b.N)
	token := make(chan bool, 40)
	for i := 0; i < b.N; i++ {
		token <- true
		go func() {
			q, err := Connect("sqlite3", "/home/kirigiri/GoPath/src/github.com/Felamande/lotdb/resource/db/lottery.2.sqlite3")
			if err != nil {
				b.Fatal(err)
			}
			_, _ = q.Sum(45).Include(1, 2, 3).Exclude(14).Result()
			wg.Done()
			<-token
		}()
	}
	wg.Wait()
}
