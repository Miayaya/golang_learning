package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Increment the value pointed to by the pointer by 10.
func AddNum(num *int) {

	fmt.Printf("After add, the num is %d \n", *num+10)
}

// Double each element in the slice.
func SliceMul(nums []int) []int {
	for i, v := range nums {
		nums[i] = v * 2
	}
	return nums
}

// Odd numbers
func Odd(num int) {
	for i := 0; i <= num; i++ {
		if i%2 == 1 {
			fmt.Printf("print the odd number %d \n", i)
		}
	}
}

// Even numbers
func Even(num int) {
	for i := 0; i <= num; i++ {
		if i%2 == 0 {
			fmt.Printf("print the even number %d \n", i)
		}
	}
}

//Task Schedule

func Schedule() {
	var wg sync.WaitGroup
	tasks := []func(){
		// func() { time.Sleep(1 * time.Second); fmt.Println("Task1 done") },
		// func() { time.Sleep(2 * time.Second); fmt.Println("Task2 done") },
		// func() { time.Sleep(4 * time.Second); fmt.Println("Task3 done") },
		func() { fmt.Println("Task1 done") },
		func() { fmt.Println("Task2 done") },
		func() { fmt.Println("Task3 done") },
	}
	for i, task := range tasks {
		wg.Add(1)
		go func(i int, task func()) {
			defer wg.Done()
			start := time.Now()
			task()
			fmt.Printf("Task %d use time %v \n", i+1, time.Since(start))

		}(i, task)
	}
	wg.Wait()
	fmt.Println("All tasks done")

}

// Define Shape interfacce
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Recangle
type Rectangle struct {
	width, height float64
}

func (r *Rectangle) Area() float64 {
	return r.width * r.height
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

// Circle
type Circle struct {
	Redius float64
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Redius * c.Redius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Redius
}

type Person struct {
	Name string
	Age  int
}
type Employee struct {
	Person
	EmployeeID string
}

func (e *Employee) ShowInfo() {
	fmt.Printf("My name is %s, age is %d and employeeid is %s.\n", e.Name, e.Age, e.EmployeeID)
}

// Channel communication
func ChannelCom() {
	var wg sync.WaitGroup
	ch := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range ch {
			fmt.Printf("Received: %d \n", msg)
		}
	}()

	go func() {
		defer close(ch)
		for i := 1; i <= 10; i++ {
			ch <- i
			fmt.Printf("Generated number %d \n", i)
		}
	}()

	wg.Wait()
	fmt.Println("All done!")
}

//Channel buffer

func ChannelBuffer() {
	var wg sync.WaitGroup
	ch := make(chan int, 50)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range ch {
			fmt.Printf("Received: %d \n", msg)
		}
	}()

	go func() {
		defer close(ch)
		for i := 1; i <= 100; i++ {
			ch <- i
			fmt.Printf("Generated number %d \n", i)
		}
	}()

	wg.Wait()
	fmt.Println("All done!")

}

type Counter struct {
	mu    sync.Mutex
	count int64
}

func NewCounter() *Counter {
	return &Counter{}
}

func (sc *Counter) SafeIncrement(m int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	for i := 0; i < 1000; i++ {
		sc.count++
	}
	fmt.Printf("goroutine %d incremented counter %d\n", m, sc.count)
}

func (sc *Counter) Increment() {
	atomic.AddInt64(&sc.count, 1)
	// fmt.Printf("goroutine %d incremented counter %d\n", m, sc.count)
}
func (sc *Counter) Value() int64 {
	return atomic.LoadInt64(&sc.count)
	// fmt.Printf("goroutine %d incremented counter %d\n", m, sc.count)
}

func (sc *Counter) GetCount() int64 {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.count
}

// Lock/unlock counter
func SafeCounter() {

	var wg sync.WaitGroup
	counter := NewCounter()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			counter.SafeIncrement(i)
		}(i)
	}

	wg.Wait()
	fmt.Printf("The total count is %d \n", counter.count)

}

// sync/atomic counter
func AtomiCounter() {
	var wg sync.WaitGroup
	counter := NewCounter()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				counter.Increment()
			}
			fmt.Printf("goroutine %d incremented counter %d\n", i, counter.count)
			// counter.Increment(i)
		}(i)
	}

	wg.Wait()
	fmt.Printf("The total count is %d \n", counter.Value())

}

func main() {
	//Pointer
	// testNum := 5
	// fmt.Printf("The testNum is %d \n", testNum)
	// AddNum(&testNum)
	//Slice
	// testSlice := []int{2, 5, 12, 3, 7, 5}
	// fmt.Printf("The testSlice is %v \n", testSlice)
	// testSlice = SliceMul(testSlice)
	// fmt.Printf("After sliceMult,the doubleSlice is %v \n", testSlice)

	// Goroutine
	// var wg sync.WaitGroup
	// wg.Add(2)
	// go func() {
	// 	defer wg.Done()
	// 	Odd(10)
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	Even(10)
	// }()
	// wg.Wait()

	//Task scheduler
	// Schedule()

	//Object-oriented approach
	// r := &Rectangle{width: 2.2, height: 3.5}
	// c := &Circle{Redius: 3}
	// var rec Shape = r
	// var cir Shape = c

	// fmt.Printf("Rectangle area is %v, perimeter is %v \n", rec.Area(), rec.Perimeter())
	// fmt.Printf("Circle area is %v, perimeter is %v \n", cir.Area(), cir.Perimeter())

	// em := &Employee{
	// 	Person:     Person{Name: "Lisa", Age: 22},
	// 	EmployeeID: "345678",
	// }
	// em.ShowInfo()

	//Channel
	// ChannelCom()
	// ChannelBuffer()
	// SafeCounter()
	AtomiCounter()

}
