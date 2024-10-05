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
	"main.go/utils"
)

var port string = "3000"
var url string = "http://dummyjson.com/products"
var urlPort = make(map[string]string)

var caching = &cobra.Command{
	Use: "caching-proxy",
	Run: func(cmd *cobra.Command, args []string) {

		r := mux.NewRouter()

		// Flagleri aldım
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatal(err)
		}
		url, err := cmd.Flags().GetString("origin")
		if err != nil {
			log.Fatal(err)
		}

		// sonra olan dosyayı okuyup urlport mapine yazdım
		utils.Readfile(urlPort)

		// mapte kontrol edip olup olmadığını öğrendim
		isHave := utils.ControlMap(url, port, urlPort)

		if isHave {
			var localurl string = fmt.Sprintf("http://127.0.0.1:" + urlPort[url] + "/products")
			locJsonChan := make(chan []byte)
			go getRequest(locJsonChan, localurl)
			//jsonFile := <-locJsonChan
			fmt.Println("Localden geldi")
			return
		}

		//alınan dosyayı ilk başta yazdım
		utils.Writefile(url, port)
		var srv = &http.Server{
			Addr:    "127.0.0.1:" + port,
			Handler: r,
		}

		if err != nil {
			log.Fatal(err)
		}

		jsonChan := make(chan []byte)

		go getRequest(jsonChan, url) // İstek yaratıp  o urlye istek attım
		jsonFile := <-jsonChan       // json dosyasını alana kadar bekleyerek json dosyasını alıyor

		api.Router(r, jsonFile) // json dosyasını yazmaya hazırlanıyor

		go serviceandListen(srv, port) // portu açıyor
		fmt.Println("Serverdan geldi")
		go getSignal(srv)
		time.Sleep(20 * time.Second) // 5 dakika bekliyoruz localin kapanması için
		// bu komut bütünü de serveri kapatıyor
		utils.Deletefile()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP shutdown error: %v", err)
		}

	},
}

func serviceandListen(srv *http.Server, port string) {
	fmt.Printf("Server is running in %v port", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func getSignal(srv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop

	utils.Deletefile()
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
	fmt.Println("istek geldi")
	jsonChan <- jsonFile
}

func init() {
	caching.Flags().StringVar(&port, "port", "3000", "--port <mumber>")
	caching.Flags().StringVar(&url, "origin", "", "--origin <url>")
}
