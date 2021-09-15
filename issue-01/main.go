package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

const DeliverySubject = "ORDER_REPORTING"
const Inbox = "ORDER_REPORTING"

func InitPredefinedStreamAndConsumer() (*nats.Conn, error) {
	nc, err := nats.Connect("127.0.0.1:4222")
	if err != nil {
		return nil, err
	}
	mgr, _ := jsm.New(nc, jsm.WithTimeout(10*time.Second))

	streamName := "ORDER"

	// init stream
	stream, err := mgr.LoadOrNewStream(
		streamName, jsm.Subjects("ORDER.>"), jsm.FileStorage(), jsm.DiscardOld(),
		jsm.LimitsRetention(), jsm.Replicas(2), jsm.MaxAge(24*time.Hour),
		jsm.MaxMessages(5000000),
	)

	if err != nil {
		return nil, err
	}

	// purge all message
	if err = stream.Purge(); err != nil {
		return nil, err
	}

	// init consumer for stream defined above
	_, err = mgr.LoadOrNewConsumer(
		streamName, DeliverySubject, jsm.DurableName(DeliverySubject), jsm.AcknowledgeExplicit(),
		jsm.DeliverySubject(Inbox), jsm.ReplayAsReceived(),
		jsm.DeliverAllAvailable(),
	)

	if err != nil {
		return nil, err
	}

	return nc, nil
}

func main() {
	nc0, err := InitPredefinedStreamAndConsumer()
	if err != nil {
		panic(err)
	}

	// 3 subscriber with deliverySubject
	// No.1
	subs := make([]*nats.Subscription, 3)
	defer func() {
		for _, sub := range subs {
			sub.Unsubscribe()
		}
	}()

	subscription, err := nc0.Subscribe(DeliverySubject, func(msg *nats.Msg) {
		fmt.Printf("consumer[0] [%s] received message: %v\n", time.Now().Format("2006-01-02 15:04:05"), msg.Data)
	})

	subs = append(subs, subscription)
	if err != nil {
		fmt.Println("consumer[0] subscribe error", err)
	}

	nc1, err := InitPredefinedStreamAndConsumer()
	if err != nil {
		panic(err)
	}

	// No.2
	subscription, err = nc1.Subscribe(DeliverySubject, func(msg *nats.Msg) {
		fmt.Printf("consumer[1] [%s] received message: %v\n", time.Now().Format("2006-01-02 15:04:05"), msg.Data)
	})

	subs = append(subs, subscription)
	if err != nil {
		fmt.Println("consumer[1] subscribe error", err)
	}

	nc2, err := InitPredefinedStreamAndConsumer()
	if err != nil {
		panic(err)
	}

	// No.3
	subscription, err = nc2.Subscribe(DeliverySubject, func(msg *nats.Msg) {
		fmt.Printf("consumer[2] [%s] received message: %v\n", time.Now().Format("2006-01-02 15:04:05"), msg.Data)
	})

	subs = append(subs, subscription)
	if err != nil {
		fmt.Println("consumer[2] subscribe error", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		nc3, err := InitPredefinedStreamAndConsumer()
		if err != nil {
			panic(err)
		}
		js, err := nc3.JetStream()
		if err != nil {
			panic(err)
		}

		for i := 0; i < 10; i++ {
			_, err = js.Publish("ORDER.1", []byte{byte(i)})
			if err != nil {
				fmt.Println("Publisher error: ", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}
