rm -rf pre-release
rm Omega系统更新*.zip
set -e 
mkdir pre-release
make clean 
make build/phoenixbuilder build/phoenixbuilder-windows-executable-x86_64.exe build/phoenixbuilder-android-executable-arm64 build/phoenixbuilder-macos-x86_64 -j4
cp build/phoenixbuilder ./pre-release/fastbuilder-linux
cp build/phoenixbuilder-windows-executable-x86_64.exe ./pre-release/fastbuilder-windows.exe
cp build/phoenixbuilder-macos-x86_64 ./pre-release/fastbuilder-macos
cp build/phoenixbuilder-android-executable-arm64 ./pre-release/fastbuilder-android

function get_hash(){
    fileName=$1
    outstr="$(md5sum $fileName)"
    hashStr="$(echo $outstr | cut -d' ' -f1)"
    echo "$hashStr"
}

echo $(get_hash ./pre-release/fastbuilder-linux) > ./pre-release/fastbuilder-linux.hash
echo $(get_hash ./pre-release/fastbuilder-windows.exe) > ./pre-release/fastbuilder-windows.hash
echo $(get_hash ./pre-release/fastbuilder-macos) > ./pre-release/fastbuilder-macos.hash
echo $(get_hash ./pre-release/fastbuilder-android) > ./pre-release/fastbuilder-android.hash

go run omega_release/compressor/main.go -in ./pre-release/fastbuilder-linux -out ./pre-release/fastbuilder-linux.brotli
go run omega_release/compressor/main.go -in ./pre-release/fastbuilder-windows.exe  -out ./pre-release/fastbuilder-windows.exe.brotli
go run omega_release/compressor/main.go -in ./pre-release/fastbuilder-macos  -out ./pre-release/fastbuilder-macos.brotli
go run omega_release/compressor/main.go -in ./pre-release/fastbuilder-android  -out ./pre-release/fastbuilder-android.brotli

cp omega_release/更新日志.txt ./pre-release
cp -r omega_release/新可用项 ./pre-release
cp omega_release/install.sh ./pre-release
cp omega_release/fileserver.go ./pre-release 
TimeStamp=$(date '+%m-%d_%H-%M')
touch ./pre-release/$TimeStamp
rsync -avP ./pre-release Tencent:/home/dai/
cd ./pre-release
zip -r pre-release.zip ./*
cd ..
mv ./pre-release/pre-release.zip Omega系统更新\(Beta@$TimeStamp\).zip