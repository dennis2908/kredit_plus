package pusherconn

import (
	"os"

	pusher "github.com/pusher/pusher-http-go/v5"
)

var (
	pusherClient *pusher.Client
)

func Connect() (pusher.Client, error) {

	pusherClient := pusher.Client{

		AppID: os.Getenv("pusher_appId"),

		Key: os.Getenv("pusher_key"),

		Secret: os.Getenv("pusher_secret"),

		Cluster: os.Getenv("pusher_cluster"),

		Secure: true,
	}
	return pusherClient, nil

}
