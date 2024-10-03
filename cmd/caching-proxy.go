package cmd

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"main.go/api"
)

var port string = "3000"
var url string = "http://dummyjson.com/products"
var urlPort = make(map[string]string)

var caching = &cobra.Command{
	Use: "caching-proxy",
	Run: func(cmd *cobra.Command, args []string) {
		r := mux.NewRouter()
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatal(err)
		}
		url, err := cmd.Flags().GetString("origin")
		if err != nil {
			log.Fatal(err)
		}
		//  mapte tutuyoruz olan port ve linki
		urlPort[url] = port
		var srv = &http.Server{
			Addr:    port,
			Handler: r,
		}
		//json dosyası geldi
		jsonFile := getRequest(url)
		api.Router(r, jsonFile)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}

		// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		// defer cancel()
		// srv.Shutdown(ctx)
	},
}

// Request gönderiyor ve responsu okuyup gönderiyor
func getRequest(url string) []byte {
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
	return jsonFile
}

func init() {
	caching.Flags().StringVar(&port, "port", "3000", "--port <mumber>")
	caching.Flags().StringVar(&url, "origin", "", "--origin <url>")
}
