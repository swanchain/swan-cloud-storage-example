package main

import (
	"fmt"
	"log"

	"github.com/FogMeta/go-mc-sdk/client"
	"github.com/filswan/go-mcs-sdk/mcs/api/bucket"
	"github.com/filswan/go-mcs-sdk/mcs/api/user"
)

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

	// upload file to mcs
	mcsClient, err := user.LoginByApikeyV2(mcsAPIKey, "")
	if err != nil {
		log.Println(err)
		return
	}
	bucketClient := bucket.GetBucketClient(*mcsClient)
	// if no bucket create it
	_, err = bucketClient.CreateBucket(bucketName)
	if err != nil {
		log.Println(err)
		return
	}
	// upload file to bucket
	if err = bucketClient.UploadFile(bucketName, fileName, filePath, true); err != nil {
		log.Println(err)
		return
	}
	// get file info from bucket
	info, err := bucketClient.GetFile(bucketName, fileName)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("uploaded info :%#v\n", info)

	// build download url for uploaded file
	gateway, err := bucketClient.GetGateway()
	if err != nil {
		log.Println(err)
		return
	}
	downloadURL := fmt.Sprintf("%s/ipfs/%s", *gateway, info.PayloadCid)

	// backup to meta ark
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

	// list backups
	data, err := metaClient.List(datasetName, 0, 10)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("backup list info :%#v\n", data)

	if len(data.DatasetList) == 0 {
		return
	}

	// rebuild dataset
	rebuildList, err := metaClient.Rebuild(dataset_id)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("rebuild info :%#v\n", rebuildList)
}
