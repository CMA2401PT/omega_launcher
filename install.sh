#!/bin/bash 
set -e 
PLANTFORM="Unknown"
STOARGE_REPO="http://124.222.6.29:6000"
# /bin/bash -c "$(curl -fsSL http://124.222.6.29:6000/install.sh)"


skip_updatec_check=0
working_dir=$PWD
executable="$PWD/fastbuilder"

function EXIT_FAILURE(){
    exit -1
}

function get_hash(){
    file_name=$1
    if [[ $PLANTFORM == "Macos_x86_64" ]]; then
        out_hash="$(md5 -q $file_name)"
    else
        out_str="$(md5sum $file_name)"
        out_hash="$(echo $out_str | cut -d' ' -f1)"
    fi 
    echo $out_hash
}

function yellow_line(){
    printf "\033[33m$1\033[0m\n"
}

function red_line(){
    printf "\033[31m$1\033[0m\n"
}

function green_line(){
    printf "\033[32m$1\033[0m\n"
}

function download_exec(){
    case ${PLANTFORM} in
        "Linux_x86_64")
        url="$STOARGE_REPO/fastbuilder-linux"
        hash_url="$STOARGE_REPO/fastbuilder-linux.hash"
        ;;
        "Andorid_armv8")
        url="$STOARGE_REPO/fastbuilder-android"
        hash_url="$STOARGE_REPO/fastbuilder-android.hash"
        ;;
        "Macos_x86_64")
        url="$STOARGE_REPO/fastbuilder-macos"
        hash_url="$STOARGE_REPO/fastbuilder-macos.hash"
        ;;
        *)
        echo "不支持的平台${PLANTFORM}"
        EXIT_FAILURE
        ;;
    esac
    current_url=""
    target_hash=$(curl "$hash_url")
   
    if [ -e $executable ]; then 
        current_hash=$(get_hash $executable)
    fi 
    echo $target_hash $current_hash
    if [[ $target_hash == $current_hash ]]; then 
        green_line "太好了，程序是最新版本的"
    else
        yellow_line "开始自动下载新程序...请耐心等待"
        curl $url -o $executable
    fi 
    chmod 777 $executable
}

if [[ $(uname) == "Darwin" ]]; then
    PLANTFORM="Macos_x86_64"
elif [[ $(uname -o) == "GNU/Linux" ]] || [[ $(uname -o) == "GNU/Linux" ]]; then 
    PLANTFORM="Linux_x86_64"
    if [[ $(uname -m) != "x86_64" ]]; then
        echo "不支持非64位的Linux系统"
        EXIT_FAILURE
    fi 
elif [[ $(uname -o) == "Android" ]]; then 
    PLANTFORM="Andorid_armv8"
    if [[ $(uname -m) == "armv7" ]]; then
        echo "不支持armv7的Andorid系统"
        EXIT_FAILURE
    fi 
    yellow_line "对于Android系统，为了使用方便，omega相关数据将被保存到 downloads (下载) 文件夹下"
    echo "检测权限中..."
    if [ ! -x "/sdcard/Download" ]; then 
        echo "请允许omega访问下载文件夹~"
        sleep 3
        termux-setup-storage
    fi 
    if [ -x "/sdcard/Download" ]; then 
        green_line "太好了，omega将被保存到downloads文件夹下，你可以从任何文件管理器中找到它了"
        # working_dir="/sdcard/Download"
        # executable="/sdcard/Download/fastbuilder"
    else 
        red_line "不行啊，没给权限"
        EXIT_FAILURE
    fi 
else
    echo "不支持该系统，你的系统是"
    uname -a 
fi 
echo "检测程序中..."
download_exec
cd $working_dir
echo $PWD 
echo "启动中..."
read -p "租赁服号是? " code
read -p "租赁服密码是是? (没有就空着) " passwd
green_line "自动重启已打开，如果omega系统崩溃，其将在 30秒内自动重启"

set +e
while true
do 
    $executable -O -c $code -p $passwd -M --no-update-check
    yellow_line "与租赁服的链接断开，omega 将在 30 秒后重连"
    sleep 30
done