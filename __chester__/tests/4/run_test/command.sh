mkdir workspace
cd workspace
chester init
chester create --silent "$(< ../echo_json.txt)"
chester test --silent
echo '{"k1":"v2"}' > ./__chester__/tests/1/expected_output.txt
chester test --silent
cd ..
rm -rf workspace