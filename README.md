# Learn Paxos

learn-paxos is a Go program for learning the basics of the Paxos algorithm. More
specifically, it is for learning an algorithm variously referred to as "Classic
Paxos" or "Single-Decree Paxos."

It my opinion that Paxos is much simpler to learn than is commonly believed.
Just  as the central concept of Bubble Sort can be described by the idea of
bubbling smaller items to the front of a list, the idea of Classic Paxos is
captured by the following conversation:

> Alice: Where would you two like to go to lunch?
> 
> Bob and Carlos: I don't know, where do you want to go?
> 
> Alice: Let's get some tacos!
> 
> Bob and Carlos: OK, let's get some tacos!

There really is not much more to the intuition behind Classic Paxos than an
idea of why the above conversation works. OK, sure, there are some identifiers,
a quorum concept, the fact that messages can be dropped or rearranged, and that
participants can become unreachable. But there are details to bubble sort too,
and such details enrich our understanding of the algorithm rather than stand in
the way of our intuition.

The code in this repository closely follows the algorithms in chapter 2 of Dr.
Heidi Howard's Ph.D. dissertation [1]. My recommended study plan for learning
Classic Paxos would be as follows:

1. Read sections 2.2 and 2.3 of Dr. Howard's dissertation. 
2. Implement algorithms 3 and 4, the proposer and acceptor algorithms, in
   a programming language of your choice. You're welcome to clone this
   repository if you'd like.
3. Read sections 2.4-2.7 of Dr. Howard's dissertation, to understand
   the safety and progress properties of Classic Paxos.

## Terminology

When I try to understand code that corresponds to an academic work, I am
frequently confused by differences in terminology between the code and the
paper. To make understanding the code in this repository easier, it adheres to
the following naming conventions:

Name in [1]     | Name in this code | Meaning
----------------| ----------------- | ---------------------------------
e               | epoch             | epoch (component of all messages)
v               | value             | proposal value
n<sub>a</sub>   | nAcceptors        | number of acceptors
Q<sub>P</sub>   | promisedAcceptors | acceptors that have promised
Q<sub>A</sub>   | acceptedAcceptors | acceptors that have accepted
f               | lastAcceptedEpoch | last accepted epoch (component of prepare message)
v               | lastAcceptedValue | last accepted value (component of prepare message)
𝛾               | candidateValue    | value that proposer will propose if it is sure no other value is chosen
e<sub>pro</sub> | promisedEpoch     | last promised epoch
e<sub>acc</sub> | acceptedEpoch     | last accepted epoch
v<sub>acc</sub> | acceptedValue     | last accepted value

## Nil values

In the presentation in [1], epochs and values may be nil. In the code in this
repository, a nil epoch is any epoch in which the *big.Int field is nil (any
epoch e for which e.nil() returns true), and the nil value is the empty string.

## Persistence

If participants in Paxos may both crash and recover, then the algorithm requires
that the participants remember some of their variables across restarts, i.e.,
they must store those variables on a durable storage medium, such as a disk. (If
participants may crash but not recover, the durability requirements disappear,
but I digress.) When first learning Paxos, the best way to understand
persistence requirements is to pay attention to what variables must be
persisted. I believe that any more detailed consideration of persistence is a
distraction when first learning Paxos, and that reason, I have omitted any
notion of persistence from this repository.

## References

[1] Heidi Howard. 2019. _Distributed Consensus Revised_. University of Cambridge
Computer Laboratory Technical Report UCAM-CL-TR-935. University of Cambridge,
Cambridge, England. Retrieved from https://www.cl.cam.ac.uk/techreports/.