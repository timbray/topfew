# Installing Topfew

Each Topfew [release](https://github.com/timbray/topfew/releases) comes with binaries built for both the x86 and ARM
flavors of Linux, MacOS, and Windows.

Topfew comes with a Makefile which is uncomplicated. Typing `make` will create an executable named `topfew`, 
created by `go build` with no options, in the `./bin` directory.

## Arch Linux

Topfew [is available](https://aur.archlinux.org/packages/topfew) in the 
[Arch User Repository](https://wiki.archlinux.org/title/Arch_User_Repository) (AUR).
If you have an AUR pacman wrapper installed you can install it directly. Otherwise, to install Topfew as an Arch package: 
```
git clone https://aur.archlinux.org/topfew.git
cd topfew
makepkg -i
```

## NixOS

Topfew [is available on NixOS](https://search.nixos.org/packages?show=topfew).

## Homebrew

On MacOS, Topfew is [is available on Homebrew](https://formulae.brew.sh/formula/topfew).
