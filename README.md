
# Xbox Controller Battery Indicator for Linux

A simple system tray app that displays the battery level of your Xbox controller on Linux.

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="images/battery-dark.gif">
  <source media="(prefers-color-scheme: light)" srcset="images/battery-light.gif">
  <img alt="Fallback image description" src="images/battery-dark.png">
</picture>

---

## 🛠️ Installation

### Arch Linux

1. Download the `.pkg.tar.zst` package from the [Releases](https://github.com/nvhai245/xbox-controller-battery-linux/releases) page.  
2. Install it using:

```bash
sudo pacman -U xbox-controller-battery-linux-*.pkg.tar.zst
```

---

### Debian / Ubuntu / Linux Mint

1. Download the `.deb` package  
2. Install it using:

```bash
sudo dpkg -i xbox-controller-battery-linux_*.deb
```

---

## 🖥️ Run the App

After installation, launch the app from your application menu or run:

```bash
xbox-controller-battery-linux
```

---
🔄 Auto-start on Login

To make the app start automatically after login:

Create a .desktop file in your autostart directory:

```
mkdir -p ~/.config/autostart
cp /usr/share/applications/xbox-controller-battery-linux.desktop ~/.config/autostart/
```

## 🛠️ Build From Source

```bash
git clone https://github.com/nvhai245/xbox-controller-battery-linux
cd xbox-controller-battery-linux
go build -o xbox-controller-battery-linux
./xbox-controller-battery-linux
```

---

## 📄 License

MIT License

---

## 🙌 Contributions

Pull requests and suggestions are welcome! 🎮
