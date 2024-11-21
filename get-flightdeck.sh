#!/bin/bash

: ${BINARY_NAME:="flightdeck"}
: ${DEBUG:="false"}
: ${INSTALL_DIR:="/usr/local/bin"}
: ${INSECURE:="false"}
: ${RELEASE:="latest"}

if [ "$DEBUG" == "1" ]; then
	set -x
fi

TMP_DIR=$(mktemp -d -t flightdeck-install-XXXXXXXXXX)

function cleanup {
	rm -rf $TMP_DIR > /dev/null
}

function fail {
	cleanup
	msg=$1
	echo "============"
	echo "Error: $msg" 1>&2
	exit 1
}

function install {
	#settings
	PROG="flightdeck-cli"
	ASPROG="flightdeck"
	GH="https://github.com"

	#bash check
	[ ! "$BASH_VERSION" ] && fail "Please use bash instead"
	[ ! -d $INSTALL_DIR ] && fail "output directory missing: $INSTALL_DIR"

	#dependency check, assume we are a standard POISX machine
	which find > /dev/null || fail "find not installed"
	which xargs > /dev/null || fail "xargs not installed"
	which sort > /dev/null || fail "sort not installed"
	which tail > /dev/null || fail "tail not installed"
	which cut > /dev/null || fail "cut not installed"
	which du > /dev/null || fail "du not installed"

	#choose an HTTP client
	GET=""
	if which curl > /dev/null; then
		GET="curl"
		if [[ $INSECURE = "true" ]]; then GET="$GET --insecure"; fi
		GET="$GET --fail -# -L"
	elif which wget > /dev/null; then
		GET="wget"
		if [[ $INSECURE = "true" ]]; then GET="$GET --no-check-certificate"; fi
		GET="$GET -qO-"
	else
		fail "neither wget/curl are installed"
	fi

	#debug HTTP
	if [ "$DEBUG" == "1" ]; then
		GET="$GET -v"
	fi

	#find OS #TODO BSDs and other posixs
	case `uname -s` in
	Darwin) OS="darwin";;
	Linux) OS="linux";;
	*) fail "unknown os: $(uname -s)";;
	esac
	#find ARCH
	if uname -m | grep -E '(arm|arch)64' > /dev/null; then
		ARCH="arm64"
		
	elif uname -m | grep 64 > /dev/null; then
		ARCH="amd64"
	elif uname -m | grep arm > /dev/null; then
		ARCH="arm" #TODO armv6/v7
	elif uname -m | grep 386 > /dev/null; then
		ARCH="386"
	else
		fail "unknown arch: $(uname -m)"
	fi
	
    #choose from asset list
	URL="https://github.com/arctir/flightdeck-cli/releases/download/${RELEASE}"
	FTYPE=""

    if [ "$RELEASE" == "latest" ]; then
        URL="https://github.com/arctir/flightdeck-cli/releases/latest/download"
   fi 
	case "${OS}_${ARCH}" in
	"darwin_amd64")
		URL="${URL}/flightdeck-cli_Darwin_x86_64.tar.gz"
		FTYPE=".tar.gz"
		;;
	"darwin_arm64")
		URL="${URL}/flightdeck-cli_Darwin_arm64.tar.gz"
		FTYPE=".tar.gz"
		;;
	"linux_386")
		URL="${URL}/flightdeck-cli_Linux_i386.tar.gz"
		FTYPE=".tar.gz"
		;;
	"linux_amd64")
		URL="${URL}/flightdeck-cli_Linux_x86_64.tar.gz"
		FTYPE=".tar.gz"
		;;
	"linux_arm64")
		URL="${URL}/flightdeck-cli_Linux_arm64.tar.gz"
		FTYPE=".tar.gz"
		;;
	*) fail "No asset for platform ${OS}-${ARCH}";;
	esac

	#got URL! download it...
	echo -n "Downloading flightdeck"
	echo -n " (${OS}/${ARCH})"
	echo "....."
	
	#enter tempdir
	mkdir -p $TMP_DIR
	cd $TMP_DIR
	if [[ $FTYPE = ".gz" ]]; then
		which gzip > /dev/null || fail "gzip is not installed"
		bash -c "$GET $URL" | gzip -d - > $PROG || fail "download failed"
	elif [[ $FTYPE = ".bz2" ]]; then
		which bzip2 > /dev/null || fail "bzip2 is not installed"
		bash -c "$GET $URL" | bzip2 -d - > $PROG || fail "download failed"
	elif [[ $FTYPE = ".tar.bz" ]] || [[ $FTYPE = ".tar.bz2" ]]; then
		which tar > /dev/null || fail "tar is not installed"
		which bzip2 > /dev/null || fail "bzip2 is not installed"
		bash -c "$GET $URL" | tar jxf - || fail "download failed"
	elif [[ $FTYPE = ".tar.gz" ]] || [[ $FTYPE = ".tgz" ]]; then
		which tar > /dev/null || fail "tar is not installed"
		which gzip > /dev/null || fail "gzip is not installed"
		bash -c "$GET $URL" | tar zxf - || fail "download failed"
	elif [[ $FTYPE = ".zip" ]]; then
		which unzip > /dev/null || fail "unzip is not installed"
		bash -c "$GET $URL" > tmp.zip || fail "download failed"
		unzip -o -qq tmp.zip || fail "unzip failed"
		rm tmp.zip || fail "cleanup failed"
	elif [[ $FTYPE = ".bin" ]]; then
		bash -c "$GET $URL" > "flightdeck-cli_${OS}_${ARCH}" || fail "download failed"
	else
		fail "unknown file type: $FTYPE"
	fi

    #search subtree largest file (bin)
	TMP_BIN=${TMP_DIR}/flightdeck-cli
    if [ ! -f "$TMP_BIN" ]; then
        fail "could not find find binary"
    fi

	#move into PATH or cwd
	chmod +x $TMP_BIN || fail "chmod +x failed"
	DEST="$INSTALL_DIR/$PROG"	
	if [ ! -z "$ASPROG" ]; then
		DEST="$INSTALL_DIR/$ASPROG"
	fi
	#move without sudo
	OUT=$(mv $TMP_BIN $DEST 2>&1)
	STATUS=$?
	# failed and string contains "Permission denied"
	if [ $STATUS -ne 0 ]; then
		if [[ $OUT =~ "Permission denied" ]]; then
			echo "Install requires sudo..."
			sudo mv $TMP_BIN $DEST || fail "sudo mv failed" 
		else
			fail "mv failed ($OUT)"
		fi
	fi
	echo "Downloaded to $DEST"
	#done
	cleanup
}
install
