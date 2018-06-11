# naksu

`naksu` is collection of helper scripts for MEB diskless Abitti/matriculation exam servers.
In real life Naksu is Onerva's (the [Abitti model teacher](https://www.abitti.fi/fi/tutustu/)) cat.

The need for some kind of helper scripts appeared to be evident as we followed support requests
and feedback from school IT staff and teachers.

These scripts are currently under planning/proof-of-concept stage. You can compile scripts in
Linux environment. To compile you need to have golang 1.7 or newer. In Debian/Ubuntu environment
install `golang-1.9` or `golang-1.10` before compiling scripts with `make all`.

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
 1. TODO: Update `naksu` scripts

### Start Virtual Server

 1. Make sure you have `vagrant` executable
 1. Make sure Oracle VirtualBox is installed
 1. Execute `vagrant up`

### TODO: Switch Between Abitti and Matriculation Examination Server

Almost the same as Update but you have to be able to select between Abitti-version `Vagrantfile` and the matriculation examination `Vagrantfile`. The latter is downloaded by the school principal.

### TODO: Make Backup of Server Hard Drive

1. Make sure you have `vagrant` executable
1. Make sure Oracle VirtualBox is installed
1. Make sure you have `VBoxManage` executable
1. `VBoxManage list hdds` -> Get UUID of the disk
1. `VBoxManage clonemedium {UUID} {destination path} --format VMDK`

## TODO

Things to consider later:

 * Icon and avoid asking administrator rights: https://stackoverflow.com/questions/31558066/how-to-ask-for-administer-privileges-on-windows-with-go
 * GUI: https://github.com/andlabs/ui

## License

[MIT](https://opensource.org/licenses/MIT)
