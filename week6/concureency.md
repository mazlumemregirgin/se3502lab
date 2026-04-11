# Don't communicate by sharing memory; share memory by communicating.

# Concurrency in Go: From Goroutines to Channels

## How to write production-grade concurrent systems with WaitGroup, Mutex, and Channel

**Backend Engineering · Go / Golang · ~12 min read**

---

The moment you take your first steps as a backend developer, you inevitably face one question: how do you do many things at the same time? Think of an HTTP server — thousands of requests arriving every second, each one fetching data from a database, connecting to an external API, processing a file. Doing these sequentially simply isn't an option. This is exactly where concurrency comes in.

Go is one of the rare languages that places this problem at the very center of its design. And in doing so, it embraces a guiding philosophy: *don't communicate by sharing memory — share memory by communicating.* In this article, we'll internalize this philosophy through practical examples.

---

## Concurrency ≠ Parallelism

Before we dive in, let's clarify two concepts that are frequently confused. Mixing these up can lead to poor architectural decisions.

**Concurrency**
- The capacity to *manage* multiple tasks
- Can work even on a single CPU
- Tasks progress sequentially but interleave with each other
- Like a cashier handling 3 customers

**Parallelism**
- Actually *doing* multiple tasks at the same time
- Requires multiple CPUs
- Tasks run truly simultaneously
- Like 3 cashiers each serving a customer at once

Go supports both concurrency and parallelism. But what makes Go special is how *cheap* concurrency is.

---

## Goroutine: The Lightweight Sibling of Threads

In most programming languages, spawning a new thread carries a serious cost. A native OS thread in Java or C++ consumes **1–2 MB** of stack memory. This makes running thousands of threads simultaneously impractical.

In Go, goroutines start with only **~2 KB** of stack and grow dynamically as needed. This difference may seem small, but it makes something remarkable possible: you can run tens of thousands of goroutines simultaneously — not like opening OS threads, but like calling a function.

### How does it work?

Go's runtime includes its own scheduler. This scheduler operates on the **M:N model**: it distributes M goroutines across N OS threads. This means thousands of goroutines run on far fewer real OS threads. You just write `go` — the runtime handles the rest.

Starting a goroutine is this simple:

```go
package main

import (
	"fmt"
	"time"
)

func helloWorld() {
	fmt.Println("Hello, goroutine!")
}

func main() {
	go helloWorld() // goroutine started
	time.Sleep(500 * time.Millisecond)
	fmt.Println("main finished")
}
```

You simply add `go` before the function call. That's it. But there's a critical detail here: when the main goroutine exits, *all* other goroutines die with it. That's why we added `time.Sleep()` — but this is not the real solution. In production code, milliseconds matter; hardcoding arbitrary wait times is both incorrect and dangerous.

> **Warning:** Trying to wait on goroutines with `time.Sleep()` can be catastrophic in latency-critical systems like fintech. The right tool is `sync.WaitGroup`.

---

## WaitGroup: "Wait for All, Then Continue"

`sync.WaitGroup` is a synchronization primitive designed specifically to wait for goroutines to finish. Under the hood, it maintains a counter: how many goroutines are we waiting on? When the counter reaches zero, the program continues.

There are three methods, each with a clear role:

| Method | Responsibility |
|---|---|
| `wg.Add(n)` | "I'm launching n goroutines, increment the counter" |
| `wg.Done()` | "This goroutine is done, decrement the counter" (use with `defer`) |
| `wg.Wait()` | "Block here until the counter reaches zero" |

```go
package main

import (
	"fmt"
	"sync"
)

func greet(wg *sync.WaitGroup) {
	defer wg.Done() // decrement counter when function returns — guaranteed
	fmt.Println("Hello, WaitGroup!")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3) // we'll launch 3 goroutines

	go greet(&wg)
	go greet(&wg)
	go greet(&wg)

	wg.Wait() // block until all are done
	fmt.Println("All goroutines completed.")
}
```

> **Best Practice:** Always define `wg.Done()` with `defer` at the top of the function. Think of `defer` like a `finally` block in Java — it runs no matter how the function exits. It guarantees the counter will be decremented even in a panic scenario.

One important note: you must pass the WaitGroup to goroutines **as a pointer** (`&wg`). If you pass it by value, each goroutine works with its own copy and nothing gets synchronized.

---

## Mutex: The Guardian of Data Integrity

Goroutines are great — but there's a lurking problem. What happens when multiple goroutines read and write the same data at the same time?

Consider a classic banking scenario: an account holds 100₺. User A and User B both want to deposit 10₺ at the same time. Both read the balance as 100₺, add 10, and write back 110₺. The result: 110₺ instead of 120₺. 10₺ has vanished.

This problem is called a **Race Condition**. The solution is `sync.Mutex`.

