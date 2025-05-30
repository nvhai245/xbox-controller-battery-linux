pkgname=xbox-controller-battery-linux
pkgver=1.0.0
pkgrel=1
arch=('x86_64')
url="https://github.com/nvhai245/xbox-controller-battery-linux"
license=('MIT')
makedepends=('go')
source=("${pkgname}-${pkgver}.tar.gz")
sha256sums=('SKIP')

build() {
	cd "${srcdir}/${pkgname}-${pkgver}"
	go build -o xbox-controller-battery-linux
}

package() {
	install -Dm755 "${srcdir}/${pkgname}-${pkgver}/xbox-controller-battery-linux" "${pkgdir}/usr/bin/xbox-controller-battery-linux"
	install -Dm644 "${srcdir}/${pkgname}-${pkgver}/xbox-controller-battery-linux.desktop" "${pkgdir}/usr/share/applications/xbox-controller-battery-linux.desktop"
	install -Dm644 "${srcdir}/${pkgname}-${pkgver}/icons/light/battery_high.png" "${pkgdir}/usr/share/icons/hicolor/64x64/apps/xbox-controller-battery-linux.png"
}
