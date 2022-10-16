<div align="center">
  <img src="almost-logo.png" height="120">
  <h1 align="center">almost</h1>
  <p align="center">An on-demand immutability based on file attributes.</p>
</div>

# almost
An on-demand immutability for VanillaOS.

> **Note**: This is a work in progress. It is not ready for production use.

### Read here
This program is meant to be used with [apx](https://github.com/vanilla-os/apx), 
an apt replacement for VanillaOS.

To use it with other distributions, be sure to re-pack it with the correct
package manager set:
- `"Almost::PkgManager::EntryPoint": "/usr/bin/apt"`

### Help
```
Usage: 
almost [options] [command]

Options:
	--help/-h		show this message
	--verbose/-v		show more verbosity
	--version/-V		show version

Commands:
	enter			set the filesystem as ro or rw until reboot
	config			show the current configuration
	check			check whether the filesystem is read-only or read-write
```
