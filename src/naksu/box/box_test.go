package box_test

import (
	"naksu/box"
	"testing"
)

func TestGetEverything(t *testing.T) {
	sampleShowVMInfoOutput1 := `
  name="SERVER1919C v114"
  groups="/"
  ostype="Debian (32-bit)"
  UUID="8c722e19-bd30-4eb3-b36a-95fc4e20c072"
  CfgFile="/opt/matti/VirtualBox_VMs/SERVER1919C v114/SERVER1919C v114.vbox"
  SnapFldr="/opt/matti/VirtualBox_VMs/SERVER1919C v114/Snapshots"
  LogFldr="/opt/matti/VirtualBox_VMs/SERVER1919C v114/Logs"
  hardwareuuid="8c722e19-bd30-4eb3-b36a-95fc4e20c072"
  memory=11827
  pagefusion="off"
  vram=24
  cpuexecutioncap=100
  hpet="off"
  chipset="piix3"
  firmware="EFI"
  cpus=7
  pae="on"
  longmode="on"
  triplefaultreset="off"
  apic="on"
  x2apic="off"
  cpuid-portability-level=0
  bootmenu="messageandmenu"
  boot1="floppy"
  boot2="dvd"
  boot3="disk"
  boot4="none"
  acpi="on"
  ioapic="on"
  biosapic="apic"
  biossystemtimeoffset=0
  rtcuseutc="on"
  hwvirtex="on"
  nestedpaging="on"
  largepages="off"
  vtxvpid="on"
  vtxux="on"
  paravirtprovider="default"
  effparavirtprovider="kvm"
  VMState="poweroff"
  VMStateChangeTime="2019-05-10T11:44:23.874000000"
  monitorcount=1
  accelerate3d="off"
  accelerate2dvideo="off"
  teleporterenabled="off"
  teleporterport=0
  teleporteraddress=""
  teleporterpassword=""
  tracing-enabled="off"
  tracing-allow-vm-access="off"
  tracing-config=""
  autostart-enabled="off"
  autostart-delay=0
  defaultfrontend=""
  storagecontrollername0="SATA Controller"
  storagecontrollertype0="IntelAhci"
  storagecontrollerinstance0="0"
  storagecontrollermaxportcount0="30"
  storagecontrollerportcount0="30"
  storagecontrollerbootable0="on"
  "SATA Controller-0-0"="/opt/matti/VirtualBox_VMs/SERVER1919C v114/box-disk001.vmdk"
  "SATA Controller-ImageUUID-0-0"="ced7cfb7-82cd-4f36-9e83-c933ba0e0220"
  "SATA Controller-1-0"="none"
  "SATA Controller-2-0"="none"
  "SATA Controller-3-0"="none"
  "SATA Controller-4-0"="none"
  "SATA Controller-5-0"="none"
  "SATA Controller-6-0"="none"
  "SATA Controller-7-0"="none"
  "SATA Controller-8-0"="none"
  "SATA Controller-9-0"="none"
  "SATA Controller-10-0"="none"
  "SATA Controller-11-0"="none"
  "SATA Controller-12-0"="none"
  "SATA Controller-13-0"="none"
  "SATA Controller-14-0"="none"
  "SATA Controller-15-0"="none"
  "SATA Controller-16-0"="none"
  "SATA Controller-17-0"="none"
  "SATA Controller-18-0"="none"
  "SATA Controller-19-0"="none"
  "SATA Controller-20-0"="none"
  "SATA Controller-21-0"="none"
  "SATA Controller-22-0"="none"
  "SATA Controller-23-0"="none"
  "SATA Controller-24-0"="none"
  "SATA Controller-25-0"="none"
  "SATA Controller-26-0"="none"
  "SATA Controller-27-0"="none"
  "SATA Controller-28-0"="none"
  "SATA Controller-29-0"="none"
  bridgeadapter1="em1"
  macaddress1="0800271DE972"
  cableconnected1="on"
  nic1="bridged"
  nictype1="virtio"
  nicspeed1="0"
  nic2="none"
  nic3="none"
  nic4="none"
  nic5="none"
  nic6="none"
  nic7="none"
  nic8="none"
  hidpointing="ps2mouse"
  hidkeyboard="ps2kbd"
  uart1="off"
  uart2="off"
  uart3="off"
  uart4="off"
  lpt1="off"
  lpt2="off"
  audio="none"
  audio_in="false"
  audio_out="false"
  clipboard="bidirectional"
  draganddrop="disabled"
  vrde="off"
  usb="off"
  ehci="off"
  xhci="off"
  SharedFolderNameMachineMapping1="media_usb1"
  SharedFolderPathMachineMapping1="/home/matti/ktp-jako"
  videocap="off"
  videocap_audio="off"
  videocapscreens=0
  videocapfile="/opt/matti/VirtualBox_VMs/SERVER1919C v114/SERVER1919C v114.webm"
  videocapres=1024x768
  videocaprate=512
  videocapfps=25
  videocapopts=
  description="digabi/ktp-qa"
  GuestMemoryBalloon=0
  `

	tables := []struct {
		showVMInfoOutput string
		expectedType     string
		expectedVersion  string
		expectedDiskUUID string
	}{
		{sampleShowVMInfoOutput1, "digabi/ktp-qa", "SERVER1919C v114", "ced7cfb7-82cd-4f36-9e83-c933ba0e0220"},
	}

	for _, table := range tables {
		box.SetCacheShowVMInfo(table.showVMInfoOutput)

		observedType := box.GetType()
		if observedType != table.expectedType {
			t.Errorf("GetType() returns '%s' instead of expected '%s'", observedType, table.expectedType)
		}

		observedVersion := box.GetVersion()
		if observedVersion != table.expectedVersion {
			t.Errorf("GetVersion() returns '%s' instead of expected '%s'", observedVersion, table.expectedVersion)
		}

		observedDiskUUID := box.GetDiskUUID()
		if observedDiskUUID != table.expectedDiskUUID {
			t.Errorf("GetDiskUUID() returns '%s' instead of expected '%s'", observedDiskUUID, table.expectedDiskUUID)
		}
	}
}