```go
package main

import (
	"fmt"
	"sync"
)

type BankAccount struct {
	mu      sync.Mutex // protective lock
	balance int
}

func (b *BankAccount) Deposit(amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	b.mu.Lock()         // lock the door — no one else can enter
	b.balance += amount // critical section: touching shared data
	b.mu.Unlock()       // unlock — next goroutine can proceed
}

func main() {
	var wg sync.WaitGroup
	account := BankAccount{balance: 0}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go account.Deposit(1, &wg)
	}

	wg.Wait()
	fmt.Println("Final Balance:", account.balance) // Always 1000
}
```

Here's how a Mutex works: when `Lock()` is called and another goroutine currently holds the lock, the new goroutine blocks and waits. Once the lock is released, it proceeds. This ensures that only one goroutine enters the critical section at a time.

> **Trade-off:** When you use a Mutex, you're temporarily making your asynchronous code sequential. But this is an acceptable trade-off: showing incorrect data is worse than showing no data at all. Only lock the minimum area that touches shared data — overly broad locking slows the entire system down.

---

## Channel: The Shared Line Between Goroutines

A Mutex protects shared memory. But what if goroutines need to *send data to each other*? This is where one of Go's most elegant ideas emerges: **Channel**.

Think of a channel like a conveyor belt. One goroutine places a product on the belt; another picks it up. They don't need to know each other directly — the shared line is enough. This is precisely what the manifesto at the beginning of this article means: instead of directly sharing memory, talk through a communication channel.

### Unbuffered Channel: Synchronous Handshake

The default channel is unbuffered. The sender waits until the receiver is ready. The receiver waits until the sender is ready. They wait for each other — this is a synchronization mechanism.

```go
package main

import (
	"fmt"
	"time"
)

func produce(ch chan string) {
	for i := 1; i <= 3; i++ {
		item := fmt.Sprintf("RawMaterial-%d", i)
		fmt.Println("Producer:", item, "ready")
		ch <- item // send to channel — waits until receiver is ready
		time.Sleep(time.Second)
	}
	close(ch) // signal: "no more data"
}

func pack(ch chan string, done chan bool) {
	for item := range ch { // read until channel is closed
		fmt.Println("Assembler:", item, "packaging...")
		time.Sleep(time.Second * 2)
		fmt.Println("Assembler:", item, "PACKAGED ✓")
	}
	done <- true // tell the main program "all done"
}

func main() {
	conveyor := make(chan string)
	done := make(chan bool)

	go produce(conveyor)
	go pack(conveyor, done)

	<-done // wait until "all done" arrives
	fmt.Println("Factory closed, all jobs finished.")
}
```

### Buffered Channel: The Inventory Buffer

With an unbuffered channel, the producer and assembler are bound to each other's pace. If the producer is much faster, it constantly waits — a bottleneck.

The solution: place an inventory buffer in between. We do this with a **buffered channel**:

```go
// Unbuffered — sender waits for receiver
conveyor := make(chan string)

// Buffered — capacity for 10 items
// Producer can keep going without waiting until the buffer is full
conveyor := make(chan string, 10)
```

The `10` here represents the capacity: the producer can push up to 10 items into the channel and continue without waiting, even if the receiver hasn't picked them up yet. When the channel is full, the producer waits again. This asynchronous communication increases the system's overall throughput.

### Channel Directions

Go lets you add directional constraints to channels, making your API safer and more readable:

- `chan string` — bidirectional (either side can use it)
- `chan← string` — send-only
- `←chan string` — receive-only

Finally: **don't forget to close the channel.** Calling `close(ch)` signals "no more data will come from this channel." A goroutine listening with `for range ch` exits the loop when the channel closes. If you don't close it, the goroutine waits forever — this is a **goroutine leak**, and in production it causes memory bloat.

---

## Putting It All Together

Each tool we've covered solves a different problem:

| Tool | Purpose |
|---|---|
| **Goroutine** | Launch tasks concurrently. Lightweight, cheap, thousands at once. |
| **WaitGroup** | Wait for goroutines to finish. Dynamic control instead of arbitrary `Sleep()`. |
| **Mutex** | Safe access to shared data. Prevents race conditions. |
| **Channel** | Communication between goroutines. Send messages instead of sharing memory. |

**When to use which?**

- Goroutines are modifying shared data → **Mutex**
- Goroutines need to pass data to each other → **Channel**
- You need to wait for goroutines to finish → **WaitGroup**

---

Go's concurrency model is woven into the language's DNA. The `go` keyword alone is powerful enough, but combined with the `sync` package and channels, it forms an ecosystem where you can build real production systems.

From fintech to game servers, from real-time data pipelines to HTTP services — you'll reach for these tools constantly. What matters is knowing *when* to choose which one — and keeping the manifesto in mind:

> *"Don't communicate by sharing memory; share memory by communicating."*

Every tool we covered in this article reflects a different dimension of this philosophy. WaitGroup is coordination, Mutex is safety, Channel is communication. When you use all three together, you understand why Go is so widely adopted.

---

*Tags: goroutine · waitgroup · mutex · channel · concurrency · golang · backend*