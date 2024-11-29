package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"main.go/api"
	"main.go/handlers"
	"main.go/model"
)

var port string = "3000"
var url string = "http://dummyjson.com/products"

var caching = &cobra.Command{
	Use: "caching-proxy",
	Run: func(cmd *cobra.Command, args []string) {

		r := mux.NewRouter()

		var db = model.Database{
			Username: "grkmkly35",
			Server:   "",
		}
		handlers.Connect(&db)

		chanCollection := make(chan bool)
		go handlers.SetCollection(&db, "linkport", chanCollection)
		// Flagleri aldım
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatal(err)
		}
		url, err := cmd.Flags().GetString("origin")
		if err != nil {
			log.Fatal(err)
		}
		item := model.LinkPort{
			Link: url,
			Port: port,
		}

		<-chanCollection

		isHave, portWarning := handlers.CheckLinkPort(&db, item)
		fmt.Println("MERHABA : ", isHave)
		if portWarning == "ACTIVE" && !isHave {
			fmt.Println("PORT :", portWarning, "please change port") // port aktif uyarısı varsa portu değiştir uyarısı veriyor
			return
		}
		if isHave {
			var localurl string = fmt.Sprintf("http://127.0.0.1:" + portWarning)
			locJsonChan := make(chan []byte)
			// cache istek gönderdim
			go getRequest(locJsonChan, localurl)
			<-locJsonChan
			fmt.Println("X-CACHE : HIT")
			fmt.Println(locJsonChan)
			return
		}

		handlers.InsertLinkPort(&db, &item)
		var srv = &http.Server{
			Addr:    "127.0.0.1:" + port,
			Handler: r,
		}
		jsonChan := make(chan []byte)

		go getRequest(jsonChan, url) // İstek yaratıp  o urlye istek attım
		jsonFile := <-jsonChan       // json dosyasını alana kadar bekleyerek json dosyasını alıyor

		api.Router(r, jsonFile) // json dosyasını yazmaya hazırlanıyor

		go serviceandListen(srv, port) // portu açıyor

		go getSignal(&db, item, srv)
		shutdownServe(&db, item, srv)
	},
}

func shutdownServe(db *model.Database, item model.LinkPort, srv *http.Server) {
	time.Sleep(2 * time.Minute) // 2 dakika bekliyoruz localin kapanması için

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	handlers.DeleteLinkPort(db, item)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}

func serviceandListen(srv *http.Server, port string) {
	fmt.Printf("Server is running in %v port", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func getSignal(db *model.Database, item model.LinkPort, srv *http.Server) {
	fmt.Println("X-CACHE : MISS")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop
	handlers.DeleteLinkPort(db, item)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}

// Request gönderiyor ve responsu okuyup gönderiyor
func getRequest(jsonChan chan []byte, url string) {
	req, err := http.NewRequest(http.MethodGet, url, nil) // request oluşturuyor fakat requesti göndermiyor hemen
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req) // oluşturulan requesti yaptırdık ve bir response döndü
	if err != nil {
		log.Fatal(err)
	}
	jsonFile, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	jsonChan <- jsonFile
}

func init() {
	cachingFile.Flags().StringVar(&port, "port", "3000", "--port <mumber>")
	cachingFile.Flags().StringVar(&url, "origin", "", "--origin <url>")
}
