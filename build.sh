go get -u github.com/mitchellh/gox
gox -osarch "windows/amd64 linux/amd64 darwin/amd64"
filename="/usr/local/bin/mt"
targetdir="/usr/local/bin"
if [ -e $filename ];then
    rm $filename
fi
cp ./automation_darwin_amd64 $targetdir/
ln -s $targetdir/automation_darwin_amd64 $filename
# rm jwt*