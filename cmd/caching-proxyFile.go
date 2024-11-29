package cmd

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	_ "github.com/spf13/cobra/cobra/cmd"
	"log"
	"main.go/api"
	"main.go/utils"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var cachingFile = &cobra.Command{
	Use: "caching-proxyFile",
	Run: cachingFileFunc(),
}
func cachingFileFunc() func(cmd *cobra.Command,args []string) {
	return func(cmd *cobra.Command, args []string) {
		// servis ve handler oluşturuldu.
		r := mux.NewRouter()
		var srv = &http.Server{
			Addr:    "127.0.0.1:" + port,
			Handler: r,
		}
		port, err := cmd.Flags().GetString("port")
		if err !=nil{
			log.Fatal(err)
		}
		url, err := cmd.Flags().GetString("origin")
		if err != nil {
			log.Fatal(err)
		}
		readerPort,err := utils.Readfile(port,url)
		if err !=nil{
			fmt.Println(err)
			return
		}
		//eğer port varsa hiç uğraşma
		if readerPort != ""{
			fmt.Println("X-Cache HIT")
			var localurl string = fmt.Sprintf("http://127.0.0.1:" + readerPort)
			var localJsonChan = make(chan []byte)
			go getRequest(localJsonChan,localurl)
			<- localJsonChan
			return
		}
		err = utils.Writefile(url,port)
		if err!=nil{
			log.Fatal(err)
		}
		// isteği attık
		var jsonGetChan = make(chan []byte)
		go getRequest(jsonGetChan,url)
		//isteği bekledik
		jsonFile := <- jsonGetChan

		api.Router(r,jsonFile)
		go serviceandListen(srv,port)
		go getSignalFile(srv,port,url)
		shutdownServeFile(srv)
	}
}
func shutdownServeFile(srv *http.Server){
	time.Sleep(2 * time.Minute) // 2 dakika bekliyoruz localin kapanması için

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := utils.Deletefile(port,url)
	if err !=nil{
		log.Fatal(err)
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}
func getSignalFile(srv *http.Server,port string ,url string){
	fmt.Println("X-CACHE : MISS")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop
	err := utils.Deletefile(port,url)
	if err!=nil{
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}

func init() {
	caching.Flags().StringVar(&port, "port", "3000", "--port <mumber>")
	caching.Flags().StringVar(&url, "origin", "", "--origin <url>")
}