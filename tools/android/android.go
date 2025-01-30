package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"google.golang.org/api/androidmanagement/v1"
	"google.golang.org/api/option"
)

// Required env vars:
var (
	androidServiceCredentials = os.Getenv("FLEET_ANDROID_SERVICE_CREDENTIALS")
	androidProjectID          = os.Getenv("FLEET_ANDROID_PROJECT_ID")
)

func main() {
	if androidServiceCredentials == "" || androidProjectID == "" {
		log.Fatal("FLEET_ANDROID_SERVICE_CREDENTIALS and FLEET_ANDROID_PROJECT_ID must be set")
	}

	command := flag.String("command", "", "")
	enterpriseID := flag.String("enterprise_id", "", "")
	flag.Parse()

	ctx := context.Background()
	mgmt, err := androidmanagement.NewService(ctx, option.WithCredentialsJSON([]byte(androidServiceCredentials)))
	if err != nil {
		log.Fatalf("Error creating android management service: %v", err)
	}

	switch *command {
	case "enterprises.delete":
		enterprisesDelete(mgmt, *enterpriseID)
	case "enterprises.list":
		enterprisesList(mgmt)
	case "policies.list":
		policiesList(mgmt, *enterpriseID)
	case "devices.list":
		devicesList(mgmt, *enterpriseID)
	default:
		log.Fatalf("Unknown command: %s", *command)
	}

}

func enterprisesDelete(mgmt *androidmanagement.Service, enterpriseID string) {
	if enterpriseID == "" {
		log.Fatalf("enterprise_id must be set")
	}
	_, err := mgmt.Enterprises.Delete("enterprises/" + enterpriseID).Do()
	if err != nil {
		log.Fatalf("Error deleting enterprise: %v", err)
	}
}

func enterprisesList(mgmt *androidmanagement.Service) {
	enterprises, err := mgmt.Enterprises.List().ProjectId(androidProjectID).Do()
	if err != nil {
		log.Fatalf("Error listing enterprises: %v", err)
	}
	if len(enterprises.Enterprises) == 0 {
		log.Printf("No enterprises found")
		return
	}
	for _, enterprise := range enterprises.Enterprises {
		log.Printf("Enterprise: %+v", *enterprise)
	}
}

func policiesList(mgmt *androidmanagement.Service, enterpriseID string) {
	if enterpriseID == "" {
		log.Fatalf("enterprise_id must be set")
	}
	result, err := mgmt.Enterprises.Policies.List("enterprises/" + enterpriseID).Do()
	if err != nil {
		log.Fatalf("Error listing policies: %v", err)
	}
	if len(result.Policies) == 0 {
		log.Printf("No policies found")
		return
	}
	for _, policy := range result.Policies {
		log.Printf("Policy: %+v", *policy)
	}
}

func devicesList(mgmt *androidmanagement.Service, enterpriseID string) {
	if enterpriseID == "" {
		log.Fatalf("enterprise_id must be set")
	}
	result, err := mgmt.Enterprises.Devices.List("enterprises/" + enterpriseID).Do()
	if err != nil {
		log.Fatalf("Error listing devices: %v", err)
	}
	if len(result.Devices) == 0 {
		log.Printf("No policies found")
		return
	}
	for _, device := range result.Devices {
		data, err := json.MarshalIndent(device, "", "  ")
		if err != nil {
			log.Fatalf("Error marshalling device: %v", err)
		}
		log.Println(string(data))
	}
}
