package main

import (
	"net/http"
	"log"
	"github.com/G-Node/gin-doi/src"
	"github.com/docopt/docopt-go"
	"os"
	"strconv"
)

func main() {
	usage := `gin doi.
Usage:
  gin-doi [--max_workers=<max_workers> --max_queue_size=<max_queue_size> --port=<port> --source=<source>
           --target=<target> --storeURL=<url> --mServer=<server> --mFrom=<from> --doiMaster=<master> --doiBase=<base>
           --sendMail]

Options:
  --max_workers=<max_workers>     The number of workers to start [default: 3]
  --max_queue_size=<max_quesize>  The The size of the job queue [default: 100]
  --port=<port>                   The server port [default: 8083]
  --source=<dsourceurl>           The Server adress from which data can be read [default: https://repo.gin.g-node.org/]
  --target=<target>               The Location for long term storgae [default: data]
  --storeURL=<url>                The base url for storage [default: http://doid.gin.g-node.org/]
  --mServer=<server>              The mailserver adress (:and port) [default: localhost:25]
  --mFrom=<from>                  The mail from adress [default: no-reply@g-node.org]
  --doiMaster=<master>            The mail adress to send info to [default: dev@g-node.org]
  --doiBase=<base>                The first part of the DOI [default: 10.12751]
  --sendMail                      Whether Mail Noticiations should really be send (Otherwise just print them)
 `

	args, err := docopt.Parse(usage, nil, true, "gin doi 0.1a", false)
	if err != nil {
		log.Printf("Error while parsing command line: %+v", err)
		os.Exit(-1)
	}
	ds := ginDoi.GinDataSource{GinURL: args["--source"].(string)}
	dp := ginDoi.DoiProvider{ApiURI:"", DOIBase:args["--doiBase"].(string)}
	mServer := ginDoi.MailServer{Adress: args["--mServer"].(string), From: args["--mFrom"].(string),
		DoSend:args["--sendMail"].(bool),
		Master: args["--doiMaster"].(string)}
	storage := ginDoi.LocalStorage{Path:args["--target"].(string), Source:ds, HttpBase:args["--storeURL"].(string),
					DProvider:dp, MServer: &mServer}
	op := ginDoi.OauthProvider{Uri:"https://auth.gin.g-node.org/api/accounts"}
	// Create the job queue.
	maxQ,err :=strconv.Atoi(args["--max_queue_size"].(string))
	if err != nil {
		log.Printf("Error while parsing command line: %+v", err)
		os.Exit(-1)
	}
	jobQueue := make(chan ginDoi.Job, maxQ)
	// Start the dispatcher.
	maxW,err :=strconv.Atoi(args["--max_workers"].(string))
	dispatcher := ginDoi.NewDispatcher(jobQueue, maxW)
	dispatcher.Run(ginDoi.NewWorker)

	// Start the HTTP handler.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ginDoi.InitDoiJob(w, r, &ds, &op)
	})
	http.HandleFunc("/do/", func(w http.ResponseWriter, r *http.Request) {
		ginDoi.DoDoiJob(w,r,jobQueue, storage, &op)
	})
	http.Handle("/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir("/assets"))))

	log.Fatal(http.ListenAndServe(":"+args["--port"].(string), nil))
}
