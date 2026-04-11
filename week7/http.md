# Go `net/http` — Server-Side Ders Notları

> **Kapsam:** HTTP istek/yanıt döngüsü, `http.Handler` arayüzü, Go 1.22+ `ServeMux`, ve Fiber ile karşılaştırma.

---

## Öğrenme Hedefleri

Bu dersin sonunda şunları yapabileceksin:

- HTTP Request/Response yaşam döngüsünü düşük seviyeli bir ağ perspektifinden açıklamak
- Go standart kütüphanesini (`net/http`) kullanarak production-ready bir HTTP sunucusu kurmak
- Go 1.22+ ile gelen gelişmiş `ServeMux` özelliklerini (metot eşleştirme ve wildcard) uygulamak
- `http.Handler` arayüzünü analiz ederek standart kütüphane ile Fiber gibi framework'ler arasındaki mimari farkı anlamak

---

## 1. Teorik Temel: HTTP Döngüsü

Kod yazmadan önce bir veri paketinin yolculuğunu zihnimizde canlaştırmalıyız. Fiber gibi üst düzey framework'ler kullandığımızda bu süreç genellikle gizlenir. Ham döngüyü anlamak, gecikme ve bağlantı sorunlarını ayıklamak için kritiktir.

### HTTP Request/Response Yaşam Döngüsü

