package main

import (
    "fmt"
    "encoding/xml"
    "io/ioutil"
    "os"
    "golang.org/x/crypto/ssh"
)
//  ----------------------
type Stats struct {
    Size     string  `xml:"size,attr"`
    Consumer string  `xml:"consumerCount,attr"`
    Enqueue  string  `xml:"enqueueCount,attr"`
    Dequeue  string  `xml:"dequeueCount,attr"`
}
type Queue struct {
    Name    string   `xml:"name,attr"`
    Stat    Stats    `xml:"stats"`
}

type XMLAmq struct {
    Res    []Queue  `xml:"queue"`
}
//  ----------------------
func sshInit( )( conf *ssh.ClientConfig, err error ){
    key, err := ioutil.ReadFile( os.Getenv("HOME") + "/.ssh/ant_rsa" )
    if err == nil {
        signer, err := ssh.ParsePrivateKey( key )
        if err == nil {
            conf = &ssh.ClientConfig{
                User: "ant",
                Auth: []ssh.AuthMethod{
                    ssh.PublicKeys( signer ),
                },
            }
        }
    }
//  --
return
}
//  ----------------------
func sshGetAmqs( ip string, config *ssh.ClientConfig )( []Queue, error){
    client, err := ssh.Dial( "tcp", ip+":22", config )
    qs := XMLAmq{}
    if err == nil {
        session, err := client.NewSession()
        if err == nil {
            defer session.Close()
            b, err := session.CombinedOutput( `curl -sS http://localhost:8161/admin/xml/queues.jsp` )
            if err == nil {
                err = xml.Unmarshal( b, &qs )
            }
        }
    }
return qs.Res, err
}
//  ----------------------
//  ----------------------
func main() {

    config, err := sshInit()
    if err != nil {
        panic( err )
    }

    qsRes, err := sshGetAmqs("192.168.1.2", config )
    if err != nil {
        panic( err )
    }

    fmt.Println( "192.168.1.2" )
    for _, q := range qsRes {
        fmt.Printf( " %-36s %s - %s - %s\n",q.Name, q.Stat.Size, q.Stat.Consumer, q.Stat.Dequeue )
    }

}
