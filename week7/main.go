package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// User modeli
// önce databaesimizde user tablosu olduğunu varsayalım ve bu tabloda name ve email alanları olsun. bu alanları struct olarak tanımlayarak json tagleri ekliyoruz.

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// normalde bu tür bir uygulamada veritabanı kullanırdık ancak bu örnekte basit olması açısından bir map kullanarak kullanıcıları saklayacağız. bu map'in key'i kullanıcı adı olacak ve value'su ise User struct'ı olacak.
var userRegistry = make(map[string]User)

func main() {
	// 1. Multiplexer (Router) tanımlıyoruz. bu aslında fiber kodlarken routes.go dosyasında yaptığımız işlemin manuel hali gibi düşünebiliriz. fiber'de c.Handle() ile yaptığımız işlemi burada http.NewServeMux() ile yapacağız.
	// fiberda app:= fiber.New() ile oluşturduğumuz app aslında bir multiplexer'dır. fiber bize bu işlemi kolaylaştırır ancak burada manuel olarak yapacağız.
	mux := http.NewServeMux()

	// 2. Route tanımları
	// Method Matching: Sadece POST veya GET isteklerine cevap verir.
	// bu aslında app.Get() veya app.Post() ile yaptığımız işlemin manuel hali gibi düşünebiliriz. fiber'de c.HandleFunc() ile yaptığımız işlemi burada mux.HandleFunc() ile yapacağız.
	mux.HandleFunc("POST /register", registerHandler)
	mux.HandleFunc("GET /user/{name}", getUserHandler)

	// 3. Sunucuyu başlatma
	fmt.Println("Server 9000 portunda çalışıyor...")
	// ListenAndServe bir blocking fonksiyondur, server açık olduğu sürece alt satıra geçmez.
	// bu fonksiyon parametre olarak port numarası ve multiplexer'ı alır. eğer sunucu başlatılırken bir hata oluşursa bu hatayı ekrana yazdırırız.

	if err := http.ListenAndServe(":9000", mux); err != nil {
		fmt.Printf("Hata: %s\n", err)
	}
}

// registerHandler: Kullanıcıyı kaydeder

// bu fonksiyon parametere olarka http.ResponseWriter ve *http.Request alır. http.ResponseWriter, sunucunun istemciye vereceği cevabı yazmak için kullanılır. *http.Request ise istemciden gelen isteği temsil eder.
// bir http isteği çok büyük olabilir bu yüzden eğer pointer yerine struct olarak alırsak bu büyük bir kopyalama işlemi olur ve performansı düşürür. bu yüzden pointer olarak alırız.
func registerHandler(w http.ResponseWriter, r *http.Request) {

	// önce bu user için bir değişken tanımlayalım. var u User ile tanımlayarak bu değişkeni User struct'ına sahip bir değişken olarak tanımlamış olduk. bu değişkeni daha sonra json'dan gelen veriyi decode etmek için kullanacağız.
	var u User

	// Fiber'deki c.BodyParser'ın manuel hali:
	// sunucuya gelen istek 1 ve 0lardan oluşan anlamsız bir veridir. bu veriyi bizim anlamlı bir şekilde kullanabilmemiz için decode etmemiz gerekir. bu işlemi json.NewDecoder(r.Body).Decode(&u) ile yapacağız. biz gelen verinin json formatında olduğunu varsayıyoruz. bu yüzden r.Body'den gelen 1 ve 0 ları json.NewDecoder(r.Body).Decode(&u) ile User struct'ına decode eder ve u değişkenine atar.
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "JSON çözülemedi", http.StatusBadRequest)
		return
	}

	if u.Name == "" {
		http.Error(w, "İsim boş olamaz", http.StatusUnprocessableEntity)
		return
	}
	// kullanıcıyı kaydetme işlemi burada gerçekleşir. userRegistry map'ine kullanıcıyı ekleriz. key olarak kullanıcı adı ve value olarak ise User struct'ını veririz.
	userRegistry[u.Name] = u

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)

	/*
		normalde bu user kayıt fonksiyonunu fiber ile yazsaydık:

		// fiberda ayrı ayrı respose, request parametreleri vermek yerine tek bir parametre olan c *fiber.Ctx alırdık.
		// c.BodyParser dediğimde requesti okuyor, c.JSON dediğimde ise responseu yazıyor.
		func registerHandler(c *fiber.Ctx) error {
			var u User
			manuel decoder yerine fiberin c.BodyParser() fonksiyonu kullanılır.
			if err := c.BodyParser(&u); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON çözülemedi"})
			}
			if u.Name == "" {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "İsim boş olamaz"})
			}
			userRegistry[u.Name] = u

			// header + status + encode tek satırda yapılırdı.
			return c.Status(fiber.StatusCreated).JSON(u)
		}
	*/
}

// getUserHandler: Path parametresinden ismi alır ve kullanıcıyı döner
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Path param alımı (Fiber'de c.Params("name"))
	// burada requestten path parametresini almak için r.PathValue("name") kullanırız. bu fonksiyon path'ten name parametresini alır ve bize verir. bu name parametresi aslında kullanıcı adını temsil eder ve biz bu kullanıcı adını kullanarak userRegistry map'inden kullanıcıyı bulacağız.
	name := r.PathValue("name")

	// burada userRegistry map'inden kullanıcıyı bulmak için user, ok := userRegistry[name] kullanırız. bu kod parçası aslında userRegistry map'inde name anahtarına sahip bir kullanıcı olup olmadığını kontrol eder. eğer böyle bir kullanıcı varsa user değişkenine atanır ve ok değişkeni true olur. eğer böyle bir kullanıcı yoksa user değişkeni boş olur ve ok değişkeni false olur. bu şekilde kullanıcıyı bulup bulmadığımızı kontrol ederiz.
	// eğer kullanıcı yoksa go hata vermez, sadece bos bir user döner ve ok false olur. bu yüzden if !ok ile kullanıcıyı bulamadığımızı kontrol ederiz ve eğer kullanıcı yoksa http.Error() ile istemciye bir hata mesajı göndeririz.
	user, ok := userRegistry[name]
	if !ok {
		http.Error(w, "Kullanıcı yok", http.StatusNotFound)
		return
	}
	// kullanıcıyı bulduktan sonra ise bu kullanıcıyı json formatında istemciye göndermek için w.Header().Set("Content-Type", "application/json") ile header'ı set ederiz ve json.NewEncoder(w).Encode(user) ile user değişkenini json formatında encode ederek istemciye göndeririz.
	w.Header().Set("Content-Type", "application/json")
	// burada ise responseda go strcutını (user) ali json metnine dönüştür ve bunu response olarak yazdırmak için json.NewEncoder(w).Encode(user) kullanırız. bu kod parçası aslında user struct'ını json formatına encode eder ve w'ye yazar. w ise http.ResponseWriter olduğu için bu json verisi istemciye gönderilir.
	json.NewEncoder(w).Encode(user)
}
