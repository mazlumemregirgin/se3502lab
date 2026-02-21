package week2

func main() {

	/*
		Fonksiyon yapısına baktığımızda Java ile farklılık gösterir
		Javada bir fonksiyondan birden fazla değer return edebilmek için bir class oluştururduk. Go'da ise fonksiyonlar birden fazla değer döndürebilirler.
		Go'da fonksiyonlar first-class citizens'dır. Bu, fonksiyonların değişkenlere atanabileceği, diğer fonksiyonlara argüman olarak geçirilebileceği ve fonksiyonlardan döndürülebileceği anlamına gelir.

		Go'da fonksiyon tanımlarken func keyword'ünü kullanırız. Fonksiyonun adı, parametreleri ve dönüş tipi belirtilir. Örneğin:

		// 1. func: Fonksiyon tanımladığını belirtir.
		// 2. Add: Fonksiyonun adı. (Büyük harfse Public, küçükse Private)
		// 3. (a int, b int): Parametreler ve tipleri.
		// 4. (int, error): Dönecek değerlerin tipleri.
		func Add(a int, b int) (int, error) {
			if a < 0 {
				// Hata durumunda: sonucun sıfır değerini ve hata mesajını dön.
				return 0, errors.New("negatif sayı girilemez")
			}
			// Başarılı durumda: sonucu ve nil (hata yok) dön.
			return a + b, nil
		}



		1. multiple return values: Go'da bir fonksiyon birden fazla değer döndürebilir.
		istediğim kadar değer döndürebilirim. ama best practice olarak 2-3'ü geçmemeye çalışırım. genellikle bir sonuç ve bir hata döndürürüm. Bu, Java'daki try-catch bloklarına benzer bir hata kontrol mekanizması sağlar.
		Örneğin:
		Java'da olsa: public float divide(float a, float b) throws Exception

		func divide(dividend, divisor float64) (float64, error) {
			if divisor == 0.0 {
				return 0.0, errors.New("cannot divide by zero") // Hata var, sonuç 0.0
			}
			return dividend / divisor, nil // Hata yok (nil), sonucu döndür
		}

		Javada bir metodun dönüş değerini kullanmasan da olur. Go'da ise eğer bir fonksiyon bir değer döndürüyorsa, o değeri kullanmak zorunda değilsin.
		Ancak, eğer fonksiyon birden fazla değer döndürüyorsa ve sen sadece birini kullanmak istiyorsan, diğerlerini _ (blank identifier) ile görmezden gelebilirsin. Örneğin:
		res, _ := divide(10, 2) // Hatayı görmezden geliyorum (Java'da try-catch'i boş bırakmak gibi, riskli!)


		2. named return values: Go'da fonksiyonların dönüş değerlerine isim verebilirsin. Bu, kodun okunabilirliğini artırır ve bazen hata ayıklamayı kolaylaştırır. Örneğin:
		func split(sum int) (x, y int) { // x ve y burada int(0) olarak doğdu bile
			x = sum * 4 / 9
			y = sum - x
			return // "x ve y'yi al ve git" demek
		}


		Java developerlar için farklar
		1. Public/Private: Java'daki public veya private keyword'leri Go'da yoktur. Fonksiyon ismini büyük harfle başlatırsan (Divide) her yerden erişilir (Public), küçük harfle başlatırsan (divide) sadece o paketten erişilir (Private).
		2. Overloading Yok: Java'daki gibi divide(int a) ve divide(int a, int b) şeklinde aynı isimli iki fonksiyon yazamazsın. Her fonksiyonun ismi benzersiz olmalıdır.
		3. Hata Kontrolü: Java'daki try-catch yerine Go'da sürekli şu kalıbı göreceksin:

			result, err := divide(10, 0)
			if err != nil {
				// Hatayı işle (logla, return et vb.)
			}

	*/

	/*
		POINTERS
		Go'da pointer'lar vardır ve C/C++'daki gibi çalışır. Bir pointer, bir değişkenin bellek adresini tutar. Go'da pointer'lar ile çalışmak için & (address of) ve * (dereference) operatörlerini kullanırız.,
		Java'da pointer'lar yoktur, ancak referans tipleri (örneğin, nesneler) benzer bir şekilde davranır. Go'da pointer'lar ile çalışmak, özellikle büyük veri yapılarıyla çalışırken performansı artırabilir.

		a. the opeators:
		& -> bir değişkenin adresini almak için kullanılır.
		* -> bir pointer'ı dereference etmek (yani, pointer'ın işaret ettiği değere erişmek) için kullanılır.

		b. pointer'ların avantajları:
		1. Bellek Verimliliği: Büyük veri yapılarıyla çalışırken, pointer'lar kopyalamak yerine veriye doğrudan erişim sağlar, bu da bellek kullanımını azaltır.
		2. Değişkenlerin Değiştirilmesi: Bir fonksiyona bir değişkenin pointer'ını geçirerek, fonksiyonun o değişkeni değiştirmesine izin verebilirsin. Bu, Java'daki referans tiplerine benzer bir davranış sağlar.
		3. Performans: Büyük veri yapılarıyla çalışırken, pointer'lar kopyalamak yerine veriye doğrudan erişim sağlar, bu da performansı artırabilir.

		c. declaration and usage:
		var p *int // p adında bir int pointer'ı tanımla
		x := 10
		p = &x
		fmt.Println("x'in değeri:", x) // 10
		fmt.Println("p'nin gösterdiği değer:", *p) // 10
		*p = 20
		fmt.Println("x'in yeni değeri:", x) // 20

	*/

	/*
		PASS-BY-VALUE SEMANTICS

		go her zaman pass-by-value dur. fonksiyona her zaman her şeyin kopyası gider. fark şu: neyin kopyası gittiği.
		int geçirirsen    →  değerin kopyası gider     (42'nin kopyası)
		pointer geçirirsen →  adresin kopyası gider    (0xc000016080'in kopyası)

		Java da pass-by-value'dur. Ancak, Java'da referans tipleri (örneğin, nesneler) pass-by-value olarak geçilirken, bu değerler aslında nesnenin kendisi değil, nesnenin bellekteki adresini tutan bir referanstır.
		Bu nedenle, Java'da bir nesne referansı geçildiğinde, fonksiyon o nesnenin içeriğini değiştirebilir çünkü referansın kopyası hala aynı nesneye işaret eder.

		Java:   referansın kopyası gider  →  aynı objeye bakılır
		Go:     adresin kopyası gider     →  aynı memory'e bakılır
		─────────────────────────────────────────────────────────
		Aynı şey. Farklı kelimeler.

	*/

	/*
		PACKAGE STRUCTURE

		fonksiyonların erişimi javadakinden daha basittir.
		Go'da bir fonksiyonun erişilebilir olup olmadığı, fonksiyonun adının büyük harfle başlayıp başlamadığına bağlıdır.
		Büyük harfle başlayan fonksiyonlar (örneğin, Add) diğer paketlerden erişilebilir (public), küçük harfle başlayan fonksiyonlar (örneğin, add) sadece tanımlandıkları paket içinde erişilebilir (private).

		paket yapısı-> go da her klasör bir pakettir. o klasörün içindeki tüm go dosyaları aynı pakete aittir. bu paketler birbirleriyle iletişim kurabilirler. java da ise package'lar klasör yapısına göre belirlenir ama go da klasör yapısı değil, package keyword'ü belirler.


		mathutils.Pi            // ✓ büyük P
		mathutils.Add(1, 2)     // ✓ büyük A
		mathutils.calculateSecret() // ✗ compiler error — küçük c
		```

		---

		## Özet
		```
		Java         →   Go
		──────────────────────
		public       →   BüyükHarf
		private      →   küçükHarf
		protected    →   yok
		package-priv →   yok
	*/
}
