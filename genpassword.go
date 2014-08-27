package main

/*
   密码字典生成 insion@live.com
*/

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
	"os"
)

var (
	CharsSet = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+`-=[]\\{}|;':\",./<>? ")

	ipath        = "./password.txt" //存储字典的路径
	ilen         = 3                //生成的密码长度
	imark uint32 = 0
)

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

	fs, e := os.OpenFile(ipath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if e != nil {
		panic(e)
	}
	defer fs.Close()

	for i := 1; i <= ilen; i++ {
		fmt.Println("i:", i)
		sc := make(chan string)

		go gen(string(CharsSet), i, sc)

		for x := range sc {
			atomic.AddUint32(&imark, 1)
			fs.WriteString(x)
			fs.WriteString(string("\n"))
			fmt.Println("Gen:", x)
		}
	}

	imarkFinal := atomic.LoadUint32(&imark)
	since := int(time.Since(starttime).Seconds())

	fmt.Println("完成消耗时间:", since, "s", "生成:", imarkFinal, "个密码")
	time.Sleep(10)
}
