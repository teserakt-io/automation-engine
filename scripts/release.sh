#!/bin/sh

GIT_TAG=$(git describe --exact-match HEAD 2>/dev/null || true)
GIT_BRANCH=`git branch | sed -n '/\* /s///p'`


#if ! [[ "${GIT_BRANCH}" == "master" ]]; then
#    echo "Releases are only performed on master!"
#    exit 1
#fi

if [[ -z "${VERSION}" && -z "${GIT_TAG}" ]]; then
    echo "Release not tagged, refusing to build."
    exit 1
fi

if ! [[ -z "${VERSION}" ]]; then
    V=$VERSION
elif ! [[ -z "${GIT_TAG}" ]]; then
    V=$GIT_TAG
else
    echo "Bug in release script."
    return 1
fi

OUTDIR=build/$V

echo "Producing release $GIT_TAG"

mkdir -p $OUTDIR

OUTDIR=$OUTDIR/linux_amd64/ GOOS=linux GOARCH=amd64 ./scripts/build.sh
OUTDIR=$OUTDIR/darwin_amd64/ GOOS=darwin GOARCH=amd64 ./scripts/build.sh
OUTDIR=$OUTDIR/windows_amd64/ GOOS=windows GOARCH=amd64 ./scripts/build.sh

LINUX_BINARIES=$(find $OUTDIR/linux_amd64 -type f -executable)
for bin in $LINUX_BINARIES; do 
    if [[ $(echo $(ldd $bin) | grep -c "\.so") -gt 1 ]]; then 
        echo "Build failed. Dynamic executable created on linux."
        exit -1
    else
        echo "$bin is a static binary."
    fi
done

mkdir -p $OUTDIR/configs/

mkdir -p $OUTDIR/configs/
cp -v configs/config.yaml.example $OUTDIR/configs/

pushd build/$V
tar cjf ../e4-ae-$V.tar.gz *
popd
