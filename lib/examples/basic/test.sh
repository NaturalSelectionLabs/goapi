echo
echo
echo '== print openapi doc =='

curl localhost:3000 -H 'Accept: application/json'

echo
echo
echo '== login =='

curl localhost:3000/login -id '{"username": "admin", "password": "123456"}'

echo
echo
echo '== get posts =='

curl 'localhost:3000/users/3/posts?keyword=sky&t=game' -H 'Cookie: token=123456'

echo
echo
echo '== validation error =='

curl 'localhost:3000/users/3/posts?keyword=A&t=game' -H 'Cookie: token=123456'
