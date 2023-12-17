
# Swan Cloud Storage Example
[![Discord](https://img.shields.io/discord/770382203782692945?label=Discord&logo=Discord)](https://discord.gg/MSXGzVsSYf)
[![Twitter Follow](https://img.shields.io/twitter/follow/0xfilswan)](https://twitter.com/0xfilswan)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)

This is a example to how to build your own Web3 cloud storage by Swan Cloud Storage Solution, or integrate the Swan Cloud Storage to your applications.

--------
Swan Cloud Storage Solution consists of [Multi-Chain Storage(MCS)](https://multichain.storage) and [MetaArk](https://github.com/FogMeta/go-mc-sdk). 

[Multi-Chain Storage (MCS)]((https://multichain.storage)) represents a revolutionary Web3 S3 storage gateway fortified by smart contract integration to accelerate decentralized storage adoption.

[MetaArk](https://github.com/FogMeta/go-mc-sdk) is a robust data archiving and rebuilding service designed to securely store and automatically recover your most critical data assets. It encodes large datasets into CAR files, distributed seamlessly across decentralized storage networks. 


## Build Web3 Cloud Storage with MCS and MetaArk

Swan Network simplifies integrating decentralized storage into applications with [Multi-Chain Storage (MCS) SDK](https://github.com/filswan/go-mcs-sdk) and [MetaArk SDK](https://github.com/FogMeta/go-mc-sdk). The steps include:
 - Upload data to IPFS clusters via MCS for availability  
 - Back up to the Filecoin network using MetaArk SDK for redundancy
 - Rebuild from MetaArk if data is lost(If needed)

## Example

There is a simple case, you can integrate the [code](https://github.com/filswan/swan-storage/blob/main/examples/mcs-meta/main.go) to your applications to build your own Web3 Storage solution.


### Use Case Code
```go
package main

import (
	"fmt"
	"log"

	"github.com/FogMeta/go-mc-sdk/client"
	"github.com/filswan/go-mcs-sdk/mcs/api/bucket"
	"github.com/filswan/go-mcs-sdk/mcs/api/user"
)

// Acquire MCS api_key 
const (
	mcsAPIKey = "MCS_xxx"

	swanKey   = "xxxxxxxxxxxxx"
	swanToken = "xxxxxxxxxxxxxxxx"

	metaServer = ""
)

func main() {
	bucketName := "bucket_name"
	fileName := "file_name"
	filePath := "/file_path"

	// Initialize the MCS client
	mcsClient, err := user.LoginByApikeyV2(mcsAPIKey, "")
	if err != nil {
		log.Println(err)
		return
	}
	bucketClient := bucket.GetBucketClient(*mcsClient)
	// if no bucket,then create it
	_, err = bucketClient.CreateBucket(bucketName)
	if err != nil {
		log.Println(err)
		return
	}
	// upload file to your bucket
	if err = bucketClient.UploadFile(bucketName, fileName, filePath, true); err != nil {
		log.Println(err)
		return
	}
	// get file info from the bucket
	info, err := bucketClient.GetFile(bucketName, fileName)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("uploaded info :%#v\n", info)

	// get the download url for uploaded files
	gateway, err := bucketClient.GetGateway()
	if err != nil {
		log.Println(err)
		return
	}
	downloadURL := fmt.Sprintf("%s/ipfs/%s", *gateway, info.PayloadCid)

	// backup the file to Filecoin network by MetaArk
	datasetName := "my_dataset_name"
	metaClient := client.NewClient(swanKey, swanToken, &client.MetaConf{
		MetaServer: metaServer,
	})
	dataset_id, err := metaClient.Backup(datasetName, &client.IpfsData{
		DownloadUrl: downloadURL,
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("backup id :%#v\n", dataset_id)

	// list the backups info 
	data, err := metaClient.List(datasetName, 0, 10)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("backup list info :%#v\n", data)

	if len(data.DatasetList) == 0 {
		return
	}

	// if needed, you can rebuild dataset from filecoin network by MetaArk
	rebuildList, err := metaClient.Rebuild(dataset_id)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("rebuild info :%#v\n", rebuildList)
}

```

### How to get the MCS `API_key` and MetaArk `metaServer`
 - MCS `API_key`: Acquired from "https://www.multichain.storage" -> setting -> Create API Key
 - MetaArk `metaServer`: if you want to get a availabe `metaServer`, please open an [issue](https://github.com/filswan/swan-cloud-storage-example/issues) to contact with Swan Cloud team 



## Use Case
[Chainnode](https://chainnode.io) as a blockchain snapshots-as-a-service provider, ensuring swift downloads of the most up-to-date node snapshots and blockchain data are available in the future. The storage and backup functions are supported by the MCS and MetaArk.
 
## References
 - MCS SDKs
	- [Python SDK](https://github.com/filswan/python-mcs-sdk)
	- [Go SDK](https://github.com/filswan/go-mcs-sdk)
 - MetaArk SDKs
 	- [Go SDK](https://github.com/FogMeta/go-mc-sdk)

















