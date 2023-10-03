cu=`pwd`
rm -rf release/packages && echo "success deletion" || echo "not success"
mkdir -p release/packages/
ver_path='github.com/h2p2f/dedicated-vault/internal/client/config.version'
build_path='github.com/h2p2f/dedicated-vault/internal/client/config.buildDate'
ca_path='github.com/h2p2f/dedicated-vault/internal/client/config.ca'
cert_path='github.com/h2p2f/dedicated-vault/internal/client/config.cert'
key_path='github.com/h2p2f/dedicated-vault/internal/client/config.key'
ver_value='0.0.4'
build_value='2023-10-03'
ca_value='/tmp/dedicated-vault/crypto/ca-cert.pem'
cert_value='/tmp/dedicated-vault/crypto/client-cert.pem'
key_value='/tmp/dedicated-vault/crypto/client-key.pem'
os_all='linux windows darwin freebsd'
arch_all='amd64 arm64'
for os in $os_all; do
    for arch in $arch_all; do
      set GOOS=$os
      set GOARCH=$arch
      if [ $os = "windows" ]; then
        go build -ldflags "-X "$ver_path"="$ver_value" -X "$build_path"="$build_value" -X "$ca_path"="$ca_value" -X "$cert_path"="$cert_value" -X "$key_path"="$key_value  -o $os"_"$arch".exe" && echo "Success build for arch "$arch" and os "$os || echo "No problem"
        mv $os"_"$arch".exe" release/packages && echo "Move success" || echo "Move not success"
      else
        go build -ldflags "-X "$ver_path"="$ver_value" -X "$build_path"="$build_value" -X "$ca_path"="$ca_value" -X "$cert_path"="$cert_value" -X "$key_path"="$key_value -o $os"_"$arch && echo "Success build for arch "$arch" and os "$os || echo "No problem"
        mv $os"_"$arch release/packages && echo "Move success" || echo "Move not success"
      fi
    done
done
echo "Success Build"
cd $cu