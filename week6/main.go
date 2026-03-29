package main

import (
	"fmt"
	"time"
)

func helloWorld() {
	fmt.Println("merhaba")
}

func main() {
	go helloWorld()

	time.Sleep(500) // bu satırı koymazsak main fonksiyonu bittiğinde goroutine de kapanır. go keydowrdu ile çağırdığımız fonksiyon bitemeden main biteceği için kapanır çalışmaz.
	// ama normalde bunu amele gibi manuel olarak eklemek bu ölçekte bir kod için normal ama gerçek hayat senaryoalına baktığımızda mesela Fintech alanında bir şirkette milisaniyelerin
	// bile çok önemi var. bu yüzden bu bekleme süresini bizim manuel olarak eklemek yerine bunu go tarafına bırakabiliyoruz. bunu da wait grouplar ile yapıyoruz.

	fmt.Println("go fonksiyonu çalıştı")
	waitgroup() // wait group ile yapılan implementasyon için waitgroup.go dosyasına git.

	// peki mesela bir banka hesabı kodu yazdığımızı düşünelim burada veri persistencelığını nasıl sağlayacağız??
	// bunu da sync.Mutex ile sağlayacağız.
	mutexlifonks()

	// channel ile ilgili implementasyonu görmek için channel.go dosyasına git.
	chanellifonks()
}
