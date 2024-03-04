package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

var locations map[int32]string

func main() {

	locations = make(map[int32]string)

	seed := flag.String("seed", "", "seed server")
	user := flag.String("username", "", "username")
	pass := flag.String("password", "", "password")
	useSha512 := flag.Bool("use512", false, "whether to use SCRAM_SHA_512")

	flag.Parse()

	seeds := []string{*seed}
	var opts []kgo.Opt
	opts = append(opts,
		kgo.SeedBrokers(seeds...),
	)

	ctx := context.Background()

	// Initialize public CAs for TLS
	opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))

	if *useSha512 {
		// Initializes SASL/SCRAM 512
		opts = append(opts, kgo.SASL(scram.Auth{
			User: *user,
			Pass: *pass,
		}.AsSha512Mechanism()))
	} else {
		// Initializes SASL/SCRAM 256
		opts = append(opts, kgo.SASL(scram.Auth{
			User: *user,
			Pass: *pass,
		}.AsSha256Mechanism()))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	admin := kadm.NewClient(client)
	brokerMetadata, err := admin.BrokerMetadata(ctx)
	if err != nil {
		panic(err)
	}
	defer admin.Close()

	for _, broker := range brokerMetadata.Brokers {
		locations[broker.NodeID] = *broker.Rack
	}

	resp, err := admin.ListTopicsWithInternal(ctx)
	if err != nil {
		panic(err)
	}
	resp.EachPartition(process)
}

func process(pd kadm.PartitionDetail) {
	topicLocations := make(map[string]int)
	for _, replicaID := range pd.Replicas {
		topicLocations[locations[replicaID]]++
	}
	for location, count := range topicLocations {
		if count > 1 {
			fmt.Printf("%s:%v has more than 1 replica in %s\n", pd.Topic, pd.Partition, location)
		}
	}
}
