package main

import "fmt"

type BenchmarkTest func(ReaderFactory)

func TestReaderBehavior1(readerFactory ReaderFactory) {
	fmt.Println("Running TestReaderBehavior1")
}

func TestReaderBehavior2(readerFactory ReaderFactory) {
	fmt.Println("Running TestReaderBehavior2")
}

func TestReaderBehavior3(readerFactory ReaderFactory) {
	fmt.Println("Running TestReaderBehavior3")
}
