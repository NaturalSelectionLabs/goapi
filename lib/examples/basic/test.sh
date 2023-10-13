echo
echo
echo '== print openapi doc =='

curl localhost:3000/openapi.json

echo
echo
echo '== favicon =='

curl localhost:3000/favicon.ico -i

echo
echo
echo '== login =='

curl localhost:3000/login -id '{"username": "a@a.com", "password": "123456"}'

echo
echo
echo '== get posts =='

curl 'localhost:3000/users/3/posts?keyword=sky&t=game' -H 'Cookie: token=123456'

echo
echo
echo '== validation error =='

curl 'localhost:3000/users/3/posts?keyword=A&t=game' -H 'Cookie: token=123456'
