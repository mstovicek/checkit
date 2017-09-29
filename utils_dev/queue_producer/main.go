package main

import (
	"fmt"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/messaging"
)

func main() {
	log := logger.NewStdout().SetDebug(true).SetComponent("utils_dev")

	fmt.Println("a")
	commitReceivedProducer, _ := messaging.NewRabbitMqCommitReceivedPublisher(messaging.RabbitMqURL, log)
	fmt.Println("d")
	defer commitReceivedProducer.Close()
	fmt.Println("b")
	uuid := messaging.NewUUIDGenerator()

	commitReceivedProducer.Publish(
		uuid.Generate(),
		messaging.CommitReceived{
			CommitHash:     "22bde6193c1e313bc71668b74652634bcbd46a49",
			Server:         "github",
			RepositoryName: "mstovicek/checkit-demo",
			Files: []string{
				"files/a/test.php",
				"files/b/test.php",
			},
		},
	)
	fmt.Println("c")

	//commitReceivedProducer.Publish(
	//	uuid.Generate(),
	//	messaging.CommitReceived{
	//		CommitHash: "hello 7s .......",
	//		Server: "github",
	//		RepositoryName: "mstovicek/checkit",
	//	},
	//)
	//
	//commitReceivedProducer.Publish(
	//	uuid.Generate(),
	//	messaging.CommitReceived{
	//		CommitHash: "hello 6s ......",
	//		Server: "github",
	//		RepositoryName: "mstovicek/checkit",
	//	},
	//)
	//
	//commitReceivedProducer.Publish(
	//	uuid.Generate(),
	//	messaging.CommitReceived{
	//		CommitHash: "hello 5s .....",
	//		Server: "github",
	//		RepositoryName: "mstovicek/checkit",
	//	},
	//)
	//
	//commitReceivedProducer.Publish(
	//	uuid.Generate(),
	//	messaging.CommitReceived{
	//		CommitHash: "hello 4s ....",
	//		Server: "github",
	//		RepositoryName: "mstovicek/checkit",
	//	},
	//)
	//
	//commitReceivedProducer.Publish(
	//	uuid.Generate(),
	//	messaging.CommitReceived{
	//		CommitHash: "hello 3s ...",
	//		Server: "github",
	//		RepositoryName: "mstovicek/checkit",
	//	},
	//)
}
