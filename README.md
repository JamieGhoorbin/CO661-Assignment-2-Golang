# CO661-Assignment 2
# Concurrent programming in Go

## Solving three problems using synchronous and asynchronous message passing in Go.

## Final grade: 95%

### **The dentist problem (part 1) - 50%**

Implementing a dentist studio system, run by one single dentist. The dentist is sitting in a treatment room, waiting to meet patients. Out of the treatment room there is a small waiting room for n patients to sit in a FIFO queue.

**The dentist.** The dentist checks for patients in the waiting room.

- If there are no patients, the dentist falls asleep.
- If there are is at least one patient, the dentist calls the first one in. The remaining patients keep waiting. During the treatment, the dentist is active while the patient is sleeping . When the dentist finishes the treatment, the patient is woken up, and the dentist checks for patients in the waiting room. And so on …

**The patient.** The patient, upon arrival, checks if the dentist is busy with other patients or sleeping.

- If the dentist is sleeping, the patient wakes the dentist up and falls asleep while being treated. The patient is woken up when the treatment is completed, and leaves (i.e., terminates).
- If the dentist is busy with another patient, the arriving patient goes in the waiting room and waits (i.e., sleeps). When the patient is woken up, the treatment starts: the patient falls asleep until being woken up at the end of the treatment.

> **My task.** Model the dentist problem. You will need to create a dentist function, and a patient function. Use (synchronous or asynchronous) channels to synchronize the activities of dentist and patient.

Constraints/hints.

1. Use message passing (no shared memory, counters, arrays, …).
2. Use asynchronous channel wait of size n for the waiting queue of patients.
3. If patients arrive and wait is full (e.g., m>n) they will block (see Go’s implementation of write operation on full channels) until a space is free. This is ok. For the moment, just let this be.
4. Use synchronous channel dent to model the state (sleeping/awake) of the dentist. The dentist should fall asleep while reading on dent, not while reading on wait.
5. Use synchronous channels, one for each patient, to model the state (sleeping/awake) of that patient.
   a. Function patient creates this as a fresh channel (type chan int).
   b. This fresh channel will be added to wait (type chan chan int) if the patient needs to stay in the waiting room.
   c. The patient sleeps while waiting to read on this fresh channel: (1) when waiting in the queue, and (2) when having the treatment done.
6. Allow the treatment to take some random time (using time.Sleep()). This needs to be at the dentist side (as the patient is sleeping while having a treatment done).

### **Introducing priorities (part 2) – 25%**

> (2.a) Some patients need emergency procedure and need to be prioritised. The dentist establishes a priority-based queue system that gives some patients priority over others. Modify your solution to part 1 so that instead of channel wait there are two channels hwait and lwait for high- and low-priority patients, respectively.
> The dentist will check that there are no high-priority patients before serving low-priority patients.
> The dentist needs to set a timer and move a patient from lwait to hwait whenever m milliseconds have passed since the last read from lwait.

> (2.b) Assume Go does not have a fair semantics. Can you identify one possibility of starvation in the scenario described in the problem statement of part 1? Your answer should be justified (you can refer to your own code if you wish). If you identify one, say what you could change your solution to make the system starvation-free. Your fix needs to be general with respect to the set-up scenario created by the main function (i.e., it cannot require assumptions on the relationship between the size of queues or number of patients initialized by main).

### **The assistant (part 3) – 25%**

> (3.a) The dentist is overwhelmed by the overhead of handling two queues (lwait and hwait) and decides to hire an assistant. Modify your solution to part 2 by introducing a new goroutine “assistant” that communicates with the patients using the queues hwait and lwait, and communicates with the dentist using one single queue wait. The dentist will not see or act on the queues hwait and hwait but only receive patients on wait. You need to modify the function dentist and main accordingly.

> (3.b) Can you identify a possibility of deadlock that affects part 2 but not part 3?
