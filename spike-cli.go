package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	aero "github.com/aerospike/aerospike-client-go"
	"github.com/brunoyin/spike-cli/spikeutils"
)

var (
	host         = "127.0.0.1"
	port         = 3000
	namespace    = "test"
	namespace2   = ""
	setName      = "scorecard"
	username     = ""
	password     = ""
	wpolicy      = aero.NewWritePolicy(0, 0)
	clientPolicy = aero.NewClientPolicy()
	binList      = "name,city,state"
)

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	// log.SetOutput(os.Stdout)
	log.Println("started")
	t1 := time.Now()
	//
	asinfoCmd := "namespaces"
	// test := flag.Bool("test", false, "Run a simple test of create, read and delete ops")
	action := flag.String("action", "info", "Run asinfo, scan or query, load")
	// ns := flag.String("ns", "test", "Aerospike namespace")
	// namedSet := flag.String("setname", "scorecard", "Aerospike set name")
	csvfile := flag.String("csv", "", "sample scorecard csv filename")
	limit := flag.Int("limit", 10, "Limit number of rows in a query, use 0 for no limit")
	// bins := flag.
	flag.StringVar(&host, "host", host, "Remote host")
	flag.IntVar(&port, "port", port, "Remote port")
	flag.StringVar(&namespace, "namespace", namespace, "Namespace")
	flag.StringVar(&namespace2, "namespace2", namespace2, "Second Namespace to load data to")
	flag.StringVar(&setName, "set", setName, "Set name")
	flag.StringVar(&asinfoCmd, "asinfo", asinfoCmd, "asinfo commands like namespaces, sets, bins")
	flag.StringVar(&username, "user", username, "Aerospike user name")
	flag.StringVar(&password, "password", password, "Aerospike password")
	flag.StringVar(&binList, "bins", binList, "bin names for query")
	// loadTestdata := flag.Bool("ns2", false, "load test data")

	flag.Parse()
	if username != "" {
		clientPolicy.User = username
		clientPolicy.Password = password
	}
	clientPolicy.Timeout = 15 * time.Second
	client, err := spikeutils.GetClient(clientPolicy, host, port)
	if err != nil {
		spikeutils.PanicOnError(err)
	}
	defer client.Close()

	switch *action {
	case "info":
		fmt.Println("Running asinfo ...")
		spikeutils.Info(client, host, port, asinfoCmd)
	case "scan":
		fmt.Printf("Running scan on %s\n", namespace)
		spikeutils.Scan(client, namespace, setName)
	case "query":
		fmt.Printf("Running query on %s namespace, %s setname, limiting return to %d\n\n", namespace, setName, *limit)
		spikeutils.Query(client, namespace, setName, strings.Split(binList, ","), *limit)
	case "load":
		fmt.Printf("Loading data file %s to %s namespace , %s setname", *csvfile, namespace, setName)

		if namespace2 != "" {
			fmt.Printf("Also loading data file %s to %s namespace , %s setname", *csvfile, namespace2, setName)
			spikeutils.LoadData(client, wpolicy, *csvfile, setName, namespace, namespace2)
		} else {
			spikeutils.LoadData(client, wpolicy, *csvfile, setName, namespace)
		}

	default:
		fmt.Printf("%s is an invalid action", *action)
	}
	fmt.Println("\n========================")
	log.Println("ended")
	t2 := time.Now().UnixNano() - t1.UnixNano()
	fmt.Printf("\n\n%f seconds used", float64(t2)/1000000000.0)
}
