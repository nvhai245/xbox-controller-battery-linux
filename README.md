
# Xbox Controller Battery for Linux

A simple system tray app that displays the battery level of your Xbox controller on Linux.

---

## 🚀 Installation

### 📦 Download Prebuilt Packages

Go to the **[Releases](https://github.com/nvhai245/xbox-controller-battery-linux/releases)** page and download the latest version for your system:

- **Debian/Ubuntu:** `xbox-controller-battery-linux_*.deb`
- **Arch Linux:** `xbox-controller-battery-linux-*.pkg.tar.zst`

---

### 🐧 Debian / Ubuntu

1. Download the `.deb` package from the [Releases](https://github.com/nvhai245/xbox-controller-battery-linux/releases) page.  
2. Install it using:

```bash
sudo dpkg -i xbox-controller-battery-linux_*.deb
sudo apt-get install -f  # fix any missing dependencies
```

---

### 🅰️ Arch Linux

1. Download the `.pkg.tar.zst` package.  
2. Install it using:

```bash
sudo pacman -U xbox-controller-battery-linux-*.pkg.tar.zst
```

---

## 🖥️ Run the App

After installation, launch the app from your application menu or run:

```bash
xbox-controller-battery-linux
```

It will show in the system tray and display the battery level of your Xbox controller.

---

## 🛠️ Build From Source

If you want to build the app yourself:

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
