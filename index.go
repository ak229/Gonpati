package main

import (
    "bytes"
    "golang.org/x/crypto/ssh"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "sync"
    "time"
    "strings"
)

// ##### Declare all the necessary structs #####
// --- GanpatiConfig ---


type GanpatiConfig struct {

	Name string `json:"name"`
	Version string `json:"version"`

}

// --- hostInfo ---
type HostInfo struct {

	Ip string `json:"ip"`
	Module []string  `json:"module"`

}

//type hostData struct {

	
//}

// ##### End of Declaration #####


// ##### All necessary Getters #####

// get ganpati config
func getGanpatiConfig() GanpatiConfig{

        raw, err := ioutil.ReadFile("./config/ganpati.json")
        if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
        }
        var c GanpatiConfig
        json.Unmarshal(raw, &c)
        return c
}


// get host information
func getHostInfo() []HostInfo {

 	raw, err := ioutil.ReadFile("./config/cluster.json")
        if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
        }
	fmt.Println(string(raw))
	var c []HostInfo
        json.Unmarshal(raw, &c)
        return c
} 


// ##### End of Getters #####


// ##### Utilities #####

func remoteCredentials() *ssh.ClientConfig {

	config := &ssh.ClientConfig{
    		User: "bigdata",
    		Auth: []ssh.AuthMethod{
        	ssh.Password("bigdata"),
    		},
	}	
	return config

}

func start(host string) {

	f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	var info string = ">>>START"
	if _, err = f.WriteString(info+"\n"); err != nil {
               panic(err)
       	}

}


func end(host string) {

        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = ">>>END"
        if _, err = f.WriteString(info+"\n"); err != nil {
               panic(err)
        }

}

func services( host string, session *ssh.Session, stdoutBuf bytes.Buffer ) {

	session.Stdout = &stdoutBuf
 	session.Run("jps | awk '{print $2}'")

	f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)	
 	var info string = "SERVICES|"+stdoutBuf.String()
	if _, err = f.WriteString(strings.Trim(strings.Replace(info,"\n",",",-1),",")+"\n"); err != nil {
        	panic(err)
 	}

}

func homeVariables( host string, session *ssh.Session, stdoutBuf bytes.Buffer ) {
	session.Stdout = &stdoutBuf
	session.Run("echo $HADOOP_HOME")

	f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	var info string = "HADOOP_HOME|"+stdoutBuf.String()

	if _, err = f.WriteString(strings.Trim(strings.Replace(info,"\n",",",-1),",")+"\n"); err != nil {
        	panic(err)
	}
	
}

func path( host string, session *ssh.Session, stdoutBuf bytes.Buffer ) {
        session.Stdout = &stdoutBuf
        session.Run("echo $PATH")

        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = "PATH|"+stdoutBuf.String()
        if _, err = f.WriteString(info+"\n"); err != nil {
                panic(err)
        }

}

func clusterID( host string, session *ssh.Session, stdoutBuf bytes.Buffer ) {
        session.Stdout = &stdoutBuf
        session.Run("cat /home/bigdata/hadoop/{data,name}/current/VERSION | grep clusterID | awk '{split($0,c,\"=\");print c[2]}'")

        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = "CLUSTERID|"+stdoutBuf.String()
        if _, err = f.WriteString(strings.Trim(strings.Replace(info,"\n",",",-1),",")+"\n"); err != nil {
                panic(err)
        }

}
// ##### End of Utilities #####





// Main
func main(){

	// intro
	intro()

	// analyze
	analyze()
}

// intro
func intro(){
	ganpatiConfig := getGanpatiConfig()
	fmt.Println("----Welcome----")
	fmt.Println(ganpatiConfig.Name + " : " + ganpatiConfig.Version)
}


// analyze
func analyze() {

	hostInfo := getHostInfo()

	// thread initialization
	var wg sync.WaitGroup
	wg.Add(len(hostInfo))

	for _,h := range hostInfo {
		
		fmt.Println("Creating thread for "+h.Ip)

		// Dial your ssh server.
		go func(ip string) {
			fmt.Println("Spawing thread for "+ip)
			conn, err := ssh.Dial("tcp", ip+":22",remoteCredentials())
			
			if err != nil {
            			fmt.Println(err)
   			}


			defer wg.Done();

			start(ip)
			hadoop_analyze(ip, conn)
			end(ip)
			}(h.Ip)

	}

	time.Sleep(3000)

	wg.Wait()
	
}

func hadoop_analyze( host string, conn *ssh.Client ) {
	session1, _ := conn.NewSession()
	session2, _ := conn.NewSession()
	session3, _ := conn.NewSession()
	session4, _ := conn.NewSession()
	var stdoutBuf bytes.Buffer

	defer conn.Close()

	services(host, session1, stdoutBuf)
	homeVariables(host,session2,stdoutBuf)
	path(host,session3,stdoutBuf)
	clusterID(host,session4,stdoutBuf)
}
