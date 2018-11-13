mkdir workspace
cd workspace
chester init
mkdir my_test
cd my_test
echo 'print "Hello From Python!"' > script.py
echo 'python script.py' > command.sh
chmod 777 command.sh
cd ..
chester create --silent my_test
chester test --silent
cd ..
rm -rf workspace
