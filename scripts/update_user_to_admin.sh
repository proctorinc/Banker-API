if [ "$#" -ne 1 ]; then
    echo "Invalid Parameters, usage: ./update_user_to_admin <user_id>"
    exit 1
fi

psql -U mattyp -d chase-data -c "UPDATE users SET role = 'ADMIN' WHERE id = '$1';"
