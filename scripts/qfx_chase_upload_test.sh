if [ "$#" -ne 1 ]; then
    echo "Invalid Parameters, usage: ./update_user_to_admin <user_id>"
    exit 1
fi

if (curl localhost:8080/query --cookie "auth-token=$1" \
  -F operations='{ "query": "mutation ($file: Upload!) { chaseOFXTransactionsUpload(file: $file) }", "variables": { "file": null } }' \
  -F map='{ "0": ["variables.file"] }' \
  -F 0=@$PWD/scripts/test-data/bank_transactions.QFX); then
  echo "done"
else
    exit 1
fi

if (curl localhost:8080/query --cookie "auth-token=$1" \
  -F operations='{ "query": "mutation ($file: Upload!) { chaseOFXTransactionsUpload(file: $file) }", "variables": { "file": null } }' \
  -F map='{ "0": ["variables.file"] }' \
  -F 0=@$PWD/scripts/test-data/cc_transactions.QFX); then
  echo "done"
else
    exit 1
fi
