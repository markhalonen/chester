mkdir workspace
cd workspace
snapper init
tree
snapper create --silent 'echo "hello world"'
snapper test --silent
cd ..
rm -rf workspace
