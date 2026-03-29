package main

import (
	"fmt"
	"sync"
)

/*
burada sync paketindeki wait group özelliğini kullnacağız.

wait group aslında isminden de anlaşılacağı üzere birden fazla çağırdığımız go keywordu ile fonkisyonlarını bir kuyruğa sokarak bunlar bitmeden main fonksiyonunun bitmemesini sağlar.
*/

func hellowithwaitgroup(wg *sync.WaitGroup) {

	// defer keywordu ile bu fonksiyon çalışıp bittikten sonra sayacımızı bir azaltacağını garanti altına almalıyız. wg.Done() ile "sayaç -1" yapıyoruz aslında
	// bu wg.Done() aslında fonkisyonun sonunda çalışır fonkisyon bitince fakat best practice olarak fonksiyon başında defer keywordu ile belirtmemiz bunun yapılacağını
	// garanti altına almamızı sağlar. javadaki finally keywordu gibi düşünebiliriz.
	defer wg.Done()

	fmt.Println("wait grorup ile merhaba diyorum")
}

func waitgroup() {

	// 1.ADIM = sync.WaitGroup fonksiyonunu ana fonksiyonumuza eklemeliyiz. bunu var wg sync.WaitGroup() ile yapıyoruz.
	// bu yazım şekli best-practice'tir. bu bizim kayıt defterimiz gibi düşünebiliriz.
	var wg sync.WaitGroup

	// 2.ADIM  = kayıt deferimize içine kaç fonkisyon sıraya sokacaksak bunu belirtmeliyiz ve deftere eklemeliyiz. bunu ise wg.Add(1) ile yapıyoruz.
	wg.Add(1)

	// 3.ADIM = kayıt defterimizi yani wg!yi adresiyle beraber go keywordu ile çalıştıracağımız fonksyona parametre olarak göndermeliyiz.
	go hellowithwaitgroup(&wg)

	// 4.ADIM = burada ise kayıt defterimizdeki sayaç 0 olana kadar ana fonkisyonumuzun bitmemesi için bit bekleme fonksiyonu eklemeliyiz wg.Wait() ile.
	// sayaç sıfır olana kadar bekleyecek ana fonksiyon yani aslında az önce manuel olarak 500ms bekleme time.Sleep() eklemek yerine bunu da dinamik hale getiriyoruz.
	wg.Wait()

	fmt.Println("wait group ile main çalıştı main bitti")
}
