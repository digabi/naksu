package mebroutines_test

import (
  "testing"
  "naksu/mebroutines"

  "os"
  "io/ioutil"
)

func TestGetVagrantFileVersionAbitti (t *testing.T) {
  sampleAbittiVagrantFileContent := `
  # Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
  VAGRANTFILE_API_VERSION = "2"

  def calc_mem(system_mem_mb)
    mem_available_to_vm = (system_mem_mb * 0.74).floor
    lower_limit_mb = ((8192-1024)*0.74).floor # subtracting 1GB because of integrated graphics cards
    if mem_available_to_vm < lower_limit_mb
      $stderr.puts "Not enough memory left for virtual machine: #{mem_available_to_vm} MB. Required minimum memory for the computer is 8 GB"
      abort
    end
    mem_available_to_vm
  end

  def get_amount_of_cpus_and_system_ram
    # Give VM (host memory * 0.74) & all but 1 cpu (logical) cores
    host = RbConfig::CONFIG['host_os']

    cpus = nil
    mem = nil

    if host =~ /darwin/
      cpus = ` + "`" + `sysctl -n hw.logicalcpu_max` + "`" + `.to_i
      mem = ` + "`" + `sysctl -n hw.memsize` + "`" + `.to_i / 1024 / 1024
    elsif host =~ /linux/
      cpus = ` + "`" + `lscpu -p | awk -F',' '!/^#/{print $1}'| sort -u | wc -l` + "`" + `.to_i
      mem = ` + "`" + `grep 'MemTotal' /proc/meminfo | sed -e 's/MemTotal://' -e 's/ kB//'` + "`" + `.to_i / 1024
    elsif host =~ /mswin|mingw|cygwin/
      cpus = ` + "`" + `wmic cpu Get NumberOfLogicalProcessors` + "`" + `.split[1].to_i
      mem = ` + "`" + `wmic computersystem Get TotalPhysicalMemory` + "`" + `.split[1].to_i / 1024 / 1024
    end

    if cpus.nil?
      $stderr.puts "Could not determine the amount of cpus"
      abort
    end
    if mem.nil?
      $stderr.puts "Could not determine the amount of system memory"
      abort
    end
    [[cpus - 1, 2].max, calc_mem(mem)]
  end

  def get_nic_type
    if ENV['NIC']
      return ENV['NIC']
    else
      return 'virtio'
    end
  end

  cpus, mem = get_amount_of_cpus_and_system_ram()
  Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
    config.vm.boot_timeout = 300
    config.vm.box = "digabi/ktp-qa"
    config.vm.box_url = "https://s3-eu-west-1.amazonaws.com/static.abitti.fi/usbimg/qa/vagrant/metadata.json"
    config.vm.provider :virtualbox do |vb|
      vb.name = "Virtual KTP v57"
      vb.gui = true
      vb.customize ["modifyvm", :id, "--ioapic", "on"]
      vb.customize ["modifyvm", :id, "--cpus", cpus]
      vb.customize ["modifyvm", :id, "--memory", mem]
      vb.customize ["modifyvm", :id, "--nictype1", get_nic_type()]
      vb.customize ['modifyvm', :id, '--clipboard', 'bidirectional']
      vb.customize ["modifyvm", :id, "--vram", 24]
    end

    config.vm.synced_folder '~/ktp-jako', '/media/usb1', id: 'media_usb1'
    config.vm.synced_folder ".", "/vagrant", disabled: true
    config.vm.network "public_network", :adapter=>1, auto_config: false
  end
`

  sampleMebVagrantFileContent := `
  # Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
  VAGRANTFILE_API_VERSION = "2"

  def get_amount_of_cpus_and_system_ram
    # Give VM (host memory - 1.5 GB) & all but 1 cpu cores
    host = RbConfig::CONFIG['host_os']

    cpus = 3
    mem = 8192

    begin
      if host =~ /darwin/
        cpus = ` + "`" + `sysctl -n hw.physicalcpu_max` + "`" + `.to_i
        mem = ` + "`" + `sysctl -n hw.memsize` + "`" + `.to_i / 1024 / 1024
      elsif host =~ /linux/
        cpus = ` + "`" + `lscpu -p | awk -F',' '!/^#/{print $2}'| sort -u | wc -l` + "`" + `.to_i
        mem = ` + "`" + `grep 'MemTotal' /proc/meminfo | sed -e 's/MemTotal://' -e 's/ kB//'` + "`" + `.to_i / 1024
      elsif host =~ /mswin|mingw|cygwin/
        cpus = ` + "`" + `wmic cpu Get NumberOfCores` + "`" + `.split[1].to_i
        mem = ` + "`" + `wmic computersystem Get TotalPhysicalMemory` + "`" + `.split[1].to_i / 1024 / 1024
      end
    rescue
    end
    [[cpus - 1, 2].max, [mem - 1536, 6144].max]
  end

  cpus, mem = get_amount_of_cpus_and_system_ram()
  Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
    config.vm.boot_timeout = 300
    config.vm.box = "digabi/ktp-k2018-45489"
    config.vm.box_url = "https://s3-eu-west-1.amazonaws.com/static.abitti.fi/usbimg/k2018-45489/vagrant/metadata.json"
    config.vm.provider :virtualbox do |vb|
      vb.name = "Virtual KTP v37"
      vb.gui = true
      vb.customize ["modifyvm", :id, "--ioapic", "on"]
      vb.customize ["modifyvm", :id, "--cpus", cpus]
      vb.customize ["modifyvm", :id, "--memory", mem]
      vb.customize ["modifyvm", :id, "--nictype1", "virtio"]
      vb.customize ['modifyvm', :id, '--clipboard', 'bidirectional']
      vb.customize ["modifyvm", :id, "--vram", 24]
    end

    config.vm.synced_folder '~/ktp-jako', '/media/usb1', id: 'media_usb1'
    config.vm.synced_folder ".", "/vagrant", disabled: true
    config.vm.network "public_network", :adapter=>1, auto_config: false
  end
`

  tables := []struct {
    vagrantfileContent string
    versionString string
  }{
    {sampleAbittiVagrantFileContent, "Abitti server (digabi/ktp-qa 57)"},
    {sampleMebVagrantFileContent, "Matric Exam server (digabi/ktp-k2018-45489 37)"},
  }

  for _, table := range tables {
    tmpfile, err := ioutil.TempFile("", "mebroutines_test_")
    if err != nil {
      t.Errorf("Could not open %s: %v", tmpfile.Name(), err)
    }

    defer os.Remove(tmpfile.Name())

    if _, err := tmpfile.WriteString(table.vagrantfileContent); err != nil {
  		t.Errorf("Could not write to %s: %v", tmpfile.Name(), err)
  	}
  	if err := tmpfile.Close(); err != nil {
      t.Errorf("Could not close %s: %v", tmpfile.Name(), err)
  	}

    versionString := mebroutines.GetVagrantFileVersion(tmpfile.Name())

    if versionString != table.versionString {
      t.Errorf("GetVagrantFileVersion returns wrong version string \"%s\" instead of \"%s\"", versionString, table.versionString)
    }
  }
}

func TestIfIntlCharsInPath (t *testing.T) {
  tables := []struct {
    path string
    result bool
  }{
    {"C:\\Users\\john.doe\\ktp-jako", false},
    {"C:\\Users\\raimo.keski-vääntö\\ktp-jako", true},
    {"C:\\Users\\john doe\\ktp-jako", false},

    {"/home/someuser/ktp-jako", false},
    {"/home/öylätti/ktp-jako", true},
    {"~/ktp-jako", true},
    {"~root/loremipsum", true},

    {"random whatever string", false},
    {"random whatever string with öljyrätti", true},
    {"wtf!", true},
    {"what?", true},
    {"/home/ktp-user/*", true},
  }

  for _, table := range tables {
    isIntl := mebroutines.IfIntlCharsInPath(table.path)
    if isIntl != table.result {
      t.Errorf("IfIntlCharsInPath gives '%t' instead of '%t' for path '%s'", isIntl, table.result, table.path)
    }
  }
}
