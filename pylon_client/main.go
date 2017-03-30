package main

import (
	"log"
	"fmt"
	servicepb "github.com/axsy-dev/pylon/pkg/service/pb"
	domainpb "github.com/axsy-dev/pylon/pkg/postgres/pb"
	"context"
	"google.golang.org/grpc"
	"io/ioutil"
	"path/filepath"
)

func main() {
	server := "localhost:8080"
	domainName := "davesdomain"
	modName := "davemod"
	modVersion := 1
	var client servicepb.PylonClient
	if con, err := grpc.Dial(server, grpc.WithInsecure()); err != nil {
		log.Fatal(fmt.Sprintf(`Could not connect to server "%s": %s`, server, err))
	} else {
		client = servicepb.NewPylonClient(con)
		client.EnsureDomain(context.Background(), &domainpb.DomainID{domainName})
	}

	payload := loadModulePayload(modName, modVersion)
	log.Print(string(payload))

	if res, err:=client.EnsureModule(context.Background(), &domainpb.ModuleMigration{
		Domain:  domainName,
		Module:  modName,
		Version: int32(modVersion),
		Payload: string(payload),
	}) ; err != nil {
		log.Panic(err)
	} else {

		log.Print(res.OldVersion)
	}


}

func loadModulePayload(name string,  version int) []byte {
	fileName := filepath.Join("modules", fmt.Sprintf("%s_%d.js", name, version))
	res, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicf("could not read file fixture file: %s", fileName)
	}
	return res
}
