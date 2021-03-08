package spikeutils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	aero "github.com/aerospike/aerospike-client-go"
)

// PanicOnError logs error before fail
func PanicOnError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func readData(fileName string) ([][]string, error) {

	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	// get headers
	// var headers []string
	// if headers, err = r.Read(); err != nil {
	// 	return [][]string{}, err
	// }
	// fmt.Println(headers)

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

// LoadData loading scorecard.csv file for testing
func LoadData(client *aero.Client, wpolicy *aero.WritePolicy, fileName string, setName string, ns ...string) {
	records, err := readData(fileName)
	PanicOnError(err)
	headers := records[0]
	fmt.Println(headers)
	fmt.Println("==========================================")
	for i, header := range headers {
		fmt.Printf("%d\t%s = %s\n", i, header, records[1][i])
	}
	for _, testns := range ns {
		checked := false
		for i, row := range records[1:] {
			key, err := aero.NewKey(testns, setName, row[0])
			PanicOnError(err)
			if !checked {
				_, err = client.Get(nil, key)
				if err == nil {
					fmt.Println("data already loaded. exiting")
					break
				}
				checked = true
			}

			// client.PutBins(wpolicy, key, aero.NewBin("name", row[1]), aero.NewBin("city", row[2]), aero.NewBin("state", row[3]), aero.NewBin("zip", row[4]))
			bins := aero.BinMap{
				"name":  row[1],
				"city":  row[2],
				"state": row[3],
				"zip":   row[4],
			}
			client.Put(wpolicy, key, bins)
			if i%50 == 0 {
				log.Printf("%d records loaded ....\n", i+1)
			}
		}
	}
}

// Info runs asinfo command
func Info(clientPolicy *aero.ClientPolicy, host string, port int, asinfo string) {
	conn, err := aero.NewConnection(clientPolicy, aero.NewHost(host, port))
	if err != nil {
		log.Fatalln(err.Error())
	}

	infoMap, err := aero.RequestInfo(conn, asinfo)
	if err != nil {
		log.Fatalln(err.Error())
	}

	cnt := 1
	for k, v := range infoMap {
		switch k {
		case "namespaces":
			fmt.Println("namespaces found: ", strings.TrimRight(v, ";"))
		case "sets":
			fmt.Println("sets found:")
			for _, line := range strings.Split(strings.TrimRight(v, ";"), ";") {
				tags := strings.Split(line, ":")
				fmt.Println("namespace: ", strings.SplitN(tags[0], "=", 2)[1])
				fmt.Println("\t", strings.SplitN(tags[1], "=", 2)[1])
			}
		case "bins":
			fmt.Println("Bins found:")
			for _, line := range strings.Split(strings.TrimRight(v, ";"), ";") {
				tags := strings.Split(line, ":")
				fmt.Println("namespace: ", tags[0])
				for i, v := range strings.Split(tags[1], ",")[2:] {
					fmt.Printf("\t%d: %s\n", i, v)
				}
			}
		default:
			fmt.Printf("%d :  %s\n     %s\n\n", cnt, k, v)
		}
		cnt++
	}
}

// Query scorecard test data
func Query(client *aero.Client, ns string, setName string, binNames []string, limit int) {
	// client := GetClient()
	// defer client.Close()
	stmt := aero.NewStatement(ns, setName, binNames...)
	rs, err := client.Query(nil, stmt)
	PanicOnError(err)

	fmt.Printf("Bin names: %s\n===========================\n", strings.Join(binNames, ", "))
	i := 0
	for res := range rs.Results() {
		i++
		if limit != 0 && i > limit {
			break
		}
		if res.Err != nil {
			// handle error here
			// if you want to exit, cancel the recordset to release the resources
			rs.Close()
			PanicOnError(res.Err)
		} else {
			// process record here
			fmt.Printf("\n%d:\t", i) //, res.Record.Bins)
			for _, name := range binNames {
				fmt.Printf("%s, ", res.Record.Bins[name])
			}
		}
	}
}

// Scan to discover bin names in a set.
// This scan does not do thorough scan every record. It takes the first record, and returns
func Scan(client *aero.Client, namespace string, setName string) {
	// client := GetClient()
	// defer client.Close()
	recordCount := 0
	begin := time.Now()
	policy := aero.NewScanPolicy()
	recordset, err := client.ScanNode(policy, client.Cluster().GetNodes()[0], namespace, setName)
	PanicOnError(err)

	fmt.Println("Bins found in namespace: ", namespace, "\nset name: ", setName)
	for rec := range recordset.Results() {
		if rec.Err != nil {
			// if there was an error, stop
			PanicOnError(err)
		}
		for k := range rec.Record.Bins {
			fmt.Println("\t", k)
		}
		// break

		recordCount++

		if (recordCount % 100000) == 0 {
			log.Println("Records ", recordCount)
		}
		break
	}

	end := time.Now()
	seconds := float64(end.Sub(begin)) / float64(time.Second)
	log.Println("Total records returned: ", recordCount)
	log.Println("Elapsed time: ", seconds, " seconds")
}

// GetClient get a client connected to the target aerospike cluster
func GetClient(clientPolicy *aero.ClientPolicy, host string, port int) (*aero.Client, error) {
	client, err := aero.NewClientWithPolicy(clientPolicy, host, port)
	return client, err
}

func init() {
	// use all cpus in the system for concurrency
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetOutput(os.Stdout)

}
