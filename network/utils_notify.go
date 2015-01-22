/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package network

import (
	"dbus/org/freedesktop/notifications"
	"pkg.linuxdeepin.com/lib/dbus"
	. "pkg.linuxdeepin.com/lib/gettext"
	"time"
)

const (
	dbusNotifyDest = "org.freedesktop.Notifications"
	dbusNotifyPath = "/org/freedesktop/Notifications"
)

const (
	notifyIconNetworkOffline            = "notification-network-offline"
	notifyIconWiredConnected            = "notification-network-wired-connected"
	notifyIconWiredDisconnected         = "notification-network-wired-disconnected"
	notifyIconWiredLocal                = "notification-network-wired-local"
	notifyIconWiredError                = notifyIconWiredDisconnected
	notifyIconWirelessConnected         = "notification-network-wireless-full"
	notifyIconWirelessDisconnected      = "notification-network-wireless-disconnected"
	notifyIconWirelessDisabled          = "notification-network-wireless-disabled"
	notifyIconWirelessLocal             = "notification-network-wireless-local"
	notifyIconWirelessError             = notifyIconWirelessDisconnected
	notifyIconVpnConnected              = "notification-network-vpn-connected"
	notifyIconVpnDisconnected           = "notification-network-vpn-disconnected"
	notifyIconProxyEnabled              = "notification-network-proxy-enabled"
	notifyIconProxyDisabled             = "notification-network-proxy-disabled"
	notifyIconNetworkConnected          = notifyIconWiredConnected
	notifyIconNetworkDisconnected       = notifyIconWiredDisconnected
	notifyIconMobile2gConnected         = "notification-network-mobile-2g-connected"
	notifyIconMobile2gDisconnected      = "notification-network-mobile-2g-disconnected"
	notifyIconMobile3gConnected         = "notification-network-mobile-3g-connected"
	notifyIconMobile3gDisconnected      = "notification-network-mobile-3g-disconnected"
	notifyIconMobile4gConnected         = "notification-network-mobile-4g-connected"
	notifyIconMobile4gDisconnected      = "notification-network-mobile-4g-disconnected"
	notifyIconMobileUnknownConnected    = "notification-network-mobile-unknown-connected"
	notifyIconMobileUnknownDisconnected = "notification-network-mobile-unknown-disconnected"
)

var (
	notifyEnabled = true
)

func enableNotify() {
	go func() {
		time.Sleep(5 * time.Second)
		notifyEnabled = true
	}()
}
func disableNotify() {
	notifyEnabled = false
}

func notify(icon, summary, body string) {
	notifier, err := notifications.NewNotifier(dbusNotifyDest, dbusNotifyPath)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("notify", icon, summary, body)
	// use goroutine to fix dbus cycle call issue
	go func() {
		notifier.Notify("Network", 0, icon, summary, body, nil, nil, 0)
		notifications.DestroyNotifier(notifier)
	}()
	return
}

func notifyNetworkOffline() {
	notify(notifyIconNetworkOffline, Tr("Disconnected"), Tr("You are now offline."))
}

func notifyAirplanModeEnabled() {
	// TODO: airplane icon missing
	notify(notifyIconNetworkOffline, Tr("Disconnected"), Tr("Airplane mode enabled."))
}

func notifyWiredCableUnplugged() {
	notify(notifyIconWiredError, Tr("Disconnected"), deviceErrorTable[GUESS_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED])
}

func notifyApModeNotSupport() {
	notify(notifyIconWirelessError, Tr("Disconnected"), Tr("Access Point mode is not supported by this device."))
}

func notifyWirelessHardSwitchOff() {
	notify(notifyIconWirelessDisabled, Tr("Network"), Tr("The hardware switch of WLAN Card is off, please switch on as necessary."))
}

func notifyProxyEnabled() {
	notify(notifyIconProxyEnabled, Tr("Network"), Tr("System proxy is set successfully."))
}
func notifyProxyDisabled() {
	notify(notifyIconProxyDisabled, Tr("Network"), Tr("System proxy has been cancelled."))
}

func notifyVpnConnected(id string) {
	notify(notifyIconVpnConnected, Tr("Connected"), id)
}
func notifyVpnDisconnected(id string) {
	notify(notifyIconVpnDisconnected, Tr("Disconnected"), id)
}
func notifyVpnFailed(id string, reason uint32) {
	notify(notifyIconVpnDisconnected, Tr("Disconnected"), vpnErrorTable[reason])
}

func getMobileConnectedNotifyIcon(mobileNetworkType string) (icon string) {
	switch mobileNetworkType {
	case moblieNetworkType4G:
		icon = notifyIconMobile4gConnected
	case moblieNetworkType3G:
		icon = notifyIconMobile3gConnected
	case moblieNetworkType2G:
		icon = notifyIconMobile2gConnected
	case moblieNetworkTypeUnknown:
		icon = notifyIconMobileUnknownConnected
	default:
		icon = notifyIconMobileUnknownConnected
	}
	return
}
func getMobileDisconnectedNotifyIcon(mobileNetworkType string) (icon string) {
	switch mobileNetworkType {
	case moblieNetworkType4G:
		icon = notifyIconMobile4gDisconnected
	case moblieNetworkType3G:
		icon = notifyIconMobile3gDisconnected
	case moblieNetworkType2G:
		icon = notifyIconMobile2gDisconnected
	case moblieNetworkTypeUnknown:
		icon = notifyIconMobileUnknownDisconnected
	default:
		icon = notifyIconMobileUnknownDisconnected
	}
	return
}

func generalGetNotifyConnectedIcon(devType uint32, devPath dbus.ObjectPath) (icon string) {
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		icon = notifyIconWiredConnected
	case NM_DEVICE_TYPE_WIFI:
		icon = notifyIconWirelessConnected
	case NM_DEVICE_TYPE_MODEM:
		var mobileNetworkType string
		dev := manager.getDevice(devPath)
		if dev != nil {
			manager.devicesLock.Lock()
			defer manager.devicesLock.Unlock()
			mobileNetworkType = dev.MobileNetworkType
		}
		icon = getMobileConnectedNotifyIcon(mobileNetworkType)
	default:
		icon = notifyIconNetworkConnected
	}
	return
}
func generalGetNotifyDisconnectedIcon(devType uint32, devPath dbus.ObjectPath) (icon string) {
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		icon = notifyIconWiredDisconnected
	case NM_DEVICE_TYPE_WIFI:
		icon = notifyIconWirelessDisconnected
	case NM_DEVICE_TYPE_MODEM:
		var mobileNetworkType string
		dev := manager.getDevice(devPath)
		if dev != nil {
			manager.devicesLock.Lock()
			mobileNetworkType = dev.MobileNetworkType
			manager.devicesLock.Unlock()
		}
		icon = getMobileDisconnectedNotifyIcon(mobileNetworkType)
	default:
		icon = notifyIconNetworkDisconnected
	}
	return
}

func notifyDeviceRemoved(devPath dbus.ObjectPath) {
	var devType uint32
	dev := manager.getDevice(devPath)
	if dev != nil {
		manager.devicesLock.Lock()
		devType = dev.nmDevType
		manager.devicesLock.Unlock()
	}
	icon := generalGetNotifyDisconnectedIcon(devType, devPath)
	msg := deviceErrorTable[NM_DEVICE_STATE_REASON_REMOVED]
	notify(icon, Tr("Disconnected"), msg)
}
