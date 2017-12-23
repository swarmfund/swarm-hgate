# HGate

This is a simple proxy for interacting with the Horizon server, which performs the signature of the request headers for authentication and also the formation of some transactions (currently only Payment).

## Table of Contents:

1. [Install](#install)
2. [Config](#config)
2. [Usage](#usage)
    - [Run](#run)
    - [Get Ledger Changes](#get-ledger-changes)
    - [Send Payment](#send-payment)

## Install 

> Before continue you need to [install Go](https://golang.org/doc/install).

To get HGate executable just run: 

```bash
go get gitlab.com/swarmfung/hgate 
```

and test installation:

```bash
hgate
```

```
NAME:
   hgate - A simple proxy for interacting with the Horizon server

USAGE:
   hgate [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     serve    start proxy
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version 
```

## Config
Config is a file in [YAML](https://en.wikipedia.org/wiki/YAML) format.

```yaml
port: 8842          # port where HGate will listen
log_level: info     # level of log output
horizon_url: https://staging.api.sun.swarm.fund #url of Horizon API
account_id: GDSS... # AccountID of the account
seed: SCR...        # SecretSeed of the account signer
```

## Usage

### Run
Run the HGate executable with argument `serve` to start proxy:

```bash
./hgate serve
```

Additional option:

- `--config [path/to/config] or -c [path/to/config]` default: "./config.yaml"


### Get Ledger changes
```bash
    curl localhost:8842/ledger_changes 
```

```json
{
  "_links": {
    "self": {
      "href": "https://staging.api.sun.swarm.fund/ledger_changes?order=asc\u0026limit=10\u0026cursor="
    },
    "next": {
      "href": "https://staging.api.sun.swarm.fund/ledger_changes?order=asc\u0026limit=10\u0026cursor=1769526530048"
    },
    "prev": {
      "href": "https://staging.api.sun.swarm.fund/ledger_changes?order=desc\u0026limit=10\u0026cursor=927712940032"
    }
  },
  "_embedded": {
    "records": [
        {
            "id": "498216218624",
            "paging_token": "498216218624",
            "ledger": 116,
            "created_at": "2017-12-22T15:51:08Z",
            "changes": [
              {
                "type_i": 0,
                "type": "created",
                "created": {
                  "last_modified_ledger_seq": 116,
                  "type_i": 0,
                  "type": "account",
                  "account": {
                    "account_id": "GDS67HI27XJIJEL7IGHVJVNHPXZLMW6F3O45OXIMKAUNGIR2ROBUKTT4",
                    "account_type_i": 5,
                    "account_type": "AccountTypeNotVerified",
                    "block_reasons_i": 0,
                    "block_reasons": [],
                    "limits": null,
                    "policies": {
                      "account_policies_type_i": 0,
                      "account_policies_types": null
                    },
                    "signers": [],
                    "thresholds": {
                      "low_threshold": 0,
                      "med_threshold": 0,
                      "high_threshold": 0
                    }
                  },
                  "asset": null,
                  "balance": null
                },
                "updated": null,
                "removed": null,
                "state": null
              },
              {
                "type_i": 1,
                "type": "updated",
                "created": null,
                "updated": {
                  "last_modified_ledger_seq": 117,
                  "type_i": 4,
                  "type": "balance",
                  "account": null,
                  "asset": null,
                  "balance": {
                    "account_id": "GDJIZI4U67IZPWV26PYMPSIQTVTZAUDDID5PLA7W54ZKW6TEB664UQZT",
                    "balance_id": "BBS5KRCNZZR2MRKXMJU2SAAXYFBTSJARCKNZVKTDWSIA62SQIERJP2GX",
                    "asset": "SUN",
                    "amount": "1500.0000",
                    "locked": "0.0000"
                  }
                },
                "removed": null,
                "state": null
              }
            ]
        },
        {
          "id": "502511177728",
          "paging_token": "502511177728",
          "ledger": 117,
          "created_at": "2017-12-22T15:51:13Z",
          "changes": [
              {
                "type_i": 1,
                "type": "updated",
                "created": null,
                "updated": {
                  "last_modified_ledger_seq": 117,
                  "type_i": 6,
                  "type": "asset",
                  "account": null,
                  "asset": {
                    "code": "SUN",
                    "owner": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
                    "details": {
                      "description": "Description",
                      "external_resource_link": "https://myasset.com",
                      "logo_id": "SUN Logo",
                      "name": "SUN name"
                    },
                    "policies": {
                      "policy": 3,
                      "policies": [
                        {
                          "name": "AssetPolicyTransferable",
                          "value": 1
                        },
                        {
                          "name": "AssetPolicyBaseAsset",
                          "value": 2
                        }
                      ]
                    },
                    "preissued_asset_signer": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
                    "available_for_issueance": "98500.0000",
                    "issued": "1500.0000",
                    "max_issuance_amount": "100000000.0000"
                  },
                  "balance": null
                },
                "removed": null,
                "state": null
              }
            ]
        }
    
    ]
  }
}                                                     
```
### Send Payment

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
