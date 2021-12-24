package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sandertv/mcwss"
	"github.com/sandertv/mcwss/protocol/event"
)

// botのトークンと連携するチャンネルのIDを入力する
const (
	TOKEN   = "NjI3NzU0NTI5ODAyMjg5MTYy.XZBP9w.q2lHCV--FDh9CeZGmepHLAd2W3g"
	CHANNEL = "825193887919767583"
)

var mcwsclient *mcwss.Player = nil

func main() {
	address := flag.String("addr", "localhost", "address")
	port := flag.Int("port", 8000, "port number")
	flag.Parse()

	server := mcwss.NewServer(&mcwss.Config{
		HandlerPattern: "/ws",
		Address:        fmt.Sprintf("%s:%d", *address, *port),
	})

	bot := BotLogin(TOKEN)
	defer bot.Close()

	server.OnConnection(func(player *mcwss.Player) {
		mcwsclient = player
		log.Printf("%sが接続しました", mcwsclient.Name())
		player.OnPlayerMessage(func(event *event.PlayerMessage) {
			if event.MessageType != "me" && event.MessageType != "chat" {
				return
			}
			log.Printf("[Minecraft]<%s> %s", event.Sender, event.Message)
			bot.ChannelMessageSend(CHANNEL, fmt.Sprintf("<%s> %s", event.Sender, event.Message))
		})
	})
	server.OnDisconnection(func(player *mcwss.Player) {
		log.Printf("%sが切断しました", player.Name())
		mcwsclient = nil
	})

	// 中断処理
	c := make(chan os.Signal)
	defer close(c)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-c
		log.Println("サーバーを終了します")
		bot.Close()
		close(c)
		os.Exit(0)
	}()

	// サーバーを開始
	log.Printf("サーバーを開始します（ポート番号: %d）", *port)
	server.Run()
}
