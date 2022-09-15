package PubsubHelper

import (
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
	"ClipboardShare/SqlHelper"
	"context"
	"time"
	"sync"	
	"os"
)

var unsubscribeTopic context.CancelFunc
var client *pubsub.Client
var	proj = "clipboardshare-357916"
var	defaultTopicName = "my-topic"
var	defaultSubName, _ = os.Hostname()
var defaultTopic *pubsub.Topic
var defaultSub *pubsub.Subscription
var dbCredsKey = "GoogleCreds"

/**
 * Panic if the checked error is not 'nil'
 */
func check (err error) {
	if (err != nil) {
		panic(err)
	}
}

/* Changes the default hard coded configuration 
 * input 
 		* authfilePath - json file with authentication details path
*/
func SetPubSubConfig(authfilePath string) {

	// update configuration
	dat, err := os.ReadFile(authfilePath)
	check(err)
	SqlHelper.InsertConfig(SqlHelper.Config{Key: dbCredsKey, Value: string(dat)})
}

/* Create a new topic if it doesn't exist - if exist, return the existing one */
func createTopicIfNotExists(client *pubsub.Client, topic string) *pubsub.Topic {
	ctx := context.Background()

	// get list of topics 
	t := client.Topic(topic)
	exist, err := t.Exists(ctx)
	check(err)
	
	// if the topic already exist
	if exist {
		return t
	}

	// if the topic doesn't exist - create a new one
	t, err = client.CreateTopic(ctx, topic)
	check(err)

	return t
}

/* Create a subscription if it doesn't exist - the subscription name is the hostname of the device */
func createSubscriptionIfNotExists(client *pubsub.Client, name string, topic *pubsub.Topic) *pubsub.Subscription {
	ctx := context.Background()

	sub := client.Subscription(name)
	exist, err := sub.Exists(ctx)
	check(err)

	if exist {
		return sub
	}

	sub, err = client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig {
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	check(err)
	
	return sub
}

/* subscribe to any topic, you should pass a channel to which the new messages will be added */
func SubscribeToTopic(client *pubsub.Client, sub *pubsub.Subscription, topic *pubsub.Topic, messageChan chan []byte) {
	ctx := context.Background()

	var mu sync.Mutex

	cctx, cancel := context.WithCancel(ctx)
	unsubscribeTopic = cancel

	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {

		// acknoledge message received 
		msg.Ack()
		
		// mutex lock
		mu.Lock()
			
		// push the message into the channel
		messageChan <- msg.Data
		
		// mutex unlock
		defer mu.Unlock()
	})
	check(err)
}

/* publish to any topic */
func publishToTopic(client *pubsub.Client, topic *pubsub.Topic, msg []byte) {
	ctx := context.Background()
	result := topic.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	_, err := result.Get(ctx)
	check(err)
}

/* publish to the default topic (hard coded in the code but can be changed using `setPubSubConfig` function) */
func DefaultPublish(msg []byte) {
	publishToTopic(client, defaultTopic, msg)
}

/* subscribe to the default topic (hard coded in the code but can be changed using `setPubSubConfig` function) */
func DefaultSubscribe(messageChan chan []byte) {
	SubscribeToTopic(client, defaultSub, defaultTopic, messageChan)
}


func init() {
	initializeClient()
}

/* Initialize the client and connect to the Google Services */
func initializeClient() {
	ctx := context.Background()
	tmpClient, err := pubsub.NewClient(ctx, proj, option.WithCredentialsJSON([]byte(SqlHelper.QueryConfigByKey(dbCredsKey).Value)))
	if (err != nil) {
		return
	}

	// check(err)

	// store the client connection
	client = tmpClient

	// create a new topic 
	defaultTopic = createTopicIfNotExists(client, defaultTopicName)

	// Create a new subscription.
	defaultSub = createSubscriptionIfNotExists(client, defaultSubName, defaultTopic)
}
