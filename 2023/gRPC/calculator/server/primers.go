package main

import (
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
)

func (*Server) Primes(in *pb.PrimeRequest, stream pb.CalculatorService_PrimesServer) error {
	log.Printf("Primes was invoked with %v\n", in)

	ch := make(chan int64)
	go primeFactorization(in.Number, ch)
	for prime := range ch {
		stream.Send(&pb.PrimeResponse{
			Result: prime,
		})
	}

	return nil
}

func primeFactorization(N int64, ch chan<- int64) {
	defer close(ch)

	for i := int64(2); i*i < N; i++ {
		for N%i == 0 {
			ch <- i
			N /= i
		}
	}
	if N > 1 {
		ch <- N
	}
}