package local

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type network []struct {
	Host     string
	HostPort string
}

// Spurious contains our local network for faked AWS resources
type Spurious struct {
	Sqs    network `json:"spurious-sqs"`
	S3     network `json:"spurious-s3"`
	Dynamo network `json:"spurious-dynamo"`
}

var spur Spurious

// CheckFlags verifies required network details are provided and if not resort to Spurious
func CheckFlags(production, sqs, s3, dynamo *string) {
	if *production == "true" {
		return
	}

	if *sqs == "" || *s3 == "" || *dynamo == "" {
		cmdOut, err := getSpuriousNetworkDetails()
		if err != nil {
			fmt.Fprintln(os.Stderr, "There was an error running 'spurious ports --json' command: ", err)
			os.Exit(1)
		}

		json.Unmarshal(cmdOut, &spur)
		setMissingSpuriousNetworkDetails(sqs, s3, dynamo)
	}
}

func getSpuriousNetworkDetails() ([]byte, error) {
	var cmdOut []byte
	var err error

	cmdName := "spurious"
	cmdArgs := []string{"ports", "--json"}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		return nil, err
	}

	return cmdOut, nil
}

func setMissingSpuriousNetworkDetails(sqs, s3, dynamo *string) {
	if *sqs == "" {
		*sqs = spur.Sqs[0].Host + ":" + spur.Sqs[0].HostPort
	}

	if *s3 == "" {
		*s3 = spur.S3[0].Host + ":" + spur.S3[0].HostPort
	}

	if *dynamo == "" {
		*dynamo = spur.Dynamo[0].Host + ":" + spur.Dynamo[0].HostPort
	}
}
