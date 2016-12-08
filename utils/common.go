package utils
// Common Utilities
import (
    "bytes"
    "golang.org/x/crypto/ssh"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    //"sync"
    //"time"
    "../models"
    "strings"
)

func GetGanpatiConfig() models.GanpatiConfig{

        raw, err := ioutil.ReadFile("./config/ganpati.json")
        if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
        }
        var c models.GanpatiConfig
        json.Unmarshal(raw, &c)
        return c
}


// get host information
func GetHostInfo() []models.HostInfo {

        raw, err := ioutil.ReadFile("./config/cluster.json")
        if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
        }
        var c []models.HostInfo
        json.Unmarshal(raw, &c)
        return c
}


func RemoteCredentials() *ssh.ClientConfig {

        config := &ssh.ClientConfig{
                User: "bigdata",
                Auth: []ssh.AuthMethod{
                ssh.Password("bigdata"),
                },
        }
        return config

}

func Start(host string) {

        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = ">>>START"
        if _, err = f.WriteString(info+"\n"); err != nil {
               panic(err)
        }

}

func End(host string) {

        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = ">>>END"
        if _, err = f.WriteString(info+"\n"); err != nil {
               panic(err)
        }

}

func Services( host string, session *ssh.Session, stdoutBuf bytes.Buffer ) {

        session.Stdout = &stdoutBuf
        session.Run("jps | awk '{print $2}'")

        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = "SERVICES|"+stdoutBuf.String()
        if _, err = f.WriteString(strings.Trim(strings.Replace(info,"\n",",",-1),",")+"\n"); err != nil {
                panic(err)
        }

}

func HomeVariables( host string, session *ssh.Session, stdoutBuf bytes.Buffer ) {
        session.Stdout = &stdoutBuf
        session.Run("echo $HADOOP_HOME")

        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = "HADOOP_HOME|"+stdoutBuf.String()

        if _, err = f.WriteString(strings.Trim(strings.Replace(info,"\n",",",-1),",")+"\n"); err != nil {
                panic(err)
        }

}

func Path( host string, session *ssh.Session, stdoutBuf bytes.Buffer ) {
        session.Stdout = &stdoutBuf
        session.Run("echo $PATH")

        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = "PATH|"+stdoutBuf.String()
        if _, err = f.WriteString(info+"\n"); err != nil {
                panic(err)
        }

}


// intro
func Intro(){
        ganpatiConfig := GetGanpatiConfig()
        fmt.Println("----Welcome----")
        fmt.Println(ganpatiConfig.Name + " : " + ganpatiConfig.Version)
}

// new session
func NewSession(conn *ssh.Client) *ssh.Session{

	session,err := conn.NewSession()

	if err != nil {

		panic(err)
	}
	return session
}