1. **DNS Çözümlemesi:** İstemci (tarayıcı/API tüketicisi) alan adını bir IP adresine çözer.
2. **TCP El Sıkışması:** 3 yönlü el sıkışma (`SYN`, `SYN-ACK`, `ACK`) güvenilir bir bağlantı kurar.
3. **TLS El Sıkışması (HTTPS):** Güvenliyse şifreleme anahtarları müzakere edilir.
4. **İstek:** İstemci biçimlendirilmiş bir metin bloğu (Headers + Body) gönderir.
5. **Sunucu İşleme (Go'nun Rolü):**
   - İşletim sistemi çekirdeği bağlantıyı kabul eder.
   - Go'nun `net/http` sunucusu bağlantıyı kabul eder ve her istek için bir **Goroutine** başlatır.
   - İstek bir `http.Request` struct'ına ayrıştırılır.
6. **Yanıt:** Sunucu, `http.ResponseWriter` aracılığıyla byte'ları sokete geri yazar.
7. **Sonlandırma:** Bağlantı kapatılır veya sonraki istekler için canlı tutulur (`Keep-Alive`).

> **Araştırma Notu:** Go'nun standart kütüphanesi, kullanıma hazır haliyle production-grade olması açısından benzersizdir. Python veya Ruby'nin aksine, Go'nun `net/http`'si Cloudflare'in edge node'ları gibi devasa sistemlere güç vermektedir.

---

## 2. Temel Arayüz: `http.Handler`

Tüm Go web ekosistemi tek bir basit arayüz etrafında döner. Bunu anlarsan, Go web sunucularını anlarsın.

```go
type Handler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}
```

### İki Temel Bileşen

**`w http.ResponseWriter`** — Yanıtı oluşturmak için kullanılan bir arayüz.

- **Katı Kural:** Body'yi yazmadan önce header'ları yazmalısın. Body'nin ilk byte'ını yazdığında, header'lar temizlenir ve artık değiştirilemez.

Bu arayüzün içinde 3 temel yetenek (metot) bulunur:
- `Header()` — Yanıtın başlıklarını (Content-Type vb.) ayarlar.
- `WriteHeader(statusCode)` — "200 OK" mi yoksa "404 Not Found" mu döneceğiz?
- `Write([]byte)` — Gövdeye (Body) asıl veriyi basar.

**`r *http.Request`** — Tüm istek verilerini (URL, Method, Headers, Body) içeren bir struct pointer'ı.

İstekle ilgili her şeyi buradan okuruz:
- `r.Method` — "POST" mu "GET" mi?
- `r.URL.Path` — Hangi adrese gelmiş?
- `r.Body` — Gelen JSON nerede?
- `r.Header` — Client hangi tarayıcıyı kullanıyor?

> **Neden Pointer?** Bir HTTP isteği çok büyük olabilir (header'lar, devasa bir JSON body, çerezler vb.). Eğer pointer yerine kopyasını alsaydık, her istekte bellekte gereksiz bir yük oluştururduk.

### Örnek 1: "Bare Metal" Handler

Bu örnek, herhangi bir multiplexer (router) olmayan bir sunucuyu göstermektedir. Porta gelen *tüm* istekleri işler.

```go
package main

import (
    "fmt"
    "net/http"
)

// MyHandler, http.Handler arayüzünü uygular
type MyHandler struct{}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 1. Header'ları ayarla
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // 200 OK

    // 2. Body'yi yaz
    // Fprint, ResponseWriter'a yazmak için kullanılır (io.Writer'dır)
    fmt.Fprintf(w, `{"status": "received", "path": "%s"}`, r.URL.Path)
}

func main() {
    handler := &MyHandler{}

    fmt.Println("Server starting on :8080...")
    // ListenAndServe sonsuza dek bloklar
    if err := http.ListenAndServe(":8080", handler); err != nil {
        panic(err)
    }
}
```

---

## 3. `ServeMux` ile Yönlendirme (Go 1.22+ Güncellemesi)

Tarihsel olarak Go'nun standart router'ı (`ServeMux`) çok sınırlıydı; geliştiriciler HTTP metotlarını (GET vs POST) veya path parametrelerini işlemek için üçüncü taraf kütüphanelere (Gorilla Mux veya Chi gibi) başvurmak zorundaydı.

**Go 1.22 (Şubat 2024) itibarıyla bu değişti.** Standart kütüphane artık şunları desteklemektedir:

- **Metot Eşleştirme:** `GET /path`
- **Wildcard'lar:** `/items/{id}`

### `ServeMux` Kavramı

`ServeMux` özünde devasa bir `switch` ifadesidir. Aynı zamanda bir `http.Handler`'dır! Bir isteği alır, URL'ye bakar ve isteği doğru özel handler'a iletir.

### Neden `http.HandleFunc` Değil de `mux := http.NewServeMux()`?

Yeni başlayanlar genelde doğrudan `http.HandleFunc()` kullanır. Bu, Go'nun arka planda yarattığı `DefaultServeMux` (varsayılan yönlendirici) üzerine yazar. `http.NewServeMux()` daha iyidir çünkü:

- **İzolasyon:** Uygulamanın farklı kısımları için farklı router'lar tanımlayabilirsin.
- **Güvenlik:** Global bir değişken (`DefaultServeMux`) kullanmak büyük projelerde beklenmedik çakışmalara yol açabilir.
- **Test Edilebilirlik:** `mux` bir nesne olduğu için test fonksiyonlarına parametre olarak gönderebilirsin.

> **Fiber Karşılaştırması:** Fiber'de `app := fiber.New()` ile oluşturduğun `app` de aslında bir multiplexer'dır. Fiber bize bu işlemi kolaylaştırır, ancak burada manuel olarak yapıyoruz.

### Multiplexer (Mux) Nedir?

"Mux" ismi Multiplexer'dan gelir. Elektronikte birden fazla sinyali tek bir çıkışa yönlendiren cihazdır. Web dünyasında:

- **Girdi:** Binlerce farklı URL ve HTTP metodu (GET, POST vb.)
- **İşlem:** Gelen isteği inceleyip "Bu istek hangi fonksiyona gitmeli?" kararını vermek
- **Çıktı:** Doğru Handler fonksiyonunu çalıştırmak

| Özellik | Fiber (`app := fiber.New()`) | Go Standart (`mux := http.NewServeMux()`) |
|---|---|---|
| Yönlendirme | `app.Get("/user", handler)` | `mux.HandleFunc("GET /user", handler)` |
| Parametre Alımı | `c.Params("id")` | `r.PathValue("id")` (Go 1.22+) |
| Middleware | `app.Use(logger)` | Biraz daha manuel (Sarmalama mantığı) |

### `ListenAndServe` Parametrelerinin Sırrı

`ListenAndServe(addr string, handler Handler)` iki şey ister:

- **`addr` (":9000"):** Hangi adresi ve portu dinleyeceği. Başındaki iki nokta (`:`) şu anlama gelir: "Bu makinedeki tüm ağ arayüzlerini (localhost, yerel IP vb.) dinle."
- **`handler` (mux):** `http.NewServeMux()` ile oluşturduğun `mux` nesnesi, içinde `ServeHTTP` metodunu barındırır. Go'da bir nesne bir arayüzün metoduna sahipse, o arayüzün ta kendisi sayılır. Yani `mux` bir `Handler`'dır.

**Fiber Karşılaştırması:** `app.Listen(":3000")` dediğinde Fiber arkada kendi server motorunu (fasthttp) çalıştırır.

> **"Blocking" (Bloklama) Mevzusu:** Bu fonksiyon blocking bir fonksiyondur. Kod buraya geldiğinde durur ve dinlemeye başlar. Eğer alt satıra geçerse, bu genellikle server'ın çöktüğü anlamına gelir. Bu yüzden genellikle şöyle kullanılır:

```go
err := http.ListenAndServe(":9000", mux)
if err != nil {
    log.Fatal(err) // Hata varsa logla ve uygulamayı kapat
}
```

### İstek Gelince Ne Oluyor? (The Magic)

Birisi `localhost:9000/register` adresine istek attığında:

1. `ListenAndServe` bu isteği yakalar.
2. Hemen yeni bir **Goroutine** (hafif thread) açar.
3. İsteği (`r`) ve yanıt yazıcısını (`w`) paketleyip `mux.ServeHTTP(w, r)` fonksiyonunu çağırır.
4. `mux` da kendi içine bakar: "/register için hangi fonksiyonu kaydetmiştin?" der ve senin yazdığın `registerHandler`'ı çalıştırır.

### Örnek 2: Modern Standart Kütüphane Routing

```go
package main

import (
    "fmt"
    "net/http"
    "strconv"
)

func main() {
    mux := http.NewServeMux()

    // 1. Metota Özgü Routing (Go 1.22'de Yeni)
    // Yalnızca GET istekleriyle eşleşir. POST istekleri otomatik olarak 405 alır.
    mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("System OK"))
    })

    // 2. Path Wildcard'ları (Go 1.22'de Yeni)
    mux.HandleFunc("GET /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        idStr := r.PathValue("id")

        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid ID format", http.StatusBadRequest)
            return
        }

        response := fmt.Sprintf("Fetching details for User ID: %d", id)
        w.Write([]byte(response))
    })

    // 3. POST için catch-all
    mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Creating a new user..."))
    })

    fmt.Println("Starting Server on :8080")
    // mux'ı handler olarak geçiriyoruz
    // ServeMux, ServeHTTP'yi uyguladığı için arayüze uyar
    http.ListenAndServe(":8080", mux)
}
```

---

## 4. Veri Modeli ve Bellek (Storage)

Backend'de her şey veridir. Fiber'de olduğu gibi önce verinin "şeklini" belirleyen bir struct yazmalıyız.

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Global bir "database" simülasyonu (Thread-safe değil ama şu anlık yeterli)
var userRegistry = make(map[string]User)
```

> **Fiber Karşılaştırması:** Fiber'de veriyi genelde `c.BodyParser(&user)` ile alırsın. Burada ise `encoding/json` paketini manuel kullanacağız.

### Neden Map ve Struct Kullanıyoruz?

Gerçek bir dünyada (PostgreSQL veya MongoDB kullanırken), veriler diskte saklanır. Ancak öğrenme aşamasında veya hızlı prototiplemede Go'nun kendi veri yapılarını kullanmak bize hız kazandırır.

- **`struct` (Şema):** Veritabanındaki bir tablo (table) veya döküman (document) yapısıdır. Verinin hangi alanlardan oluşacağını belirler.
- **`map` (Koleksiyon):** Veritabanının kendisidir. Key-Value yapısı sayesinde bir kullanıcıyı ismiyle veya ID'siyle aradığımızda $O(1)$ karmaşıklığında bulunmasını sağlar.

### Kritik Fark: Kalıcılık (Persistence)

Fiber veya Spring ile bir uygulama yazdığında, uygulamayı kapatıp açarsan veritabanın (DB) hariciyse veriler kalır. Ancak `map` kullanırsak:

- **Dezavantaj:** Uygulama (process) durdurulduğu an tüm RAM temizlenir ve kullanıcılar silinir.
- **Avantaj:** Hiçbir dış bağımlılığın (PostgreSQL kurulumu, Docker vb.) olmaz, sadece saf Go koduyla mantığı kurarsın.

### "Concurrency" (Eşzamanlılık) Uyarısı

Go'daki standart `map` yapısı **"thread-safe" değildir**.

> **Sorun:** Aynı anda iki farklı kişi `/register` isteği atarsa (yani iki farklı Goroutine aynı map'e yazmaya çalışırsa), Go "concurrent map writes" hatası verip uygulamayı çökertecektir.

> **Çözüm:** Gerçek projelerde ya bir veritabanı kullanılır ya da map'in yanına bir `sync.Mutex` (kilit mekanizması) eklenir.

---

## 5. JSON Decoding: Byte'lardan Struct'a

### HTTP İsteği Nedir? (The Raw Reality)

İstemci, sunucuna veriyi bir karakter dizisi (raw text) olarak gönderir. TCP bağlantısı üzerinden gelen şey aslında şuna benzer:

```
POST /register HTTP/1.1
Content-Type: application/json

{"name": "emre", "email": "emre@mail.com"}
```

Bilgisayar için bu sadece bir byte yığınıdır. Go bu byte'ları otomatik olarak bir struct'a dönüştürmez; çünkü gelen verinin ne olduğunu (JSON mu, XML mi, yoksa düz metin mi) bilemez.

### `r.Body` Nedir? (The Stream)

`r.Body` aslında bir **Musluktur (Stream)**. Veri internetten bir kerede küt diye gelmez; byte byte akar. Go'da buna `io.ReadCloser` denir:

- **Read:** İçinden veri okunabilir.
- **Closer:** Okuma bitince musluğun kapatılması gerekir (standart kütüphane bunu bizim yerimize yapar).

### `json.NewDecoder(r.Body)` — The Translator

Burada şunu diyoruz: "Ey Go! Sana bir musluk (`r.Body`) veriyorum. Bu musluktan akan byte'ların JSON formatında olduğunu biliyorum. Sen bu musluğun başına otur ve gelenleri dinlemeye başla."

`NewDecoder` bir mekanizma kurar; henüz hiçbir şeyi çözmemiştir, sadece musluk ile senin aranda bir "çevirmen" görevlendirmiştir.

### `.Decode(&u)` — The Filling

İşte sihir burada gerçekleşiyor:

1. `Decode` metodu musluğu açar ve byte'ları okumaya başlar.
2. Okuduğu her bir JSON anahtarını ("name", "email"), senin verdiğin `&u` (User struct'ı) içindeki alanlarla eşleştirir.

> **Neden `&u` (Pointer)?** Çünkü `Decode` fonksiyonunun, fonksiyon dışındaki bir değişkenin içini doldurabilmesi (mutate etmesi) gerekir. Eğer pointer vermezsen, fonksiyon kendi içinde bir kopyayla çalışır ve senin asıl `u` değişkenin boş kalır.

### Struct Tag'lerinin Gücü

`User` struct'ını tanımlarken şöyle yazarız:

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

`json:"name"` kısımları (bunlara **Struct Tags** denir), `Decode` fonksiyonuna şu talimatı verir: "Gelen JSON içinde küçük harfle 'name' diye bir anahtar görürsen, onu al ve benim büyük harfle başlayan `Name` alanıma yaz."

> JSON'da genelde küçük harf (camelCase), Go'da ise export edilebilir olması için büyük harf (PascalCase) kullanıldığı için bu "eşleştirme" işlemi kritiktir.

### Verinin Kablodan Struct'a Yolculuğu (Kare Kare)

1. **Girdi:** `[123, 34, 110, 97, 109, 101, ...]` — Ham byte dizisi `{ " n a m e ...`
2. **Decoder:** Byte'ları tek tek okur, süslü parantezi görür ("Hah, bir nesne başlıyor" der).
3. **Mapping:** "name" anahtarını bulur, struct'ındaki `json:"name"` etiketli alanı arar.
4. **Atama:** Bulduğu değeri ("emre") o alanın içine yazar.
5. **Sonuç:** `u.Name` yazınca "emre" cevabını veren bir Go değişkeni elde edilir.

> **Fiber'de Neden Bunu Düşünmüyorsun?** Fiber'de `c.BodyParser(&u)` dersin. Fiber arkada şuna bakar: "Header'da `Content-Type: application/json` var mı? Varsa, go encoding/json kütüphanesini çağır." Saf Go'da bu çeviri işleminin bizzat sorumlusu sensin.

> **Neden manuel daha verimli olabilir?** `NewDecoder`, tüm body'nin belleğe yüklenmesini beklemeden akış (stream) üzerinden veriyi okuyabilir. Bu, çok büyük JSON'larda Fiber'in bazı yöntemlerinden daha az RAM tüketir.

---

## 6. Tam Proje Kodu: `main.go`

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// User modeli
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Basit in-memory storage
var userRegistry = make(map[string]User)

func main() {
    // 1. Multiplexer (Router) tanımlıyoruz
    mux := http.NewServeMux()

    // 2. Route tanımları (Go 1.22+ syntax)
    mux.HandleFunc("POST /register", registerHandler)
    mux.HandleFunc("GET /user/{name}", getUserHandler)

    // 3. Sunucuyu başlatma
    fmt.Println("Server 9000 portunda çalışıyor...")
    if err := http.ListenAndServe(":9000", mux); err != nil {
        fmt.Printf("Hata: %s\n", err)
    }
}

// registerHandler: Kullanıcıyı kaydeder
func registerHandler(w http.ResponseWriter, r *http.Request) {
    var u User

    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, "JSON çözülemedi", http.StatusBadRequest)
        return
    }

    if u.Name == "" {
        http.Error(w, "İsim boş olamaz", http.StatusUnprocessableEntity)
        return
    }

    userRegistry[u.Name] = u

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(u)
}

// getUserHandler: Path parametresinden ismi alır ve kullanıcıyı döner
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")

    user, ok := userRegistry[name]
    if !ok {
        http.Error(w, "Kullanıcı yok", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

---

## 7. Handler'ları Adım Adım Anlamak

### `registerHandler` — Adım Adım Akış

#### Adım 1: Hazırlık
Gelen veriyi içine dökeceğimiz boş bir `User` struct'ı oluştururuz.

#### Adım 2: Okuma (Decoding)
`r.Body` içindeki ham byte'ları JSON'dan `User` struct'ına çözeriz.

```go
if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
    http.Error(w, "JSON çözülemedi", http.StatusBadRequest)
    return
}
```

`if err := ... ; err != nil` — Go'nun meşhur hata kontrolüdür. Eğer gelen veri geçerli bir JSON değilse (örneğin süslü parantez kapatılmadıysa), `Decode` fonksiyonu bir hata döner ve biz 400 Bad Request göndeririz.

#### Adım 3: İş Mantığı (Logic)
Çözdüğümüz kullanıcıyı `userRegistry` map'ine ekleriz.

```go
userRegistry[u.Name] = u
```

"Eğer bunu yapmazsan, fonksiyon bittiği an `u` değişkeni yok olur (çöpe gider). Map'e kaydederek, bir sonraki GET isteğinde o kullanıcıyı isminden bulabilmeyi garanti ediyoruz."

#### Adım 4: Yanıt
Kullanıcıya başarı mesajı döneriz.

```go
w.Header().Set("Content-Type", "application/json")  // Mektubun üzerine etiket yapıştırmak
w.WriteHeader(http.StatusCreated)                    // Mührü basmak (201)
json.NewEncoder(w).Encode(u)                         // Yanıtı paketleyip göndermek
```

**`w.Header().Set("Content-Type", "application/json")`** — İstemciye ne gönderdiğimizi söyleriz. Fiber'de `c.JSON()` dediğinde, Fiber bu satırı senin yerine otomatik olarak yazar.

**`w.WriteHeader(http.StatusCreated)`** — HTTP dünyasının evrensel dili. `http.StatusCreated` aslında 201 sayısıdır.
- `200 OK`: "Her şey yolunda."
- `201 Created`: "Her şey yolunda ve senin istediğin şeyi başarıyla oluşturdum."
- **Kritik Kural:** Bunu `w.Write()` veya `Encode()` yapmadan hemen önce yapmalısın. Body yazılmaya başladıktan sonra status kodu değiştiremezsin.

**`json.NewEncoder(w).Encode(u)`** — `Decoder` işleminin tam tersi (Encoding / Marshaling). `json.NewEncoder(w)` çıktıyı doğrudan internet kablosuna yazar; `.Encode(u)` ise `u` nesnesini JSON'a çevirir ve `w` üzerinden istemciye gönderir.

### Büyük Resim: Verinin "U-Dönüşü"

```
GİRİŞ:   İstemciden JSON geldi  →  Decoder: Go Struct'ına çevirdi
İŞLEM:   Go Struct'ı Map'e kaydedildi
ÇIKIŞ:   Go Struct'ı tekrar JSON'a çevrildi (Encoder)  →  İstemciye geri gitti
```

> **Neden geri gönderiyoruz?** Bu bir backend standardıdır. İstemciye "Bak, gönderdiğin veriyi aldım ve bu şekilde kaydettim, her şey yolunda" teyidini vermektir.

### Fiber Karşılaştırması: `registerHandler`

```go
// Fiber ile yazılsaydı:
func registerHandler(c *fiber.Ctx) error {
    var u User
    if err := c.BodyParser(&u); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON çözülemedi"})
    }
    if u.Name == "" {
        return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "İsim boş olamaz"})
    }
    userRegistry[u.Name] = u
    return c.Status(fiber.StatusCreated).JSON(u) // header + status + encode tek satırda
}
```

Fiber'de `w` (yazıcı) ve `r` (okuyucu) diye iki ayrı nesneyle uğraşmak yerine, bunlar tek bir Context (`c`) objesi içinde birleştirilir.

---

### `getUserHandler` — Path Parametresi ile Kullanıcı Getirme

#### Adım 1: Path Parametresini Yakalamak

```go
name := r.PathValue("name")
```

Go 1.22'nin yeni özelliği ile URL'deki `{name}` kısmına ne yazıldıysa (örneğin "emre") onu bir string olarak alırız.

> **Fiber Karşılaştırması:** `c.Params("name")`

#### Adım 2: Map'ten Sorgulama — "Comma ok" Idiom

```go
user, ok := userRegistry[name]
```

- `user` — Eğer "emre" map'te varsa, onun `User` struct'ını getirir.
- `ok` — Bir boolean (true/false). Eğer "emre" kayıtlıysa `true`, yoksa `false` döner.

> Go `map`'ten olmayan bir key istendiğinde hata vermez; sadece boş bir değer döner ve `ok` değişkeni `false` olur. Bu yüzden `if !ok` ile kullanıcıyı bulamadığımızı kontrol ederiz.

#### Adım 3: Bulamazsak — Hata Dönme

```go
if !ok {
    http.Error(w, "Kullanıcı yok", http.StatusNotFound)
    return
}
```

`http.Error` çok pratiktir: hem `w.WriteHeader(404)` yapar hem de body'ye hata mesajını yazar.

#### Adım 4: Bulursak — Encode ve Gönder

```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(user)
```

> **Önemli Farkındalık:** Fiber veya Spring'de bir nesne döndüğünde (`return user`), framework arkada Encode işlemini senin için yapar. Saf Go'da ise "Explicit over Implicit" (Açıklık kapalılıktan iyidir) kuralı gereği, her şeyi adım adım senin yazman istenir. Bu başta zahmetli gelse de, hata aldığında tam olarak nerede sorun olduğunu şak diye anlarsın.

---

## 8. `net/http` ile Interface Gücü: Composition

`json.NewEncoder(w)` fonksiyonu parametre olarak `io.Writer` ister. Bizim `http.ResponseWriter` bir arayüzdür ve içinde `Write` metodu barındırdığı için, JSON encoder doğrudan HTTP yanıtına veri yazabilir. **Bu Go'nun Composition gücüdür.**

**Hata Yönetimi:** Framework'lerde hatalar genelde `return c.Status(500).SendString("err")` ile yönetilir. Saf Go'da `http.Error(w, ...)` fonksiyonu hem status kodunu set eder hem de body'ye hata mesajını yazar; ardından fonksiyonu `return` ile kesmek gerekir.

---

## 9. Kritik Analiz: `net/http` vs. Fiber

Fiber, `net/http`'ye değil, farklı bir motor olan **`fasthttp`**'ye dayanır.

| Özellik | Standart `net/http` | Fiber (fasthttp) |
|---|---|---|
| **Bellek Modeli** | Her istek/yanıt için yeni nesneler tahsis eder (daha güvenli, debug edilmesi daha kolay). | **Zero-Allocation:** Bellek tasarrufu için nesneleri yeniden kullanır. Son derece hızlı, ancak referanslar çok uzun tutulursa race condition'a yol açabilir. |
| **API Stili** | Düşük seviyeli, açık. Manuel JSON decode/encode gerektirir. | **Express.js stili.** `c.JSON()` veya `c.BodyParser()` gibi yüksek seviyeli metotlar. |
| **Uyumluluk** | Go ekosistemiyle %100 uyumlu. | Standart `http.Handler` middleware'iyle **uyumsuz** (adaptör gerektirir). |
| **Context** | Standart `context.Context`'i yoğun şekilde kullanır. | Kendi context mekanizmasını kullanır (bridge'ler mevcut). |

### Fiber Bizim İçin Neyi Kolaylaştırıyor?

**1. Context (`c *fiber.Ctx`) Kavramı**
`net/http`'de `w` (yazıcı) ve `r` (okuyucu) diye iki ayrı nesneyle uğraşıyordun. Fiber bunları tek bir `Context` (`c`) objesi içinde birleştirir.

**2. BodyParser (Akıllı Çevirmen)**
Standart kütüphanede `json.NewDecoder(r.Body).Decode(&u)` yazarken Fiber'de `c.BodyParser(&u)` dersin. Fiber sadece JSON değil, `form-data` veya `xml` olarak gelse de otomatik olarak ayrıştırır.

**3. Tek Satırda Cevap (Method Chaining)**
Saf Go'da 3-4 adımda yaptığın şeyi Fiber'de zincirleme yapabilirsin:
```go
return c.Status(201).JSON(u)
// Bu satır: Header'ı application/json yapar + Status'ü 201 yapar + Encode eder + Gönderir
```

**4. Gelişmiş Routing**
`net/http` 1.22 öncesinde metot ayırmak imkansızdı. Fiber'de `app.Post`, `app.Get`, `app.Put` gibi metotlar çok nettir; `/user/:name` gibi parametreler çok daha güçlüdür.

**5. "Zero-Allocation" ve Hız**
Fiber, Go'nun standart kütüphanesini kullanmak yerine `fasthttp` üzerine kuruludur. Çok yüksek trafikli sistemlerde daha az bellek harcar; çünkü request nesnelerini çöpe atmaz, yeniden kullanır.

### Neden `net/http` Önce Öğretiliyor?

1. **Evrensellik:** Her Go geliştiricisi bunu bilir. Fiber spesifiktir.
2. **Middleware:** Bir `http.Handler`'ı nasıl sarmalayacağını anlamak, Go'daki tüm loglama, kimlik doğrulama ve izlemenin temelidir.
3. **Debugging:** Fiber'in "sihri" başarısız olduğunda, sorunu düzeltmek için temel HTTP ilkelerini (header'lar, status kodları, byte'lar) anlamanız gerekir.

> **Ne zaman Fiber kullanmalısın?** Hızlıca bir mikroservis yazman gerekiyorsa, ekipçe standart bir yapı istiyorsanız ve "boilerplate" (tekrar eden) kod yazmaktan kaçınmak istiyorsan Fiber harikadır.

---

## 10. Özet: Akış Şeması

```
main()
  └── ListenAndServe(":9000", mux)
        └── [İstek Geldi] → Yeni Goroutine Aç
              └── mux.ServeHTTP(w, r)
                    ├── "POST /register" → registerHandler(w, r)
                    │     ├── json.Decode → User struct
                    │     ├── userRegistry[u.Name] = u
                    │     └── json.Encode → 201 Created
                    └── "GET /user/{name}" → getUserHandler(w, r)
                          ├── r.PathValue("name")
                          ├── userRegistry[name] → (user, ok)
                          └── json.Encode → 200 OK  |  http.Error → 404
```

---

*Notlar derste işlenen konuları kapsamaktadır. Fiber ile karşılaştırmalı örnekler, saf Go'daki "elle yapılan" işlemlerin framework'lerin arkasında nasıl çalıştığını anlamak için eklenmiştir.*