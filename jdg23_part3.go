// (3b) Yes. Part 2 in the dentist function, if a lwait patient is moved (upgraded)
// to hwait and hwait is full, then it blocks and we reach a deadlock. I was able
// to observe this behaviour when decreasing the hwait size and placing print
// statements to see where the dentist function was getting blocked. This will not
// happen in part 3 providing the wait channel size is large enough. We use a fan
// in approach for patients into the wait queue, with priority given to hwait patients.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Create a synchronous channel
	dent := make(chan chan int)
	// Create a high priority wait queue
	hwait := make(chan chan int, 100)
	// Create a low priority wait queue
	lwait := make(chan chan int, 5)
	// Create a wait queue
	wait := make(chan chan int, 150)

	// Spawn an assistant
	go assistant(hwait, lwait, wait)

	// Spawn a dentist
	go dentist(wait, dent)

	// Sleep to delay creation of patients
	time.Sleep(2 * time.Second) // jdg23

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
	time.Sleep(30 * time.Second)
}

// The Dentist
func dentist(wait chan chan int, dent <-chan chan int) {
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

// The Assistant
func assistant(hwait chan chan int, lwait <-chan chan int, wait chan<- chan int) {
	// Timer for moving patients from lwait to hwait
	timer := time.NewTimer(500 * time.Millisecond)
	for {
		// Aging - move patient from lwait to wait whenever time period has passed
		select {
		case <-timer.C:
			select {
			case lpatient := <-lwait:
				wait <- lpatient // Issue: will block here if wait is full
			default:
				// Default case here prevents blocking if lwait empty
			}
			// Reset the timer after checking lwait
			timer = time.NewTimer(500 * time.Millisecond)
		default:
			// Default case here prevents blocking if lwait empty
		}

		// Prioritising hwait over lwait patients
		select {
		case hpatient := <-hwait:
			wait <- hpatient
		default:
			select {
			case lpatient := <-lwait:
				wait <- lpatient
			default:
				// Default case here prevents blocking if lwait empty
			}
			// Reset the timer after checking lwait
			timer = time.NewTimer(50 * time.Millisecond)
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
