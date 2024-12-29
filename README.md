# Hello World Go Project

This is a simple Go project to verify that the Go development environment is set up correctly.

## How to Run

1. Navigate to the project directory:
    ```sh
    cd hello-world
    ```

2. Run the Go program:
    ```sh
    go run main.go
    ```



Doubts


# Requirement 1

1 - I will store transactions right? I need to care about the product traded? 
2 - In general a transaction is related to a owner (user or other system...). I need to save this too? The owner of the transaction.
3 - The transaction can change? 

# Requirement 2

https://fiscaldata.treasury.gov/datasets/treasury-reporting-rates-exchange/treasury-reporting-rates-of-exchange

https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange?fields=record_date,country,exchange_rate,currency,effective_date&filter=country:eq:Argentina&sort=-record_date&page[number]=1&page[size]=1&format=json

page[number]=3&page[size]=2&
