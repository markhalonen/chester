mkdir workspace
cd workspace
git init | grep -wo Initialized
tree -a
echo "Hello git" > new_file.txt
git add new_file.txt
git commit -m "A great new file" | grep -wo "new_file.txt"
git log | grep -wo "A great new file"
tree
cd ..
rm -rf workspace
