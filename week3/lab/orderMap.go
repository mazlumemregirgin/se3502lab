package lab

import "fmt"

func OrderMap() {

	// Go'da map'ler sırasızdır. Yani, map'e eklenen öğelerin sırası garanti edilmez.
	// Bu nedenle, map'leri sıralı bir şekilde tutmak için ek bir veri yapısı kullanmamız gerekir.

	colors := map[string]string{}
	colors["AliceBlue"] = "#F0F8FF"
	colors["Coral"] = "#FF7F50"
	colors["DarkGray"] = "#A9A9A9"
	colors["ForestGreen"] = "#228B22"
	colors["Indigo"] = "#4B0082"
	colors["Lime"] = "#00FF00"
	colors["Navy"] = "#000080"
	colors["Orchid"] = "#DA70D6"
	colors["Salmon"] = "#FA8072"

	for key, value := range colors {
		fmt.Printf("%s:%s, ", key, value)
	}
}
