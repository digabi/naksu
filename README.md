# naksu

`naksu` is collection of helper scripts for MEB stickless Abitti/matriculation exam servers.
In real life Naksu is Onerva's (the [Abitti model teacher](https://www.abitti.fi/fi/tutustu/)) cat.

The need for some kind of helper scripts appeared to be evident as we followed support requests
and feedback from school IT staff and teachers.

These scripts are currently under planning/proof-of-concept stage. The executables can be downloaded from [GitHub](https://github.com/digabi/naksu/releases/latest). Download either Windows or Linux version and execute the file in the OS-related zip. After this naksu updates itself when executed.

## Plan for the Scripts

### Fresh Install / Update

 1. Make sure you have `vagrant` executable
 1. Make sure Oracle VirtualBox is installed
 1. Create `~/ktp` if it does not exist
 1. Create `~/ktp-jako` if it does not exist
 1. If `~/ktp/Vagrantfile` exists execute `vagrant destroy -f` (we expect that there is existing installation)
 1. Delete `~/ktp/Vagrantfile.bak`
 1. Rename `~/ktp/Vagrantfile` to `~/ktp/Vagrantfile.bak`
 1. Retrieve `http://static.abitti.fi/usbimg/qa/vagrant/Vagrantfile` to `~/ktp/Vagrantfile`
 1. Execute `vagrant box update`
 1. Execute `vagrant box prune`

### Switch Between Abitti and Matriculation Examination Server

 Same as install/update procedure but the user is able to select `Vagrantfile`. The
 file must be downloaded beforehand by the school principal.

### Start Virtual Server

 1. Make sure you have `vagrant` executable
 1. Make sure Oracle VirtualBox is installed
 1. Execute `vagrant up`

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

`naksu` can be used in Linux and Windows environments. The compiling is supported
only on Linux. You need at least Go 1.7 to build naksu. In
Debian/Ubuntu environment install `golang-1.9` or `golang-1.10`.

Make sure `go` points to your compiler or set `GO` to point your go binary (in `Makefile`).

### Requirements

 * Install `libgtk-3-dev` which is required by `libui`.
 * Run `make update_libs` to get all required golang dependencies.

### Building Linux version

`make linux`

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

## TODO

Things to consider later:

 * Avoid asking administrator rights (Windows): [Github](https://stackoverflow.com/questions/31558066/how-to-ask-for-administer-privileges-on-windows-with-go)

## License

[MIT](https://opensource.org/licenses/MIT)

## Changelog

### 1.4.1 (05-SEP-2018)

 * Always create debug file to `~ktp/naksu_lastlog.txt` or `TEMP/naksu_lastlog.txt`
 * Give more informative error message when vagrant exits with `--macaddress` error message
 * Bugfix: get Linux temporary path from `TempDir()` instead of `TEMP` env var

### 1.4.0 (31-AUG-2018)

 * Show current server version
 * Show status messages when working with vagrant/VBoxManage
 * UI Language support (Finnish, Swedish)

### 1.3.0 (13-AUG-2018)

 * Restart with `x-terminal-emulator` if started via file manager (Linux)
 * Warn user if there is less than 5 Gb free disk before install, update or backup
 * Bugfix: Don't panic if media does not have vendor/model (Linux)
 * Get vendor/model strings for removable media directly from WMI without `wmic` (Windows)
 * Application icon (Windows)
 * Added logging (use `-debug` switch to show debug messages)

### 1.2.0 (09-AUG-2018)

 * User can create a backup (clone) from VM hard drive
 * Continues execution although binary update fails

### 1.1.0 (05-JUL-2018)

 * User can switch between Abitti and Matriculation Exam servers
