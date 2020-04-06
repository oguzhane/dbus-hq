package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const intro = `
<node>
	<interface name="com.github.oguzhane.btkb">
		<method name="Foo">
			<arg direction="out" type="s"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `</node> `

type lowerCaseExport struct{}

func (export lowerCaseExport) Foo() (string, *dbus.Error) {
	return "bar", nil
}

type bluezProfile struct{}

func (bluezProfile) Cancel() *dbus.Error {
	return nil
}

func (bluezProfile) NewConnection(path dbus.ObjectPath, fs dbus.UnixFDIndex, props map[string]dbus.Variant) *dbus.Error {
	return nil
}

func (bluezProfile) Release() *dbus.Error {
	return nil
}

func (bluezProfile) RequestDisconnection(path dbus.ObjectPath) *dbus.Error {
	return nil
}

type btkService struct{}

func (btkService) SendKeys(modifier byte, keys []byte) *dbus.Error {
	return nil
}

// func (f foo) Foo() (string, *dbus.Error) {
// 	fmt.Println(f)
// 	return string(f), nil
// }
func dbusServer() {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	blzp := bluezProfile{}
	btkserv := btkService{}

	conn.Export(blzp, "/bluez/oguzhane/profile", "org.bluez.Profile1")
	conn.Export(btkserv, "/org/oguzhane/btkservice", "org.oguzhane.btkservice")

	n := &introspect.Node{
		Name: "/bluez/oguzhane/profile",
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name:    "org.bluez.Profile",
				Methods: introspect.Methods(blzp),
			},
		},
	}

	nserv := &introspect.Node{
		Name: "/org/oguzhane/btkservice",
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name:    "org.oguzhane.btkservice",
				Methods: introspect.Methods(btkserv),
			},
		},
	}

	conn.Export(introspect.NewIntrospectable(n), "/bluez/oguzhane/profile", "org.freedesktop.DBus.Introspectable")
	conn.Export(introspect.NewIntrospectable(nserv), "/org/oguzhane/btkservice", "org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName("com.github.oguzhane.btkb",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
	fmt.Println("Listening on com.github.oguzhane.btkb / /bluez/oguzhane/profile ...")
	registerProfile(conn)
	fmt.Println("Profile registered")
	select {}
}

func pressToExit() {
	fmt.Println("press button to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

func registerProfile(conn *dbus.Conn) {
	sdp := `<?xml version="1.0" encoding="UTF-8" ?>

	<record>
		<attribute id="0x0001">
			<sequence>
				<uuid value="0x1124" />
			</sequence>
		</attribute>
		<attribute id="0x0004">
			<sequence>
				<sequence>
					<uuid value="0x0100" />
					<uint16 value="0x0011" />
				</sequence>
				<sequence>
					<uuid value="0x0011" />
				</sequence>
			</sequence>
		</attribute>
		<attribute id="0x0005">
			<sequence>
				<uuid value="0x1002" />
			</sequence>
		</attribute>
		<attribute id="0x0006">
			<sequence>
				<uint16 value="0x656e" />
				<uint16 value="0x006a" />
				<uint16 value="0x0100" />
			</sequence>
		</attribute>
		<attribute id="0x0009">
			<sequence>
				<sequence>
					<uuid value="0x1124" />
					<uint16 value="0x0100" />
				</sequence>
			</sequence>
		</attribute>
		<attribute id="0x000d">
			<sequence>
				<sequence>
					<sequence>
						<uuid value="0x0100" />
						<uint16 value="0x0013" />
					</sequence>
					<sequence>
						<uuid value="0x0011" />
					</sequence>
				</sequence>
			</sequence>
		</attribute>
		<attribute id="0x0100">
			<text value="Raspberry Pi Virtual Keyboard" />
		</attribute>
		<attribute id="0x0101">
			<text value="USB > BT Keyboard" />
		</attribute>
		<attribute id="0x0102">
			<text value="Raspberry Pi" />
		</attribute>
		<attribute id="0x0200">
			<uint16 value="0x0100" />
		</attribute>
		<attribute id="0x0201">
			<uint16 value="0x0111" />
		</attribute>
		<attribute id="0x0202">
			<uint8 value="0x40" />
		</attribute>
		<attribute id="0x0203">
			<uint8 value="0x00" />
		</attribute>
		<attribute id="0x0204">
			<boolean value="false" />
		</attribute>
		<attribute id="0x0205">
			<boolean value="false" />
		</attribute>
		<attribute id="0x0206">
			<sequence>
				<sequence>
					<uint8 value="0x22" />
					<text encoding="hex" value="05010906a101850175019508050719e029e715002501810295017508810395057501050819012905910295017503910395067508150026ff000507190029ff8100c0050c0901a1018503150025017501950b0a23020a21020ab10109b809b609cd09b509e209ea09e9093081029501750d8103c0" />
				</sequence>
			</sequence>
		</attribute>
		<attribute id="0x0207">
			<sequence>
				<sequence>
					<uint16 value="0x0409" />
					<uint16 value="0x0100" />
				</sequence>
			</sequence>
		</attribute>
		<attribute id="0x020b">
			<uint16 value="0x0100" />
		</attribute>
		<attribute id="0x020c">
			<uint16 value="0x0c80" />
		</attribute>
		<attribute id="0x020d">
			<boolean value="true" />
		</attribute>
		<attribute id="0x020e">
			<boolean value="false" />
		</attribute>
		<attribute id="0x020f">
			<uint16 value="0x0640" />
		</attribute>
		<attribute id="0x0210">
			<uint16 value="0x0320" />
		</attribute>
	</record>
	`
	// profile := "/bluez/oguzhane/profile"
	var profile dbus.ObjectPath = "/bluez/oguzhane/profile"
	uuid := "00001124-0000-1000-8000-00805f9b34fb"
	opts := map[string]dbus.Variant{
		"ServiceRecord":         dbus.MakeVariant(sdp),
		"Role":                  dbus.MakeVariant("server"),
		"RequireAuthentication": dbus.MakeVariant(false),
		"RequireAuthorization":  dbus.MakeVariant(false),
	}
	err := conn.Object("org.bluez", "/org/bluez").Call("org.bluez.ProfileManager1.RegisterProfile", 0, profile, uuid, opts).Store()
	fatalErr(err)
}
func main() {
	dbusServer()
	pressToExit()
	//clean up connection on exit
	/*defer api.Exit()

	app, err := service.NewApp("hci0")
	fatalErr(err)
	defer app.Close()
	app.SetName("go_bluetooth")

	service1, err := app.NewService()
	fatalErr(err)
	app.AddService(service1)
	err = service1.Expose()

	fatalErr(err)
	err = app.Run()
	fatalErr(err)


	*/
	// // c, _ := serv.NewChar()
	// // c.OnWrite(func(c *Char, value []byte) ([]byte, error) {

	// // })
	// devc, _ := device.NewDevice("hci0", "3C:F8:62:EA:82:87")
	// devc.SetProperty("class", uint32(0x002540))
	// devc.SetProperty("piscan", true)
	// profile.NewProfile1("/bluez/ogz/btkb_profile", dbus.ObjectPath("org.bluez.Profile1"))
	// manager, _ := profile.NewProfileManager1()
	// opts := map[string]interface{}{
	// 	"ServiceRecord": `<?xml version="1.0" encoding="UTF-8" ?>

	// 	<record>
	// 		<attribute id="0x0001">
	// 			<sequence>
	// 				<uuid value="0x1124" />
	// 			</sequence>
	// 		</attribute>
	// 		<attribute id="0x0004">
	// 			<sequence>
	// 				<sequence>
	// 					<uuid value="0x0100" />
	// 					<uint16 value="0x0011" />
	// 				</sequence>
	// 				<sequence>
	// 					<uuid value="0x0011" />
	// 				</sequence>
	// 			</sequence>
	// 		</attribute>
	// 		<attribute id="0x0005">
	// 			<sequence>
	// 				<uuid value="0x1002" />
	// 			</sequence>
	// 		</attribute>
	// 		<attribute id="0x0006">
	// 			<sequence>
	// 				<uint16 value="0x656e" />
	// 				<uint16 value="0x006a" />
	// 				<uint16 value="0x0100" />
	// 			</sequence>
	// 		</attribute>
	// 		<attribute id="0x0009">
	// 			<sequence>
	// 				<sequence>
	// 					<uuid value="0x1124" />
	// 					<uint16 value="0x0100" />
	// 				</sequence>
	// 			</sequence>
	// 		</attribute>
	// 		<attribute id="0x000d">
	// 			<sequence>
	// 				<sequence>
	// 					<sequence>
	// 						<uuid value="0x0100" />
	// 						<uint16 value="0x0013" />
	// 					</sequence>
	// 					<sequence>
	// 						<uuid value="0x0011" />
	// 					</sequence>
	// 				</sequence>
	// 			</sequence>
	// 		</attribute>
	// 		<attribute id="0x0100">
	// 			<text value="Raspberry Pi Virtual Keyboard" />
	// 		</attribute>
	// 		<attribute id="0x0101">
	// 			<text value="USB > BT Keyboard" />
	// 		</attribute>
	// 		<attribute id="0x0102">
	// 			<text value="Raspberry Pi" />
	// 		</attribute>
	// 		<attribute id="0x0200">
	// 			<uint16 value="0x0100" />
	// 		</attribute>
	// 		<attribute id="0x0201">
	// 			<uint16 value="0x0111" />
	// 		</attribute>
	// 		<attribute id="0x0202">
	// 			<uint8 value="0x40" />
	// 		</attribute>
	// 		<attribute id="0x0203">
	// 			<uint8 value="0x00" />
	// 		</attribute>
	// 		<attribute id="0x0204">
	// 			<boolean value="false" />
	// 		</attribute>
	// 		<attribute id="0x0205">
	// 			<boolean value="false" />
	// 		</attribute>
	// 		<attribute id="0x0206">
	// 			<sequence>
	// 				<sequence>
	// 					<uint8 value="0x22" />
	// 					<text encoding="hex" value="05010906a101850175019508050719e029e715002501810295017508810395057501050819012905910295017503910395067508150026ff000507190029ff8100c0050c0901a1018503150025017501950b0a23020a21020ab10109b809b609cd09b509e209ea09e9093081029501750d8103c0" />
	// 				</sequence>
	// 			</sequence>
	// 		</attribute>
	// 		<attribute id="0x0207">
	// 			<sequence>
	// 				<sequence>
	// 					<uint16 value="0x0409" />
	// 					<uint16 value="0x0100" />
	// 				</sequence>
	// 			</sequence>
	// 		</attribute>
	// 		<attribute id="0x020b">
	// 			<uint16 value="0x0100" />
	// 		</attribute>
	// 		<attribute id="0x020c">
	// 			<uint16 value="0x0c80" />
	// 		</attribute>
	// 		<attribute id="0x020d">
	// 			<boolean value="true" />
	// 		</attribute>
	// 		<attribute id="0x020e">
	// 			<boolean value="false" />
	// 		</attribute>
	// 		<attribute id="0x020f">
	// 			<uint16 value="0x0640" />
	// 		</attribute>
	// 		<attribute id="0x0210">
	// 			<uint16 value="0x0320" />
	// 		</attribute>
	// 	</record>
	// 	`,
	// 	"Role":                  "server",
	// 	"RequireAuthentication": false,
	// 	"RequireAuthorization":  false,
	// }
	// manager.RegisterProfile("/bluez/ogz/btkb_profile", "00001124-0000-1000-8000-00805f9b34fb", opts)
	// devc.ConnectProfile("00001124-0000-1000-8000-00805f9b34fb")
	// devc.Properties.Class
	// devc.ConnectProfile()
}
