## spike-cli: a simple command line tool for Aerospike server

Aerospike does not come with a good GUI management tool. And it does not have anything for Windows.

spike-cli is a single binary that runs on both on Linux and Windows with no dependency. It's designed to help developers working with Aerospike. It can:

1. run asinfo commands. asinfo is an Aerospike command line util to retrieve server information.
1. scan a set to show the bin names that the set has
1. run query to show data
1. include test data to demo how to use this tool

### build and testing
```bash
# In Linux
go mod init github.com/brunoyin/testaero
go get github.com/aerospike/aerospike-client-go

env GOOS=linux go build
env GOOS=windows go build

# To test, start a test Aerospike server: the included Aerospike config file uses 64 MB file instead of default 4 GB.
mkdir $PWD/aerospike/data
docker run -d \
    -v $PWD/aerospike/data:/opt/aerospike/data \
    -v $PWD/aerospike/conf:/opt/aerospike/etc \
    --name aerospike -p 3000:3000 -p 3001:3001 -p 3002:3002 -p 3003:3003 \
    aerospike /usr/bin/asd --foreground --config-file /opt/aerospike/etc/aerospike.conf

go test ./spikeutils/

```

```powershell
# On Windows, in Powershell
$env:GOOS = 'linux'; go build

$env:GOOS = 'windows'; go build
```

### Usage
```powershell
.\spike-cli.exe -h

Usage of spike-cli.exe:
  -action string
        Run test, asinfo, scan or query, load (default "info")

   === info: default action ===
  -asinfo string
        asinfo commands like namespaces, sets, bins (default "namespaces")

  === load ===
  -csv string
        sample scorecard csv filename
  -namespace string
        Namespace (default "test")
  -namespace2 string
        Optional second Namespace to load data to

  === scan ===
  -namespace string
        Namespace (default "test")
  -set string
        Set name (default "scorecard")
   
  === query ===
  -namespace string
        Namespace (default "test")
  -set string
        Set name (default "scorecard")
  -bins string
        bin names for query (default "name,city,state")  
  -limit int
        Limit number of rows in a query, use 0 for no limit (default 10)
  
   === server connection: common to all action commands ===
  -host string
        Remote host (default "127.0.0.1")
  -port int
        Remote port (default 3000)
  -user string
        Aerospike user name
  -password string
        Aerospike password
```

### info

run asinfo commands: https://www.aerospike.com/docs/reference/info/

```powershell
# run asinfo show namespaces
.\spike-cli.exe
# same as
.\spike-cli.exe -action info -asinfo namespaces
# show sets
.\spike-cli.exe -asinfo sets
# show bins
.\spike-cli.exe -asinfo sets
# get server verion
.\spike-cli.exe -asinfo build
.\spike-cli.exe -asinfo node
.\spike-cli.exe -asinfo service
.\spike-cli.exe -asinfo namespace/test

```

### load test data

You do not need to load data to test if you already have data. It works only with scorecard-recent.csv .

```powershell
# load test data from scorecard-recent.csv to namespace test, set name scorecard
.\spike-cli.exe -action load -csv .\scorecard-recent.csv -namespace test
# load to 2 namespaces: test and test2
.\spike-cli.exe -action load -csv .\scorecard-recent.csv -namespace test -namespace2 test2

```

### scan

purpose is not to retrieve data but to discover bins in a set. Because Aerosipke is schema less, asinfo commands do not tell bin names in a set. 

The scan here is fast because it examine only one record assuming schema is enforced using client side programming.

```powershell
# to discover bins in set name scorecard, namespace test
.\spike-cli.exe -action scan # default namespace = test, default set name is scorecard
# to discover bins in set name scorecard, namespace test2
.\spike-cli.exe -action scan -namespace test2 -set scorecard
# Running scan on test2
# Bins found in namespace:  test2
# set name:  scorecard
#          state
#          zip
#          name
#          city
# 2021/03/06 01:36:37 Total records returned:  1
# 2021/03/06 01:36:37 Elapsed time:  0.0199532  seconds

# ========================
```

### query

View data. You will need namespace, set, and bin names to run a query. the default values work only with the test data provided.

By default, it limits 10 rows to return. You can use 0 to return all records.

```powershell
# Return 10 rows name, city, state from scorecard in test namespace
.\spike-cli.exe -action query

# Return 20 rows name, city from scorecard in test namespace
.\spike-cli.exe -action query -bins name,city,zip -limit 20

# Running query on test namespace, scorecard setname, limiting return to 20

# Bin names: name, city, zip
# ===========================

# 1:      The Art Institute of St Louis, St Charles, 63303,
# 2:      Auburn University at Montgomery, Montgomery, 36117-3596,
# 3:      Lewis-Clark State College, Lewiston, 83501-2698,
# 4:      Gannon University, Erie, 16541-0001,
# 5:      Galen College of Nursing-Tampa Bay, Saint Petersburg, 33716,
# 6:      Epic Bible College, Sacramento, 95841,
# 7:      Toccoa Falls College, Toccoa Falls, 30598,
# 8:      Los Angeles Mission College, Sylmar, 91342-3200,
# 9:      Antelope Valley College, Lancaster, 93536-5426,
# 10:     Eastfield College, Mesquite, 75150-2099,
# 11:     Sterling College, Sterling, 67579,
# 12:     Minnesota State Community and Technical College, Fergus Falls, 56537-1000,
# 13:     Brown Mackie College-Dallas, Bedford, 76021,
# 14:     University of Northwestern Ohio, Lima, 45805,
# 15:     University of Pittsburgh-Greensburg, Greensburg, 15601-5860,
# 16:     Ner Israel Rabbinical College, Baltimore, 21208,
# 17:     Argosy University-Phoenix Online Division, Phoenix, 85021,
# 18:     Cambridge College of Healthcare & Technology, Delray Beach, 33484,
# 19:     Canada College, Redwood City, 94061-1099,
# 20:     Stevens-Henager College, Murray, 84123-5671,
# ========================

```