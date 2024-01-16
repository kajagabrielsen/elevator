#include <stdlib.h>
// the address we are listening for messages on
// we have no choice in IP, so use 0.0.0.0, INADDR_ANY, or leave the IP field empty
// the port should be whatever the sender sends to
// alternate names: sockaddr, resolve(udp)addr, 
// InternetAddress addr;

char addr[] = "0.0.0.0";
int port = 30000;

// a socket that plugs our program to the network. This is the "portal" to the outside world
// alternate names: conn
// UDP is sometimes called SOCK_DGRAM. You will sometimes also find UDPSocket or UDPConn as separate types
//recvSock = new Socket(udp)
int socket = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP);

// bind the address we want to use to the socket
recvSock.bind(addr)


// a buffer where the received network data is stored
byte[1024] buffer  

// an empty address that will be filled with info about who sent the data
InternetAddress fromWho 

loop {
    // clear buffer (or just create a new one)
    
    // receive data on the socket
    // fromWho will be modified by ref here. Or it's a return value. Depends.
    // receive-like functions return the number of bytes received
    // alternate names: read, readFrom
    numBytesReceived = recvSock.receiveFrom(buffer, ref fromWho)
    
    // the buffer just contains a bunch of bytes, so you may have to explicitly convert it to a string
    
    // optional: filter out messages from ourselves
    if(fromWho.IP != localIP){
        // do stuff with buffer
    }
}