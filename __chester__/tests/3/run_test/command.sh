mkdir workspace
cd workspace
chester init
chester create --silent "python no_file.py" | grep -wo "python: can't open file 'no_file.py'"
cd ..
rm -rf workspace
