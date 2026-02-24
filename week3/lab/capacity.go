package lab

import "fmt"

func Capacity() {

	var slowSlice []int
	for i := 0; i < 1000000; i++ {
		slowSlice = append(slowSlice, i)
	}

	fmt.Println("slow", slowSlice)
	fmt.Println("Length:", len(slowSlice))
	fmt.Println("Capacity:", cap(slowSlice))

	fastSkice := make([]int, 0, 1000000)
	for i := 0; i < 1000000; i++ {
		fastSkice = append(fastSkice, i)
	}

	fmt.Println("fast", fastSkice)
	fmt.Println("Length:", len(fastSkice))
	fmt.Println("Capacity:", cap(fastSkice))

}
