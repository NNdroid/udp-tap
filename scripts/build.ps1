$RELEASE_BIN_DIR='.\bin\'
$BINARY_NAME='udp-tap'
$PACKAGE_NAME='udp-tap/pkg/common'
function removeDir() {
    Remove-Item -Path $RELEASE_BIN_DIR -Recurse
}

function createDir() {
    if (Test-Path -Path $RELEASE_BIN_DIR) {
        echo "path exists"
    } else {
        echo "path not exists"
        New-Item -Path $RELEASE_BIN_DIR -ItemType Directory
    }
}

function goBuild() {
    param(
        [string]$os,
        [string]$arch
    )
    $suffix=''
    if ($os -like "windows") {
        $suffix='.exe'
    }
    $versionCode=Get-Date -format "yyyyMMdd"
    $goVersion=go version
    $gitHash=git log --pretty=format:'%h' -n 1
    $buildTime=git log --pretty=format:'%cd' -n 1
    set CGO_ENABLED=0
    go env -w CGO_ENABLED=0
    set GOOS=$os
    go env -w GOOS=$os
    set GOARCH=$arch
    go env -w GOARCH=$arch
    go build -o $RELEASE_BIN_DIR$BINARY_NAME-${os}_$arch$suffix -ldflags "-w -s -X '$PACKAGE_NAME.Version=1.0.$versionCode' -X '$PACKAGE_NAME.GoVersion=$goVersion' -X '$PACKAGE_NAME.GitHash=$gitHash' -X '$PACKAGE_NAME.BuildTime=$buildTime'" ./cmd/main.go
}

function main() {
    removeDir
    go clean
    go mod tidy
    createDir
    goBuild linux amd64
    goBuild linux arm64
    goBuild linux 386
    goBuild linux arm
}

main