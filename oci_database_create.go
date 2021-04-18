package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/database"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"log"
	"os"
	"time"
)

type Options struct {
	DbName     string `short:"d" long:"db-name" description:"Database name" required:"true"`
	DbWorkload string `short:"l" long:"db-workload" description:"Database workload : OLTP or DSS" required:"false" default:"OLTP"`
	DbHomeId   string `short:"o" long:"db-home-id" description:"Database home OCID" required:"true"`
	//VMClusterId string `short:"c" long:"vm-cluster-id" description:"VM cluster OCID" required:"true"`
	WaitForState  string `short:"w" long:"wait-for-state" description:"Wait for state : AVAILABLE, TERMINATED, etc." required:"false"`
	AdminPassword string `short:"p" long:"admin-password" description:"Database password" required:"true"`
	// DB version not implemented because DB Home OCID is used
	//DBVersion           string `short:"v" long:"db-version" description:"Database version : 19.0.0.0, 11.2.0.4, 12.2.0.1" required:"true"`
	DBUniqueName        string `short:"u" long:"db-unique-name" description:"Database Unique Name" required:"false"`
	CharSet             string `short:"s" long:"character-set" description:"Character Set" required:"false" default:"AL32UTF8"`
	NCharSet            string `short:"n" long:"national-character-set" description:"National Character Set" required:"false" default:"AL16UTF16"`
	PDBName             string `short:"b" long:"pdb-name" description:"PDB name" required:"false"`
	TDEWalletPassword   string `short:"x" long:"tde-wallet-password" description:"TDE Wallet Password" required:"false"`
	WaitIntervalSeconds int    `short:"i" long:"wait-interval-seconds" description:"Wait Interval Seconds" default:"30" required:"false"`
	MaxWaitSeconds      int    `short:"m" long:"max-wait-seconds" description:"Max Wait Seconds" default:"3600" required:"false"`
	DryRun              bool   `short:"t" long:"dry-run" description:"Display request only" required:"false"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {

	//parse flags
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	if options.DryRun {
		optionsJSON, err := json.MarshalIndent(options, "", "  ")
		helpers.FatalIfError(err)

		log.Printf("This will create a database with the following options : \n%s", optionsJSON)
		os.Exit(0)
	}

	//Run database creation
	createDBResp, err := createDB(options)
	helpers.FatalIfError(err)

	dbJSON, err := dbCreateRespToJson(createDBResp)
	helpers.FatalIfError(err)
	fmt.Printf("%s", dbJSON)

	if options.WaitForState != "" {
		//Wait for status
		stateReached, err := waitForStatus(*createDBResp.Id, options)
		helpers.FatalIfError(err)

		if stateReached {
			log.Printf("Database State is :  %s\n", options.WaitForState)
		}

	}

}

func dbCreateRespToJson(dbCreateResponse database.CreateDatabaseResponse) ([]byte, error) {
	return json.MarshalIndent(dbCreateResponse.Database, "", "  ")
}

func waitForStatus(databaseId string, dbCreateOptions Options) (bool, error) {
	timeout := time.After(time.Duration(dbCreateOptions.MaxWaitSeconds) * time.Second)
	tick := time.Tick(time.Duration(dbCreateOptions.WaitIntervalSeconds) * time.Second)
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		case <-timeout:
			return false, errors.New("timed out before reaching state " + dbCreateOptions.WaitForState)
		case <-tick:
			currentDBLifecycleStatus, err := getLifeCycleStatus(databaseId)
			helpers.FatalIfError(err)
			log.Printf("Current database status is : %s", currentDBLifecycleStatus)
			if string(currentDBLifecycleStatus) == dbCreateOptions.WaitForState {
				return true, nil
			}
		}
	}

}

func getLifeCycleStatus(databaseId string) (database.DatabaseLifecycleStateEnum, error) {
	c, clerr := database.NewDatabaseClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	getDatabaseRequest := database.GetDatabaseRequest{
		DatabaseId:      &databaseId,
		OpcRequestId:    nil,
		RequestMetadata: common.RequestMetadata{},
	}

	getDBResp, err := c.GetDatabase(context.Background(), getDatabaseRequest)
	helpers.FatalIfError(err)

	return getDBResp.LifecycleState, err
}

func createDB(dbCreateOptions Options) (database.CreateDatabaseResponse, error) {
	c, clerr := database.NewDatabaseClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	createDatabaseDetails := database.CreateDatabaseDetails{
		DbName:        &dbCreateOptions.DbName,
		AdminPassword: &dbCreateOptions.AdminPassword,
		//DbUniqueName:            &uniqueName,
		//PdbName:                 &dbCreateOptions.PDBName,
		CharacterSet:  &dbCreateOptions.CharSet,
		NcharacterSet: &dbCreateOptions.NCharSet,
		//DbWorkload:              "OLTP",
	}

	if options.DBUniqueName != "" {
		createDatabaseDetails.DbUniqueName = &dbCreateOptions.DBUniqueName
	}

	if options.TDEWalletPassword != "" {
		createDatabaseDetails.TdeWalletPassword = &dbCreateOptions.TDEWalletPassword
	}

	if options.PDBName != "" {
		createDatabaseDetails.PdbName = &dbCreateOptions.PDBName
	}

	if options.DbWorkload == "OLTP" {
		createDatabaseDetails.DbWorkload = "OLTP"
	} else {
		createDatabaseDetails.DbWorkload = "DSS"
	}

	createNewDatabaseDetails := database.CreateNewDatabaseDetails{
		DbHomeId: &dbCreateOptions.DbHomeId,
		Database: &createDatabaseDetails,
		//DbVersion: &dbCreateOptions.DBVersion,
	}

	createDatabaseRequest := database.CreateDatabaseRequest{
		CreateNewDatabaseDetails: createNewDatabaseDetails,
		OpcRetryToken:            nil,
		OpcRequestId:             nil,
		RequestMetadata:          common.RequestMetadata{},
	}

	createDBRest, err := c.CreateDatabase(context.Background(), createDatabaseRequest)
	//helpers.FatalIfError(err)

	return createDBRest, err

}
