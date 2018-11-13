mkdir workspace
cd workspace
chester init
tree
chester create --silent 'echo "hello world"'
chester test --silent
cd ..
rm -rf workspace
