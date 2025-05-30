pkgname=xbox-controller-battery-linux
pkgver=1.0.0
pkgrel=1
arch=('x86_64')
url="https://github.com/nvhai245/xbox-controller-battery-linux"
license=('MIT')
makedepends=('go')
source=()
sha256sums=()

build() {
	cd "$srcdir"
	go build -o xbox-controller-battery-linux .
}

package() {
	install -Dm755 "$srcdir/xbox-controller-battery-linux" "$pkgdir/usr/bin/xbox-controller-battery-linux"
	install -Dm644 "$srcdir/xbox-controller-battery-linux.desktop" "$pkgdir/usr/share/applications/xbox-controller-battery-linux.desktop"
	install -Dm644 "$srcdir/xbox-controller-battery-linux.png" "$pkgdir/usr/share/icons/hicolor/64x64/apps/xbox-controller-battery-linux.png"
	cp -r "$srcdir/icons" "$pkgdir/usr/share/icons/hicolor/64x64/apps/icons"
}
