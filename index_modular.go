package main

import (
//    "bytes"
    "golang.org/x/crypto/ssh"
    "fmt"
    "sync"
    "time"
    "./utils"
    "./modules/hadoop_standalone"
)

// Main
func main(){

	// intro
	utils.Intro()

	// analyze
	analyze()
}


// analyze
func analyze() {

	hostInfo := utils.GetHostInfo()

	// thread initialization
	var wg sync.WaitGroup
	wg.Add(len(hostInfo))

	for _,h := range hostInfo {
		
		fmt.Println("Creating thread for "+h.Ip)

		// Dial your ssh server.
		go func(ip string) {
			fmt.Println("Spawing thread for "+ip)
			conn, err := ssh.Dial("tcp", ip+":22",utils.RemoteCredentials())
			
			if err != nil {
            			fmt.Println(err)
   			}


			defer wg.Done();

			utils.Start(ip)
			hadoop_standalone.RunTests(ip, conn)
			utils.End(ip)
			}(h.Ip)

	}

	time.Sleep(3000)

	wg.Wait()
	
}
