# naksu

`naksu` is collection of helper scripts for MEB stickless Abitti/matriculation exam servers.
In real life Naksu is Onerva's (the [Abitti model teacher](https://www.abitti.fi/fi/tutustu/)) cat.

The need for some kind of helper scripts appeared to be evident as we followed support requests
and feedback from school IT staff and teachers.

Naksu is currently in production. The executables can be downloaded from [GitHub](https://github.com/digabi/naksu/releases/latest). Download either Windows or Linux version and execute the file in the OS-related zip. After this naksu updates itself when executed.

## Overview of executed operations

### Fresh Install / Update

1.  Make sure you have `vagrant` executable
1.  Make sure Oracle VirtualBox is installed
1.  Create `~/ktp` if it does not exist
1.  Create `~/ktp-jako` if it does not exist
1.  If `~/ktp/Vagrantfile` exists execute `vagrant destroy -f` (we expect that there is existing installation)
1.  Delete `~/ktp/Vagrantfile.bak`
1.  Rename `~/ktp/Vagrantfile` to `~/ktp/Vagrantfile.bak`
1.  Retrieve `http://static.abitti.fi/usbimg/qa/vagrant/Vagrantfile` to `~/ktp/Vagrantfile`
1.  Execute `vagrant box update`
1.  Execute `vagrant box prune`

### Switch Between Abitti and Matriculation Examination Server

Same as install/update procedure but the user is able to select `Vagrantfile`. The
file must be downloaded beforehand by the school principal.

### Start Virtual Server

1.  Make sure you have `vagrant` executable
1.  Make sure Oracle VirtualBox is installed
1.  Execute `vagrant up`

### Make Backup of Server Hard Drive

1. Make sure you have `vagrant` executable
1. Make sure you have `VBoxManage` executable
1. Get VirtualBox VM id of the vagrant default machine from `~/ktp/.vagrant/machines/default/virtualbox/id`
1. `VBoxManage showvminfo -machinereadable {Machine UUID}` -> Get `Disk UUID` from `SATA Controller-ImageUUID-0-0`
1. `VBoxManage clonemedium {Disk UUID} {destination path} --format VMDK`
1. `VBoxManage closemedium {destination path}`

Since the cloned disks can be quite large the user might want to select the media for the save.
Unfortunately, libui SaveFile dialog [does not support folders](https://github.com/andlabs/libui/issues/314).

## Compiling

Compilation is usually done in Docker container. This means that you can compile Naksu in almost any environment
that supports Docker and make.

You can also build Naksu with [Travis CI](https://travis-ci.org). See [`TRAVIS-CI.md`](TRAVIS-CI.md) for more.

### Requirements

- Install make
- Install `dep` using [these instructions](https://golang.github.io/dep/docs/installation.html)
- Install Docker daemon

### Build both versions using Docker

`make docker`

Resulting binaries are saved to `bin/` directory as `naksu` and `naksu.exe`

Build also produces release bundles as `naksu_linux_amd64.zip` and `naksu_windows_amd64.zip`.

### Compiling test version on Mac OS X

Older `libui` version we use to compile production releases has issues on Mac OS X (specifically, starting the application gives error about undo).

Workaround for this is to unpin the production release of `libui` from `Gopgk.toml` *locally* and recompile the binary. The diff for the change required is following:

```
[[constraint]]
  name = “github.com/andlabs/ui”
-  revision = “6c3bda44d3039e3721c06516be3ab9ce9cbd48cc”
+  branch = “master”
+  #revision = “6c3bda44d3039e3721c06516be3ab9ce9cbd48cc"
```

**Do not commit this to git, since the master version seems to cause run-time linking failure on Windows**

After this, update dependencies with

`make update_libs`

and build a Mac OS X binary with

`make darwin-docker`

Mac OS X test version can be started with command `./bin/naksu-darwin
## Compilation details

Preferred way to compile `naksu` is using Docker.

`naksu` can be used in Linux and Windows environments. The compiling is supported
only on Linux. You need at least Go 1.7 to build naksu. In
Debian/Ubuntu environment install `golang-1.9` or `golang-1.10`.

Make sure `go` points to your compiler or set `GO` to point your go binary (in `Makefile`).

### Building Linux version

- Install `libgtk-3-dev` which is required by `libui`.
- Run `make update_libs` to get all required golang dependencies.
- `make linux`

### Cross-Compiling Windows version

Windows version can be cross-compiled using mingw-w64. You need at least 5.0 of
mingw-w64 libs. Build it from source if your pre-packaged version is older:

1. Make sure your pre-packaged mingw-w64 is installed (Debian/Ubuntu: `mingw-w64`).
1. Get mingw-w64 from https://sourceforge.net/projects/mingw-w64/files/mingw-w64/mingw-w64-release/
1. Build:

```
mkdir ~/mingw-w64/current
./configure --prefix=$HOME/mingw-w64/current --host=x86_64-w64-mingw32
make -j4
make install
```

1. By default the `Makefile` expects mingw-w64 at `$HOME/mingw-w64/current`.
   This can be changed by editing `CGO_LDFLAGS="-L$(HOME)/mingw-w64/current/lib"`
   in the `naksu.exe` rule. The path should point to `lib` under your mingw-w64 install path.

After installing mingw-w64:

`make windows`

### Windows Icon

Windows version is built with icon file. Building `src\naksu.syso` is done with
[rsrc](https://github.com/akavel/rsrc).

## Troubleshooting

In case of trouble execute naksu with `-debug` switch. If naksu can't find your `vagrant`/`VBoxManage` executable(s) you can use `VAGRANTPATH` and `VBOXMANAGEPATH` environment variables to set these by hand:

```
VAGRANTPATH=/opt/vagrant/latest/bin/vagrant naksu
VBOXMANAGEPATH=D:\Oracle\VirtualBox\VBoxManage.exe naksu
```

However, please report these problems since we would like to make naksu as easy to use as possible.

## Publishing

 1. Invent new version number. Naksu uses [semantic versioning](https://semver.org/)
 1. Update Changelog (below) and set version number in `src/naksu/naksu.go` (`const version = "X.Y.Z"`)
 1. Push changes, and wait for [Travis CI](https://travis-ci.org/digabi/naksu) to finish build
 1. Go to [releases](https://github.com/digabi/naksu/releases). There should be a new draft release with zipped Linux and Windows binaries attached. Press the Edit button.
 1. Enter the new version number to "Tag version" field (e.g. "v1.6.0" - note the "v")
 1. Fill "Release title" Field
 1. Copy release note markup from `README.md` to "Describe this release" textarea
 1. Click "Publish release"

## License

[MIT](https://opensource.org/licenses/MIT)

## Changelog

### 1.11.0 (18-FEB-2020)

- Naksu detects the presence of hardware virtualization (VT-X), Hyper-V, and Windows Hypervisor Platform. A warning is shown on startup if:
  1. VT-X is not enabled
  1. Hyper-V or Hypervisor is enabled in Windows
- Network status is displayed in the UI. A warning message in red is displayed in the following situations:
  1. the link speed is the selected interface too low (< 1 Gbit/s)
  1. if no interface has been selected and the link speed of the slowest interface is < 1 Gbit/s
  1. a wireless interface has been selected
- Naksu detects and attempts to fix an issue where older VirtualBox versions (6.x < 6.0.8 or 5.x < 5.2.30) end up with duplicate hard disks in the configuration file. In this situation, Naksu should no longer end up in a state where it can't start up.
- Backup fails instantly with a warning if trying to make a backup larger than 4GB on a FAT32 formatted disk.
- Virtualbox Host-Only Ethernet Adapter is hidden from network device selection (should show only physical devices).
- Fix: network device listing shows link speeds correctly as Mbit/s (not MB/s)
- Fix: crash when cancelling file selection

### 1.10.0 (07-JUN-2019)

- Allow user to select physical network interface from Naksu. This allows user to bypass interface selection in terminal.
- Fix `lsblk` output handling so Naksu works on Ubuntu 19.04

### 1.9.0 (16-MAY-2019)

- Get version of the installed VM from VirtualBox instead of Vagrantfile to
  avoid showing wrong version numbers after a failed install/update
- Store settings file naksu.ini to home directory instead of the directory
  where Naksu was executed
- VirtualBox network adapter type (e.g. virtio) can be changed from the Naksu UI

### 1.8.1 (19-MAR-2019)

- Available disk space is shown in backup media selection list
- Added user feedback to server destroy operation
- Show string "no name" for external disks that have empty volume name
- Bugfix: Calculate available disk space correctly
- Bugfix: Correctly handle nil WMI return value if external backup media does not have a volume name

### 1.8.0 (27-FEB-2019)

- Show available Abitti update in UI to prevent administrators running outdated server version
- Notify administrators if they are starting matriculation exam box with live internet connection
- Added debug logging on startup
- Bugfix: Linux did not open File Manager if `xdg-open` was not installed
- Bugfix: Windows was not able to enumerate connected removable devices
- Bugfix: Windows was not able to get free disk space
- Bugfix: Linux miscalculated free disk space

### 1.7.0 (19-DEC-2018)

User interface

- User interface remembers user language selection
- Automatic self-update can be disabled via ini-file or command-line parameter. This setting is meant for centralized installations and should not be enabled otherwise. User will still get notified
  if `naksu` is out of date via warning popup.

Other changes

- Naksu has minimal command-line interface. see `naksu -h`
- Naksu tries to persist some settings to `naksu.ini`. This ini file is stored to same path as `naksu`

### 1.6.1 (14-DEC-2018)

User interface

- Network connection check uses 4 second timeout to prevent network check getting stuck indefinitely
- Start server button is disabled if server has not yet been installed (with message to first install the server)

### 1.6.0 (3-DEC-2018)

User interface

- Added button for removing exams (`vagrant destroy -f`)
- Added button for removing servers (delete `~/.vagrant.d`, `~/.VirtualBox`, `~/ktp` and `~/VirtualBox VMs`)
- Detect installed VM version from `~/ktp/Vagrantfile` instead of Vagrant data files
- Changed the low disk warning level from 5 Gb to 50 Gb
- Added . (period) to the list of accepted characters in the home directory path
- User can cancel confirmation dialogs by closing the window

Other changes

- Major rewrite of code
- Build and test on a Docker environment instead of local workstation
- Added support for compiling and executing Naksu on Darwin (needed by some of the developers)

### 1.5.0 (2-OCT-2018)

User interface

- Hide management buttons behind a checkbox
- Removed Exit button in favour of Open ktp-jako button
- Refuse to install/update without network connection by disabling the buttons

Other changes

- Warn user if home path contains 8-bit characters
- Refuse to import Vagrantfile from ~/ktp/Vagrantfile
- Added desktop as a valid backup target location
- Make sure the selected backup target location is writeable
- Hide 'Invalid state while booting' errors from UI
- Use term "virtual server" instead of "stickless server"
- Improvements which make it easier to build Naksi using a continuous integration server

### 1.4.1 (05-SEP-2018)

- Always create debug file to `~ktp/naksu_lastlog.txt` or `TEMP/naksu_lastlog.txt`
- Give more informative error message when vagrant exits with `--macaddress` error message
- Bugfix: get Linux temporary path from `TempDir()` instead of `TEMP` env var

### 1.4.0 (31-AUG-2018)

- Show current server version
- Show status messages when working with vagrant/VBoxManage
- UI Language support (Finnish, Swedish)

### 1.3.0 (13-AUG-2018)

- Restart with `x-terminal-emulator` if started via file manager (Linux)
- Warn user if there is less than 5 Gb free disk before install, update or backup
- Bugfix: Don't panic if media does not have vendor/model (Linux)
- Get vendor/model strings for removable media directly from WMI without `wmic` (Windows)
- Application icon (Windows)
- Added logging (use `-debug` switch to show debug messages)

### 1.2.0 (09-AUG-2018)

- User can create a backup (clone) from VM hard drive
- Continues execution although binary update fails

### 1.1.0 (05-JUL-2018)

- User can switch between Abitti and Matriculation Exam servers
