package main

import (
	"se3502lab/week1/lab"
)

func main() {

	println("Hello, World!")

	/*
		WHY GO?

		Go was born at Google in 2007 out of frustration with existing languages
		for programming large software systems. It was designed to be simple,
		efficient, and easy to use, with a focus on concurrency and scalability.


		Go is compiled.
		Source code is translated directly into machine code for the spesific processor arcgitecture before it runs

		Go is statically typed.
		variable types are known at compile time. this cathces many errors before the code ever runs

		Go'da garbage collection var. bellek yönetimini manuel olarak yapamanı beklemez; bunun yerine
		arka planda çalışan ve artık kullanılmayan bellek alanlarını otomatik oalrak temizleyen bir sistem kullanılır.


		Go'da sadece 25 keyword vardır. Ne kadar az keyword var ise o kadar kolay öğrenilir ve okunur.


	*/

	/*
		WORKSPACE SETUP

		Eskiden GOPATH adında bir ortam değişkeni vardı ve tüm Go projeleri bu dizin altında bulunurdu.

		zorunlu klasör yapısı ise şöyle idi:
		$GOPATH/
			src/    # source files
			pkg/    # compiled package objects
			bin/    # compiled executable binaries

		Ancak, Go 1.11 ile birlikte modüller tanıtıldı ve bu zorunlu klasör yapısı kaldırıldı.
		Go modülleri, proje bağımlılıklarını yönetmek için kullanılan bir sistemdir. Proje kök dizininde go.mod dosyası oluşturulur ve bu dosya proje bağımlılıklarını tanımlar.
		Go modülleri sayesinde, projeler herhangi bir dizinde oluşturulabilir ve bağımlılıklar daha kolay yönetilebilir hale gelmiştir.


		go.mod dosyası nedir?
		go.mod dosyası, Go modüllerini tanımlayan bir dosyadır. Proje kök dizininde bulunur ve proje bağımlılıklarını yönetmek için kullanılır. Bu dosya, projenin hangi modüllere ihtiyaç duyduğunu ve bu modüllerin hangi sürümlerini kullanacağını belirtir.

		içeriği genellikle şu şekildedir:
		module myproject
		go 1.18
		require (
			github.com/some/dependency v1.2.3
		)


	*/

	/*
		VARIABLES, CONSTANTS, AND ZERO VALUES

		Go is explicit. you must declare variables before using them
		Bunun için iki yöntem vardır:
		1. var keyword -> used when you want to be spesific about the type or declare a variable without initializing it immediately.
			var age int = 30
			var name string = "Alice"

		2. short variable declaration (:=) -> inside functions, you can omit var and the type. go infers the type based on the value.
			func main(){
				age := 30 -> inferred as int
				name := "Alice" -> inferred as string
			}



		Go'da değişkenler sıfır değerine sahiptir. Bu, bir değişkenin değeri atanmadan önce otomatik olarak belirli bir değere sahip olduğu anlamına gelir. Sıfır değerler, türlere göre farklılık gösterir. Örneğin:
		int -> 0
		string -> "" (boş string)
		bool -> false
		float -> 0.0


		Constants are declared with const. They are created at compile time and cannot be changed.
		const Pi = 3.14
		const StatusOK = 200

	*/

	/*
		CONTROL FLOWS

		1. for döngüsü: Go'da while diye bir keyword yoktur. for döngüsü hem for gibi hem de while gibi davranabilir.
		for kelimesinin 3 farklı yapısı vardır:
			a. standard for loop:
			for i := 0; i < 5; i++ {
				println(i)
			}
			b. while'mış gibi gözüken for:
			count:= 0
			for count < 5 { -> sadece koşul yazarım, koşul doğru olduğu sürece döner. aynı while gibi
				println(count)
				count++
			}

			c. infinite loop:
			for { -> koşul yok, sonsuz döngü
				println("Hello")
				break -> sonsuz döngüyü kırmak için break kullanılır.
			}


		2. If/Else statements. -> standard if/else yapısı vardır. condiiton için paranteze gerek yoktur. ama braceslar zorunludur.
		if age >= 18 {
			println("You are an adult.")
		} else {
			println("You are a minor.")
		}

		3. Switch statements -> multiple conditions için kullanılır. Go'da switch, diğer dillerdeki gibi break gerektirmez. her case kendi içinde break yapar.
		switch day {
		case "Monday":
			println("Start of the week.")
		case "Friday":
			println("End of the week.")
		default:
			println("Midweek.")
		}
	*/

	/*
					LAB PRACTICES
				Exercise 1: Fahrenheit to Celsius Converter
				//se3502lab/week1/lab/c2f.go
				Create a CLI tool that prompts the user for a temperature in Fahrenheit, converts it to Celsius using the formula $C = (F - 32) \times \frac{5}{9}$, and prints the result

				Exercise 2: The Classic "FizzBuzz"
				//se3502lab/week1/lab/fizzbuzz.go
		    	Write a program that prints the numbers from 1 to 50. However:

		    	- For multiples of **3**, print "Fizz" instead of the number.
		    	- For multiples of **5**, print "Buzz" instead of the number.
		    	- For numbers which are multiples of **both 3 and 5**, print "FizzBuzz".

				Exercise 3: "Guess the Number" Game
				//se3502lab/week1/lab/guess.go
		    	Create a game where the computer selects a random number between 1 and 100. The user must guess the number.

		    	- If the guess is too high, print "Too high!"
		    	- If the guess is too low, print "Too low!"
		    	- If the guess is correct, print "You got it!" and exit the program.

	*/

	lab.FizzBuzz()
	lab.C2F()
	lab.Guess()
}
