# HGate

This is a simple proxy for interacting with the Horizon server, which performs the signature of the request headers for authentication and also the formation of some transactions (currently only Payment).



## Build sources 

> Before start you need to [install Go](https://golang.org/doc/install).

To get HGate executable just run: 
```bash
go install gitlab.com/swarmfung/hgate 
```

## Usage

Run the HGate executable with argument `serve` to start proxy:

```bash
./hgate serve
```

Additional option:

- `--config [path/to/config] or -c [path/to/config]` default: "./config.yaml"

## Config
Config is a file in [YAML](https://en.wikipedia.org/wiki/YAML) format.

```yaml
port: 8842   #port where HGate will listen
log_level: warn   #level of log output
horizon_url: http://128.199.229.140:8011    #url of Horizon API
seed: SADB...   #your SecretSeed
```

## Available operations

### Payment
Send **POST** request to `http://localhost:8842/send_payment` with following JSON body:

- Method - POST
- Content Type - application/json

```json
{
    "destination_id": "GDTWK...",
    "amount": "10.4242",
    "asset": "SUN",
    "subject": "Proxy Test",
    "reference": "Some unique reference",
    "pay_fee_instead_dest": false
}
```

where:

- `destination_id` - accountId of recipient of payment;
- `amount` - number of coins;
- `asset` - asset code;
- `subject` - some description for payment;
- `reference` - unique value;
- `pay_fee_instead_dest` - true - if your want to pay recipient fee, otherwise the "false"


Response:

- Success:
    ``` json
    {
        "type": "success",
        "title": "Success",
        "status": 200
    }
    ```

- Bad Request:
    ``` json
    {
        "type": "bad_request",
        "title": "Bad Request",
        "status": 400,
        "detail": "The request you sent was invalid in some way",
        "extras": {
            "invalid_field": "field_name",
            "reason": "some reasons"
        }
    }
    ```

- Submission Error:
    ```json
    {
        "type": "",
        "title": "Transaction Failed",
        "status": 400,
        "detail": "The transaction failed when submitted to the network. The `extras.result_codes` field on this response contains further details.",
        "extras": {
            "envelope_xdr": "AAAAAEDI....",
            "result_codes": {
                "transaction": "tx_failed",
                "operations": [
                    "op_malformed"
                ]
            },
            "result_xdr": "AAAAAAAAAAD/////AAAAAQAAAAAAAAAB/////wAAAAA="
        }
    }
    ```

| Error | Description |
| ----- | ----------- |
|`op_malformed` | the operation was malformed in some way.|
|`op_underfunded` | the operation failed due to a lack of funds. |
|`op_reference_duplication` | the payment with the same reference already submitted|
|`op_stats_overflow` | the amount of payments is beyond the allowed limits |
|`op_limits_exceeded` | the amount of payments is beyond the allowed limits|
|`op_fee_mismatched` |the fee thats specified in operation incorrect |
|`op_balance_not_found` | the destination balance_id is not exist in the system|
|`op_balance_account_mismatched` | the source balance_id doesn't match with source account |
|`balance_assets_mismatched` | one of the balance_ids from operation doesn't match with the asset of payment|
|`src_balance_not_found` | the source balance_id is not exist in the system|
