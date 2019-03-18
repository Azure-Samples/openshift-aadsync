package main

import (
	"flag"

	aadsyncclient "pkg/controller/client"
	controllerconfig "pkg/controller/config"
	mapper "pkg/controller/mapper"
	msgraphclient "pkg/msgraph/client"
	array "pkg/util/array"
	log "pkg/util/log"

	logrus "github.com/sirupsen/logrus"
)

var (
	logLevel = flag.String("loglevel", "Info", "Valid values are Debug, Info, Warning, Error")
)

func main() {

	flag.Parse()
	logrus.SetLevel(log.SanitizeLogLevel(*logLevel))
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	log := logrus.NewEntry(logrus.StandardLogger())
	log.Infof("####################################################")
	log.Infof("AAD Sync Controller")
	log.Infof("####################################################")

	controllerConfig, err := controllerconfig.NewControllerConfigFromFile()
	if err != nil {
		log.Fatal(err)
	}

	// Create clients to interact with the MS Graph API and the aadgroupsyncs.aad.microsoft.com CRDs
	msgraphClient := msgraphclient.NewClient(log)
	//aadsyncClient := aadsyncclient.NewClientForLocal(controllerConfig.Namespace, log)
	aadsyncClient := aadsyncclient.NewClient(controllerConfig.Namespace, log)

	// Create mapper to map between MS Graph and aadgroupsyncs.aad.microsoft.com constructs
	mapper := mapper.NewClient(log)

	// Create new or update existing aadgroupsyncs.aad.microsoft.com CRDs to match MS Graph
	// and Controller Config
	groupsToSync := controllerConfig.Groups
	for _, groupID := range groupsToSync {

		msGraphGroup, err := msgraphClient.GetGroup(groupID)
		if err != nil {
			log.Error(err)
		}

		log.Infof("----------------------------------------------------")
		log.Infof("Processing AAD Group ID from MS Graph: %s", msGraphGroup.ID)
		log.Infof("----------------------------------------------------")

		aadsyncGroup, err := aadsyncClient.Get(groupID)
		if err != nil {
			log.Error(err)
		}

		if aadsyncGroup == nil {

			log.Info("Creating new aadgroupsyncs.aad.microsoft.com entry")
			aadsyncGroup := mapper.CreateFromMSGraphGroup(msGraphGroup)
			_, err := aadsyncClient.Create(aadsyncGroup)
			if err != nil {
				log.Error(err)
			}

		} else {

			log.Info("Updating existing aadgroupsyncs.aad.microsoft.com entry")
			aadsyncGroup := mapper.UpdateFromMSGraphGroup(msGraphGroup, aadsyncGroup)
			_, err := aadsyncClient.Update(aadsyncGroup)
			if err != nil {
				log.Error(err)
			}
		}
	}

	log.Infof("----------------------------------------------------")
	log.Infof("Deleting non-configured aadgroupsyncs.aad.microsoft.com entries")
	log.Infof("----------------------------------------------------")

	// Delete all existing aadgroupsyncs.aad.microsoft.com CRDs that do not exist in Controller Config
	aadsyncGroupList, err := aadsyncClient.List()
	if err != nil {
		log.Error(err)
	}

	for _, group := range aadsyncGroupList {
		if !array.ContainsString(groupsToSync, group.ObjectMeta.Name) {
			err := aadsyncClient.Delete(group.ObjectMeta.Name)
			if err != nil {
				log.Error(err)
			}
		}
	}
}
