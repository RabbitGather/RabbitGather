/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(cmd *cobra.Command, args []string) {
	conn, err := amqp.Dial("amqp://sayken:kingkingjin@rabbitmq:5672/")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	//failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go consume(ch, q)

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatal(err)
	}
	//failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
	var forever chan struct{}
	<-forever
	//	a echo server on port 8080
	//mx := http.NewServeMux()
	//mx.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	// response everything received
	//	b, err := io.ReadAll(r.Body)
	//	if err != nil {
	//		log.Println(err)
	//		w.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//
	//	b = append(b, []byte("BBBBBBBBBBB")...)
	//	_, err = w.Write(b)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//})
	//err := http.ListenAndServe(":3001", mx)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
}

func consume(ch *amqp.Channel, q amqp.Queue) {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
	}
}
