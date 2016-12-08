// Hadoop Analysis
package hadoop_standalone

import(
	"bytes"
	"golang.org/x/crypto/ssh"
	"strings"
	"os"
	"fmt"
	"../../utils"
)

func RunTests(host string, conn *ssh.Client){

	var stdoutBuf bytes.Buffer
	defer conn.Close()
	// all tests will come here
	fmt.Println("Running Tests...")

	utils.Services(host, utils.NewSession(conn), stdoutBuf)
	utils.HomeVariables(host,utils.NewSession(conn),stdoutBuf)
	utils.Path(host,utils.NewSession(conn),stdoutBuf)
	ClusterID(host,utils.NewSession(conn),stdoutBuf)


}


func ClusterID( host string, session *ssh.Session, stdoutBuf bytes.Buffer ) {
        session.Stdout = &stdoutBuf

	session.Run("cat /home/bigdata/hadoop/{data,name}/current/VERSION | grep clusterID | awk '{split($$0,c,\"=\");print c[2]}'")
        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = "CLUSTERID|"+stdoutBuf.String()
        if _, err = f.WriteString(strings.Trim(strings.Replace(info,"\n",",",-1),",")+"\n"); err != nil {
                panic(err)
        }

}
