package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	memEnv := os.Getenv("MEMBERS")
	if memEnv == "" {
		log.Fatal("Must specify a set of members")
	}

	members := strings.Split(memEnv, ",")

	log.Printf("Starting Mongo server")
	cmd := exec.Command("mongod", "--replSet", "rs0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	defer cmd.Wait()

	if initiate := os.Getenv("INITIATOR"); initiate != "true" {
		return
	}

	if err := pingWait(members); err != nil {
		log.Fatalf("Error ping wait: %s", err)
	}

	// At this point, mongod may or may not have gotten to the point where it can
	// receive instructions to configure the replica set.  Thus, we have to wait
	// until its ready by sleeping.
	time.Sleep(5 * time.Second)
	log.Printf("Setting up replica set with members: %+v", members)
	if err := setupReplicaSet(members); err != nil {
		log.Printf("Failed setup replica set: %s", err)
	}
}

func setupReplicaSet(members []string) error {
	rsInit := `
		rs.initiate(
			{
				_id: "rs0",
				version: 1,
				members: [%s]
			}
		)`

	var memObjs []string
	for i, m := range members {
		memObjs = append(memObjs, fmt.Sprintf(`{_id: %d, host : "%s"}`, i, m))
	}

	rsInit = fmt.Sprintf(rsInit, strings.Join(memObjs, ","))
	cmd := exec.Command("mongo", "--eval", rsInit)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// pingWait blocks until all `addrs` are pingable.
func pingWait(addrs []string) error {
	var err error
	for _, fullAddr := range addrs {
		addr := strings.Split(fullAddr, ":")[0]
		log.Printf("Pinging %s", addr)
		for pinged := false; !pinged; pinged, err = ping(addr) {
			if err != nil {
				return err
			}
			time.Sleep(time.Second)
		}
		log.Printf("Successfully pinged %s", addr)
	}
	return nil
}

func ping(addr string) (bool, error) {
	err := exec.Command("ping", "-c1", addr).Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
