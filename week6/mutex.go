package main

import (
	"fmt"
	"sync"
)

/*

Bir banka codebaseinde bir para yatırma fonksiyonunu geliştirdiğimizi varsayalım.
go goroutine kapsamında ama aslında go routine dışında genel bir sorun var bu para yatırma işlemlerinde:
aynı anda birden fazla kullanıcının para yatırma veya para çekme işlemi yaptığını düşünelim.

Bakiye:100 ve kullanıcı_1 10 tl para yatırmak istiyor. ve bakiyenin 110 tl olması bekleniyor. aynı anda kullanıcı_2 de para yatırmak istiyor ama bunlar aynı anda işlemi gerçekleştirmeye
çalıştıkları için anlık durumları göremiyorlar bu yüzden kullanıcı_2 de bakiyeyi 100tl olarak görüyor ve para yatırınca bakiye 110tl olarak gözüküyor

aslında 10+10 = 20tl yatırıldı fakat bu akış güvenli ve tutarlı şekilde ilerlemediği için sistemde bakiye doğru güncellenemedi ve bakiye 100tl gözüküyor.const

Bu sorunu MUTEX ile çözebiliriz.
*/

// önce banka hesabı structını oluşturdum içinde bizim için aynı anda yapılan işlemleri biri yaparken diğeriinin işlem yapmasını engelleyecek olan mutex tipinde mu'yu ekledim.
type BankaHesabı struct {
	mu     sync.Mutex
	bakiye int
}

func (b *BankaHesabı) ParaYatır(miktar int, wg *sync.WaitGroup) {
	defer wg.Done()

	// burada sync.Mutex'in Lock() fonksiyonunu kullanarak işlem başlangıcında başka birinin daha bu işleme aynı aynda girmesini engelliyoruz. fonksiyonu kitliyoruz aslında.
	b.mu.Lock()

	fmt.Println("para yatıyor.")
	b.bakiye += miktar
	// içeride logic işlemler yapıldıktan sonra ise tekrardan sync.Mutex'in Unlock() fonksiyonunu kullanarak kilidi açıyoruz. kilit açılınca diğer go worker gelip yeniden çalıştırabiliyor bu fonkisyonu
	b.mu.Unlock()
}

/*
Bu sayede kullanıcı_1 geldi para yatırma işlemi yaparken fonkisyon mutex ile kilitlendi önce kullanıcı_1in işleminin bitmesini bekledi. kullanıcı_1'in işlemi bitince
kullanıcı_2 fonksiyonu kullandı. sorun çözüldü artık hesapta bakiye 110tl değil 120tl olarak gözüküyor.
*/

/*
aslında go routine ile işleri asenkron hale getirmeye çalışırken burada kısa süreliğinde işlemleri sequential'a dnüştürüyoruz. ama bu aşamada bunu yapmalıyız çünkü burada da bir trade-off var
verinin tutatlılığı (data integrity) yanlış verinin gösterilmesi hiç veri gelmemesinden daha kötü. bu işlemin daha sağlıklı olmasını sağlamak için mutex kullanmalıyız
*/

// her yeri lock() ile kitlememliyiz çünkü ne kadar fazla alanı locklarsak diğer goroutineler o kadar çok bekler. sadece goroutinelerin ortak veriye dokunduğu yeri locklamak gerekiyor.
func mutexlifonks() {

	var wg sync.WaitGroup
	hesap := BankaHesabı{bakiye: 0}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go hesap.ParaYatır(1, &wg)
	}

	wg.Wait()
	fmt.Println("Son Bakiye:", hesap.bakiye)
}
