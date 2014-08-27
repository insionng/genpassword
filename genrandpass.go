package main

/*
   密码字典生成 insion@live.com
*/

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"os"
)

var (
	cache        = make([][]byte, 0)
	mutex        sync.RWMutex
	isTimeout    bool
	maxCacheSize = 100
	firstTime    time.Time
	timeOut      = 20 * time.Second

	CharsSet = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+`-=[]\\{}|;':\",./<>? ")

	ipath        = "./password.txt" //存储字典的路径
	ilen         = 4                //生成的密码长度
	imark uint32 = 0
)

func CacheInsert(path string, data []byte) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	if len(cache) == 0 {
		firstTime = time.Now()
		isTimeout = false
	} else {
		isTimeout = time.Now().Sub(firstTime) > timeOut
	}

	cache = append(cache, data)
	if len(cache) >= maxCacheSize || isTimeout {

		buf := [][]byte(cache)

		fs, e := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if e != nil {
			return e
		}
		defer fs.Close()

		for _, v := range buf {
			fs.WriteString(string(v))
			fs.WriteString(string("\n"))
		}

		cache = make([][]byte, 0)
		return err
	}
	return nil
}

func gen(charset string, n int, sc chan string) {

	for i, c := range charset {
		if n == 1 {
			sc <- string(c)
		} else {
			var ssc = make(chan string)
			go gen(charset[:i]+charset[i+1:], n-1, ssc)
			for k := range ssc {
				sc <- fmt.Sprintf("%v%v", string(c), k)
			}
		}
	}
	close(sc)

}

func main() {

	starttime := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())

	sc := make(chan string)

	/*
		go gen(string(CharsSet), ilen, sc)
		for x := range sc {
			atomic.AddUint32(&imark, 1)
			CacheInsert(ipath, []byte(x))
			fmt.Println("Gen:", x)
		}
	*/
	go gen(string(CharsSet), ilen, sc)

	fs, e := os.OpenFile(ipath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if e != nil {
		panic(e)
	}
	defer fs.Close()

	for x := range sc {
		atomic.AddUint32(&imark, 1)
		fs.WriteString(x)
		fs.WriteString(string("\n"))
		fmt.Println("Gen:", x)
	}

	imarkFinal := atomic.LoadUint32(&imark)
	since := int(time.Since(starttime).Seconds())

	for {
		fmt.Println("完成消耗时间:", since, "s", "生成:", imarkFinal, "个密码")
	}
}
