if [ "$#" -ne 1 ]; then
    echo "Invalid Parameters, usage: ./update_user_to_admin <user_id>"
    exit 1
fi

if (curl localhost:8080/query --cookie "auth-token=$1" \
  -F operations='{ "query": "mutation ($file: Upload!) { chaseCSVTransactionsUpload(file: $file) }", "variables": { "file": null } }' \
  -F map='{ "0": ["variables.file"] }' \
  -F 0=@$PWD/scripts/test-data/transactions.notcsv); then
  echo "done"
else
    exit 1
fi
