# Packer Hyper-V Examples

## Basic Examples

### Windows Server 2008 Standard

This example creates Windows Server 2008 Standard vagrant box.

* **win2008-standard.json** This is as basic as the configuration can be. The VM will have"
  * 127GB hard disk and 1GB RAM
  * floppy/win2008-standard/Autounattend.xml sets
    * Creates a administrator user with username: vagrant, password: vagrant
    * [Disables the computer password change](http://misheska.com/blog/2013/07/26/windows-7-automated-install-settings/#turn-off-computer-password)
    * [Disables the 'Set Network Location Prompt'](http://misheska.com/blog/2013/07/26/windows-7-automated-install-settings/#really-disable-set-network-location-prompt)
* **win2008-standard.bat** runs packer