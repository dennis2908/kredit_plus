package pusherconn

import (
	pusher "github.com/pusher/pusher-http-go/v5"
)

var (
	pusherClient *pusher.Client
)

func Connect() (pusher.Client, error) {
	pusherClient := pusher.Client{

		AppID: "PUSHER_APP_ID",

		Key: "PUSHER_APP_KEY",

		Secret: "PUSHER_APP_SECRET",

		Cluster: "PUSHER_APP_CLUSTER",

		Secure: true,
	}
	return pusherClient, nil

}
