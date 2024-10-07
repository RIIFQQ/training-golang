package main

import (
	"fmt"
	"sync"
	"time"
)

type BankAccount struct {
	RWmutex sync.RWMutex
	Balance int
}

func (account *BankAccount) AddBalance(amount int) {
	account.RWmutex.Lock()
	account.Balance += amount
	//account.Balance += amount artinya sama dengan account.Balance = account.Balance + amount
	account.RWmutex.Unlock()
}

func (account *BankAccount) GetBalance() int {
	account.RWmutex.RLock()
	balance := account.Balance
	account.RWmutex.RUnlock()
	return balance
}

func main() {
	account := BankAccount{}

	for i := 1; i <= 1000; i++ {
		go func() {
			for j := 1; j <= 100; j++ {
				account.AddBalance(1)
				fmt.Println(account.GetBalance())
			}
		}()
	}

	time.Sleep(10 * time.Second)
	fmt.Println("Final Balance:", account.GetBalance())
}
