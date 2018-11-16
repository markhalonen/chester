mkdir workspace
cd workspace
chester init

mkdir my_test
cd my_test
echo 'import json' > script.py
echo 'print json.dumps({"k1": "v1"})' >> script.py
echo 'python script.py' > command.sh
chmod 777 command.sh
cd ..

chester create --silent my_test
chester test --silent

# Change the content on purpose so it fails
echo '{"k1":"v2"}' > ./__chester__/tests/1/expected_output.txt

chester test --silent
cd ..
rm -rf workspace