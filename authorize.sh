JWT=$(curl -plaintext -u gomicrojas123@gmail.com:Admin123 -i http://payment.com/login)
curl -X POST -H "Authorization: Bearer $JWT" --data @authorize_payload.json http://payment.com/customer/payment/authorize
