// (2b) On the assumption that Go does not have fair semantics, a possibility
// for starvation in part 1 could occur when the dentist is retrieving a
// patient (goroutine) from the wait channel. Some goroutines waiting in the
// channel may never get the chance to make progress if more patients are
// being added to the wait queue.

// I can’t think of a solution to this problem using channels. We could use
// a shared memory approach using the sync.Mutex package to lock during
// critical sections. Then dentist and patients could be some struct and we
// could maintain a list (slice) of soem custom type struct and append
// patients to the list. The first element could be dequeued using someQueue[1:].

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Synchronous channel to model (sleeping/awake) state of dentist
	dent := make(chan chan int)
	// High priority waiting queue
	hwait := make(chan chan int, 100)
	// Low priority waiting queue
	lwait := make(chan chan int, 5)
	// Spawn a dentist
	go dentist(hwait, lwait, dent)

	// Sleep to delay creation of patients
	time.Sleep(3 * time.Second) // jdg23

	high := 10
	low := 3

	// Spawn low priorty and high priorty patients
	for i := low; i < high; i++ {
		go patient(hwait, dent, i)
	}

	for i := 0; i < low; i++ {
		go patient(lwait, dent, i)
	}

	// Sleep to prevent main thread ending
	time.Sleep(40 * time.Second)
}

// The Dentist
func dentist(hwait chan chan int, lwait <-chan chan int, dent <-chan chan int) {
	// Timer for moving patients from lwait to hwait
	timer := time.NewTimer(1000 * time.Millisecond)
	for {
		// Aging - move patient from lwait to hwait whenever time period has passed
		select {
		case <-timer.C:
			select {
			case lpatient := <-lwait:
				hwait <- lpatient // Issue: will block if hwait if full
			default:
				// Default case here prevents blocking if lwait empty
			}
			// Reset the timer after checking lwait
			timer = time.NewTimer(1000 * time.Millisecond)
		default:
			// Default case here prevents blocking if lwait empty
		}

		// Prioritising hwait over lwait patients
		select {
		case hpatient := <-hwait:
			// Sleep for duration of treatment
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			hpatient <- 2
			hpatient <- 3
		default:
			select {
			case lpatient := <-lwait:
				lpatient <- 1
				// Sleep for duration of treatment
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
				lpatient <- 2
				lpatient <- 3
				// Reset the timer after checking lwait
				timer = time.NewTimer(1000 * time.Millisecond)
			default:
				fmt.Println("Dentist is sleeping...")
				nextPatient := <-dent
				nextPatient <- 1
				// Sleep for duration of treatment
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
				nextPatient <- 2
				nextPatient <- 3
			}
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
