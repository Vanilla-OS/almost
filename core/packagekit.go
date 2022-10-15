package core

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

// org.freedesktop.PackageKit.Offline.UpdatePrepared
// org.freedesktop.PackageKit.Offline.UpgradePrepared

const (
	DBUS_NAME              = "org.freedesktop.PackageKit"
	DBUS_PATH              = "/org/freedesktop/PackageKit"
	DBUS_OFFLINE_INTERFACE = "org.freedesktop.PackageKit.Offline"
)

func dbusConn() *dbus.Conn {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		panic(err)
	}

	return conn
}

func PackageKitUpdatePrepared() bool {
	conn := dbusConn()
	defer conn.Close()

	obj := conn.Object(DBUS_NAME, DBUS_PATH)
	val, err := obj.GetProperty(DBUS_OFFLINE_INTERFACE + "." + "UpdatePrepared")
	if err != nil {
		fmt.Println("error:", err)
		return false // shouldn't never happen really but who knows
	}
	return val.Value().(bool)
}

func PackageKitUpgradePrepared() bool {
	conn := dbusConn()
	defer conn.Close()

	obj := conn.Object(DBUS_NAME, DBUS_PATH)
	val, err := obj.GetProperty(DBUS_OFFLINE_INTERFACE + "." + "UpgradePrepared")
	if err != nil {
		fmt.Println("error:", err)
		return false // ..again, should never happen
	}

	return val.Value().(bool)
}
