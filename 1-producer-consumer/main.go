//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

func producer(stream Stream, tweet_queue chan *Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(tweet_queue)
			return
		}

		tweet_queue <- tweet
	}
}

func consumer(tweet_queue chan *Tweet, done chan bool) {
	for {
		t, ok := <-tweet_queue
		if ok {
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		} else {
			done <- true
		}
		
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	queue := make(chan *Tweet)
	done := make(chan bool)

	// Producer
	go producer(stream, queue)

	// Consumer
	go consumer(queue, done)

	<- done

	fmt.Printf("Process took %s\n", time.Since(start))
}
