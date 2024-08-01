while true ;
    do 
    curl -X GET http://127.0.0.1:8080/art?id=1  -D - -o /dev/null
done