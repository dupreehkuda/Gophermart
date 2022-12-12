# Gophermart bonus accrual service

#### *DISCLAIMER* 
*This API is created as a task on a Yandex Practicum course. Some decisions were made just because of a specific assignment.*

## Description

This is a bonus accrual system and it works in pair with accrual service in ./cmd/accrual. For more detailed view on the task you can read [this](SPECIFICATION.md). 

The service can register/login user and return a JWT-token in cookies for future operations, add order or withdrawal, get all user's orders and withdrawals.

## Architecture

I have implemented a three layer architecture: `handler` -> `action` -> `storage`. (Action is a layer of business-logic.)

This architecture allows us to change each part as we want and it will work while each part implements specific [interface](internal/interfaces/interfaces.go).

There is also service layer which is a layer for connecting to accrual service. Only business-logic knows about it.

## Endpoints

### Get endpoints

- `GET /api/user/balance`
  - Handler for getting user balance
  - Only accessible after authorization
  - Example output:
  ```JSON
   {
	   "Current": 469.11,
	   "Withdrawn": 0
   }
  ```

- `GET /api/user/orders`
  - Handler for showing all user's orders
  - Only accessible after authorization
  - Example output:
  ```JSON
   [
	   {
		   "number": "8942677371",
		   "status": "PROCESSED",
		   "accrual": 469.11,
		   "uploaded_at": "2022-10-31T20:29:23Z"
	   }
   ]
  ```

- `GET /api/user/balance/withdrawals`
  - Handler for getting all user's withdrawals (user's spent points)
  - Only accessible after authorization
  - Example output:
  ```JSON
   [
      {
         "order": "8942677371",
         "sum": 111.25,
         "processed_at": "2022-10-31T20:29:39Z"
      }
   ]
  ```
  
### Post endpoints

- `POST /api/user/register`
  - Handler for registering new user
  - Returns JWT-token in cookies if your login is unique
  - Example input:
   ```JSON
   {
      "login": "first",
      "password": "37"
   }
  ```

- `POST /api/user/login`
  - Handler for authorizing existing user
  - Returns JWT-token in cookies if your credentials are correct
  - Example input:
   ```JSON
   {
      "login": "first",
      "password": "37"
   }
  ```

- `POST /api/user/orders`
  - Handler for adding new order to process it
  - Validates order number with Luhn algorithm

- `POST /api/user/balance/withdraw`
  - Handler for withdrawing balance for an order
  - Works if there is enough funds
  - Example input:
  ```JSON
   {
      "order": "8942677371",
      "sum": 111.25
   }
  ```

All test cases can be found in [this file](insomnia_requests.json). It is an insomnia requests collection. Every endpoint has a request there.

## Deploy

To deploy just use 
```sh
docker-compose up
```

This will run postgres, accrual service and gophermart. 

## Conclusion

I am very bad at describing things I made some time ago. I hope that nobody will read it and it will just serve me as a lesson.

This project showed me that there is nothing hard in software development (hope I would not regret of what I said).
Again 80% of the work is just about designing the service: sketches on miro board, re-reading the task several times, thinging of every litle thing and so on. Other 20% is just about coding and debugging.


