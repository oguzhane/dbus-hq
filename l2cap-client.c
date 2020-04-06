#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <bluetooth/bluetooth.h>
#include <bluetooth/l2cap.h>
#include "l2cap_socket.h"
 
int main(int argc, char **argv)
{
	struct sockaddr_l2 addr = { 0 };
	int sock;
	const char *sample_text = "L2CAP Simple Example Done";

	printf("Start Bluetooth L2CAP client, server addr %s\n", L2CAP_SERVER_BLUETOOTH_ADDR);
	
	/* allocate a socket */
	sock = socket(AF_BLUETOOTH, SOCK_SEQPACKET, BTPROTO_L2CAP);

	/* set the outgoing connection parameters, server's address and port number */
	addr.l2_family = AF_BLUETOOTH;								/* Addressing family, always AF_BLUETOOTH */
	addr.l2_psm = htobs(L2CAP_SERVER_PORT_NUM);					/* server's port number */
	str2ba(L2CAP_SERVER_BLUETOOTH_ADDR, &addr.l2_bdaddr);		/* server's Bluetooth Address */

	/* connect to server */
	if(connect(sock, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
		perror("failed to connect");
		exit(1);
	}

	/* send a message */
	printf("connected...\n");
	send(sock, sample_text, strlen(sample_text), 0);
	
	close(sock);
	return 0;
}