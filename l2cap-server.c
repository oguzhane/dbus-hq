#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <bluetooth/bluetooth.h>
#include <bluetooth/l2cap.h>
#include <bluetooth/hidp.h>
#include "l2cap_socket.h"


int l2cap_listen(const bdaddr_t *bdaddr, unsigned short psm, int lm, int backlog)
{
	struct sockaddr_l2 addr;
	struct l2cap_options opts;
	int sk;

	if ((sk = socket(PF_BLUETOOTH, SOCK_SEQPACKET, BTPROTO_L2CAP)) < 0) {
		perror ("Cannot create a L2CAP server socket");
		return -1;
	}

	memset(&addr, 0, sizeof(addr));
	addr.l2_family = AF_BLUETOOTH;
	bacpy(&addr.l2_bdaddr, bdaddr);
	addr.l2_psm = htobs(psm);

	if (bind(sk, (struct sockaddr *) &addr, sizeof(addr)) < 0) {
		perror ("Cannot bind a L2CAP server socket");
		goto fail;
	}

	if (setsockopt(sk, SOL_L2CAP, L2CAP_LM, &lm, sizeof(lm)) < 0) {
		perror ("Cannot set socket options");
		goto fail;
	}

	memset(&opts, 0, sizeof(opts));
	opts.imtu = HIDP_DEFAULT_MTU;
	opts.omtu = HIDP_DEFAULT_MTU;
	opts.flush_to = 0xffff;

	if (setsockopt(sk, SOL_L2CAP, L2CAP_OPTIONS, &opts, sizeof(opts)) < 0) {
		perror ("Cannot set L2CAP socket options");
		goto fail;
	}

	if (listen(sk, backlog) < 0) {
		perror ("Cannot listen to a L2CAP server socket");
		goto fail;
	}

	return sk;
fail:
	close(sk);
	return -1;
}

int l2cap_connect(bdaddr_t *src, bdaddr_t *dst, unsigned short psm)
{
	struct sockaddr_l2 addr;
	struct l2cap_options opts;
	int sk;

	if ((sk = socket(PF_BLUETOOTH, SOCK_SEQPACKET, BTPROTO_L2CAP)) < 0) {
		perror ("Cannot create a L2CAP client socket");
		return -1;
	}

	memset(&addr, 0, sizeof(addr));
	addr.l2_family  = AF_BLUETOOTH;
	bacpy(&addr.l2_bdaddr, src);

	if (bind(sk, (struct sockaddr *) &addr, sizeof(addr)) < 0) {
		perror ("Cannot bind a L2CAP client socket");
		close(sk);
		return -1;
	}

	memset(&opts, 0, sizeof(opts));
	opts.imtu = HIDP_DEFAULT_MTU;
	opts.omtu = HIDP_DEFAULT_MTU;
	opts.flush_to = 0xffff;

	if (setsockopt(sk, SOL_L2CAP, L2CAP_OPTIONS, &opts, sizeof(opts)) < 0) {
		perror ("Cannot set L2CAP socket options");
		goto fail;
	}

	memset(&addr, 0, sizeof(addr));
	addr.l2_family  = AF_BLUETOOTH;
	bacpy(&addr.l2_bdaddr, dst);
	addr.l2_psm = htobs(psm);

	if (connect(sk, (struct sockaddr *) &addr, sizeof(addr)) < 0) {
		perror ("Cannot connect to a L2CAP client socket");
		goto fail;
	}

	return sk;
fail:
	close(sk);
	return -1;
}

int l2cap_accept(int sk, bdaddr_t *bdaddr)
{
	struct sockaddr_l2 addr;
	socklen_t addrlen;
	int nsk;

	memset(&addr, 0, sizeof(addr));
	addrlen = sizeof(addr);

	if ((nsk = accept(sk, (struct sockaddr *) &addr, &addrlen)) < 0) {
		perror ("Cannot accept a L2CAP connection");
		return -1;
	}

	if (bdaddr)
		bacpy(bdaddr, &addr.l2_bdaddr);

	return nsk;
}

int main(int argc, char **argv)
{
	int sintr, scontrol;	/* server sockets */

		/* Prepare the server sockets, in case a client will connect. */
	sintr = l2cap_listen(L2CAP_SERVER_BLUETOOTH_ADDR, INTR_PORT_NUM, 0, 1);
	if (sintr == -1) {
		// close (input);
		return 0;
	}

	scontrol = l2cap_listen(L2CAP_SERVER_BLUETOOTH_ADDR, CTRL_PORT_NUM, 0, 1);
	if (scontrol == -1) {
		// close (input);
		close (sintr);
		return 0;
	}

	int c;
	printf("press any key to exit..");
   	c = getchar();
	return 0;
}

void do_connect(){
		struct sockaddr_l2 scont_addr = { 0 }, scont_rem_addr = { 0 };
	int scont_socket, ccont_socket;
	unsigned int cont_opt = sizeof(scont_rem_addr);

	scont_socket = socket(AF_BLUETOOTH, SOCK_DGRAM, BTPROTO_L2CAP);

	scont_addr.l2_family = AF_BLUETOOTH;						/* Addressing family, always AF_BLUETOOTH */
	bacpy(&scont_addr.l2_bdaddr, L2CAP_SERVER_BLUETOOTH_ADDR);	/* Bluetooth address of local bluetooth adapter */
	scont_addr.l2_psm = htobs(CTRL_PORT_NUM);					/* port number of local bluetooth adapter */

	if(bind(scont_socket, (struct sockaddr *)&scont_addr, sizeof(scont_addr)) < 0) {
		perror("failed to bind scont");
		exit(1);
	}

	//---------
	
	struct sockaddr_l2 sintr_addr = { 0 }, sintr_rem_addr = { 0 };
	int sintr_socket, cintr_socket;
	unsigned int intr_opt = sizeof(sintr_rem_addr);

	sintr_socket = socket(AF_BLUETOOTH, SOCK_DGRAM, BTPROTO_L2CAP);

	sintr_addr.l2_family = AF_BLUETOOTH;						/* Addressing family, always AF_BLUETOOTH */
	bacpy(&sintr_addr.l2_bdaddr, L2CAP_SERVER_BLUETOOTH_ADDR);					/* Bluetooth address of local bluetooth adapter */
	sintr_addr.l2_psm = htobs(INTR_PORT_NUM);			/* port number of local bluetooth adapter */

	if(bind(sintr_socket, (struct sockaddr *)&sintr_addr, sizeof(sintr_addr)) < 0) {
		perror("failed to bind sintr");
		exit(1);
	}




	/***/
	printf("listening\n");
	listen(scont_socket, 1);
	listen(sintr_socket, 1);
	printf("listening2\n");

	ccont_socket = accept(scont_socket, (struct sockaddr *)&scont_rem_addr, &cont_opt);	/* return new socket for connection with a client */
	printf("accept1\n");
	
	cintr_socket = accept(sintr_socket, (struct sockaddr *)&sintr_rem_addr, &intr_opt);	/* return new socket for connection with a client */
	printf("accept2\n");

	/**/
	int c;
	printf("press any key to exit..");
   	c = getchar();

	/* close connection */
	close(scont_socket);
	close(ccont_socket);
}