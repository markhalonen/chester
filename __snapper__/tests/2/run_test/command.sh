mkdir workspace
cd workspace
snapper init
mkdir my_test
cd my_test
echo 'print "Hello From Python!"' > script.py
echo 'python script.py' > command.sh
chmod 777 command.sh
cd ..
snapper create --silent my_test
snapper test --silent
cd ..
rm -rf workspace
