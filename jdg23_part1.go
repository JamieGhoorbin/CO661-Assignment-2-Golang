package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Size of the wait queue
	n := 5
	// No patients to create
	m := 5
	// Synchronous channel to model (sleeping/awake) state of dentist
	dent := make(chan chan int)
	// Waiting queue
	wait := make(chan chan int, n)

	// Spawns one dentist and m patients with id i
	go dentist(wait, dent)
	time.Sleep(3 * time.Second)
	for i := 0; i < m; i++ {
		go patient(wait, dent, i)
		time.Sleep(1 * time.Second)
	}
	// Sleep to prevent main thread ending
	time.Sleep(8000 * time.Millisecond)
}

// The Dentist
func dentist(wait <-chan chan int, dent <-chan chan int) {
	for {
		select {
		case nextPatient := <-wait: // Dentist checks waiting room
			nextPatient <- 1
			// Sleep for duration of treatment
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			nextPatient <- 2
			nextPatient <- 3
		default: // Dentist sleeps until patient arrives
			fmt.Println("Dentist is sleeping...")
			nextPatient2 := <-dent
			nextPatient2 <- 1
			// Sleep for duration of treatment
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			nextPatient2 <- 2
			nextPatient2 <- 3
		}
	}
}

// The Patient
func patient(wait chan<- chan int, dent chan<- chan int, id int) {
	fmt.Printf("Patient %d wants a treatment\n", id)
	// Synchronous channel for patient
	pChan := make(chan int)
	// defer close(pChan)
	// Check if the dentist is available.
	select {
	case dent <- pChan: // Wake up dentist and start treatment
	default: // Patient enters waiting room
		fmt.Printf("Patient %d is waiting\n", id)
		wait <- pChan
	}
	<-pChan // 1 - Receive when patient starts treatment
	fmt.Printf("Patient %d is having a treatment\n", id)
	<-pChan // 2 - Receive when patient is woken up
	fmt.Printf("Patient %d has shiny teeth!\n", id)
	<-pChan // 3 - Done
}
