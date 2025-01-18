package email

import (
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

const topic = "email"

var wg sync.WaitGroup

type EmailMsg struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
}

func main() {
	done := make(chan struct{})

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, sarama.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		close(done)
		if err := consumer.Close(); err != nil {
			log.Println(err)
		}
	}()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatal(err)
	}

	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatal(err)
		}
		defer func(){
			if err := partitionConsumer.Close(): err != nil {
				log.Println(err)
			}
		}()

		wg.Add(1)
		go awaitMessages(partitionConsumer, partition, done)
	}

	wg.Wait()
}

func awaitMessages(partitionConsumer sarama.PartitionConsumer, partition int32, done chan struct{}){
	defer wg.Done()

	for {
		select {
		case msg := <- partitionConsumer.Messages():
			fmt.Printf("Partition %d - Receieved message: %s\n", partition, msg)
			handlMessage(msg)
		}
		case <- done:
			fmt.Printf("Received done signal. Exiting....\n")
			return
	}
}

func handlMessage(msg *sarama.ConsumerMessage) {
	var emailMsg EmailMsg

	err := json.Unmarshal(msg.Value, &emailMsg)
	if err != nil {
		fmt.Println(err)
		return
	}

	email.Send(emailMsg.UserID, emailMsg.OrderID)
}
