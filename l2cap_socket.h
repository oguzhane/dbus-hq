#ifndef L2CAP_SOCKET_H
#define L2CAP_SOCKET_H

/* server channel */
#define L2CAP_SERVER_PORT_NUM			0x0011			
#define L2CAP_SERVER_PORT_INTR			0x0013			

#define CTRL_PORT_NUM			0x0011			
#define INTR_PORT_NUM			0x0013	

/* destination address */
#define L2CAP_SERVER_BLUETOOTH_ADDR (&(bdaddr_t) {{0x3c,0xf8,0x62,0xea,0x82,0x87}})	// "3C:F8:62:EA:82:87"	/* it should be modified in your environment */

// #define BDADDR_ANY  (&(bdaddr_t) {{0, 0, 0, 0, 0, 0}})

#endif