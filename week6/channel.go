package main

import (
	"fmt"
	"time"
)

// Don't communicate by sharing memory; share memory by communicating

/*
go routinelerimizde iki ayrı go keywordu ile başlattığımız fonkisyonun birbirileriyel haberleşmesini, birbirine veri gönderme işini nasıl yapacağız??
bu sıkıntı birşey çünkü go routinler birbirinden bağımsız halde çalışan birbirlerinden habersiz fonksiyonlardır.
bunları ancak channel dediğimiz yapılar ile bir arada iletişimini sağlayabiliriz.


*/

// Burası bizim öncelik olarak ikinci goya gidecek verimizi düzenleyecek başka bir fonksiyon.
// önce verimiz buraya girecek sonrasında düzenlene ana fonkisyonumuzda oluşturdğumuz channela aktarılıp ikinci fonksiyona geçecek.
func processData(channel chan string) {
	for i := 1; i <= 3; i++ {
		urun := fmt.Sprintf("Hammadde-%d", i)
		fmt.Println("Üretici: ", urun, " hazırlandı.")

		// oluşturulan ürünü, veriyi ortak kanala fırlatıyoruz.
		channel <- urun

		time.Sleep(time.Second)
	}
	// işimiz bittiğinde bu kanalı kapatmalıyız ki diğer fonkisyon sonsuza kadar beklemesin. çünkü diğer fonkisyon channel açık olduğu sürece sanki bir veri daha gelecekmiş gibi bekler kapanmaz.
	close(channel)
}

// burası ise go keywordu ile başlattığımız ikinci fonkisyonumuz. ilk fonksiyondan düzenlenmiş verileri ortak kanaldan channeldan çekerek başka bir işlem yapar.
func packegeIt(channel chan string, done chan bool) {
	// channel kapanana kadar kanaldan gelen verileri okuması için for döngüsüne sokuyoruz. kanaldan veri geldiği takdirde hemen işlemi yapacak.
	for urun := range channel {
		fmt.Println("Montajcı: ", urun, " paketleniyor")
		time.Sleep(time.Second * 2)
		fmt.Println("Montajcı: ", urun, " PAKETLENDİ")
	}
	// Tüm paketleme bittiğinde ana programa "tamam" diyoruz
	done <- true
}

func chanellifonks() {

	// iki ayrı go routineinin birbiriyle ilişki kurabilmesi için ortak bir hat kuruyoruz bunu make(chan string ile yapıyoruz)
	transferChannel := make(chan string)

	// bu da en son programın bitip bitmediğini anlamamızı sağlayan kanal için olacak.
	isFinished := make(chan bool)

	// iki workerı da burada başlatıyoruz. birbirileirin tetikleyecekler kanallardan iletişim kuracaklar.
	go processData(transferChannel)
	go packegeIt(transferChannel, isFinished)

	// ana fonksiyonumuz isFinished kanalından true mesajı gelene kadar burada bekleyecek.
	<-isFinished
	fmt.Println("Fabrika kapandı, tüm işler bitti.")
}

// bu kod unbuffered channelin bir implementasyonu. YA ÜRETİCİ ÇOK HIZLI PAKETLEYİCİ ÇOK YAVAŞSA??

/*
Üretici bir ürün yapar, montajcı onu elinden alana kadar kapıda bekler. Yani ikisi de en yavaş olanın hızına mahkumdur.


buffered channel ile bu sorunu çözebiliriz:
make(chan string, 10) dersek, araya bir "kasa/stok" koymuş oluruz. Üretici 10 tane ürünü hızlıca yapıp kasaya bırakır ve işine bakabilir. Montajcı ise arkadan yavaş yavaş o kasadan ürünleri çeker.

transferHatti := make(chan string, 10) satırının yanına: "Buradaki 10 rakamı, alıcı veriyi almasa bile üreticinin 10 adet veriyi kanala fırlatıp beklemeden yoluna devam edebileceği kapasiteyi temsil eder."


*/
