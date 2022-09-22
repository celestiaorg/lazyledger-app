package coretestutil

var genesis = `
{
    "genesis_time": "2022-08-06T04:00:36.009315437Z",
    "chain_id": "qgb-e2e",
    "initial_height": "1",
    "consensus_params": {
      "block": {
        "max_bytes": "22020096",
        "max_gas": "-1",
        "time_iota_ms": "1000"
      },
      "evidence": {
        "max_age_num_blocks": "100000",
        "max_age_duration": "172800000000000",
        "max_bytes": "1048576"
      },
      "validator": {
        "pub_key_types": [
          "ed25519"
        ]
      },
      "version": {}
    },
    "app_hash": "",
    "app_state": {
      "auth": {
        "params": {
          "max_memo_characters": "256",
          "tx_sig_limit": "7",
          "tx_size_cost_per_byte": "10",
          "sig_verify_cost_ed25519": "590",
          "sig_verify_cost_secp256k1": "1000"
        },
        "accounts": [
          {
            "@type": "/cosmos.auth.v1beta1.BaseAccount",
            "address": "celestia1ly9zghpffkw7gca42qkqm3awmw2zatc62yey3m",
            "pub_key": null,
            "account_number": "0",
            "sequence": "0"
          },
          {
            "@type": "/cosmos.auth.v1beta1.BaseAccount",
            "address": "celestia1s67n6xvmslyyjj0ea930vrzhtcsty8aus3rpfq",
            "pub_key": null,
            "account_number": "0",
            "sequence": "0"
          },
          {
            "@type": "/cosmos.auth.v1beta1.BaseAccount",
            "address": "celestia1l5y4swfghcfz8enf2l2c7fe3r6lna4rad4n862",
            "pub_key": null,
            "account_number": "0",
            "sequence": "0"
          },
          {
            "@type": "/cosmos.auth.v1beta1.BaseAccount",
            "address": "celestia1qde6jpujzx2nptvdjj5zzarv6wg0tklp4ymwv8",
            "pub_key": null,
            "account_number": "0",
            "sequence": "0"
          }
        ]
      },
      "bank": {
        "params": {
          "send_enabled": [],
          "default_send_enabled": true
        },
        "balances": [
          {
            "address": "celestia1qde6jpujzx2nptvdjj5zzarv6wg0tklp4ymwv8",
            "coins": [
              {
                "denom": "utia",
                "amount": "2000000000000000"
              }
            ]
          },
          {
            "address": "celestia1s67n6xvmslyyjj0ea930vrzhtcsty8aus3rpfq",
            "coins": [
              {
                "denom": "utia",
                "amount": "2000000000000000"
              }
            ]
          },
          {
            "address": "celestia1ly9zghpffkw7gca42qkqm3awmw2zatc62yey3m",
            "coins": [
              {
                "denom": "utia",
                "amount": "2000000000000000"
              }
            ]
          },
          {
            "address": "celestia1l5y4swfghcfz8enf2l2c7fe3r6lna4rad4n862",
            "coins": [
              {
                "denom": "utia",
                "amount": "2000000000000000"
              }
            ]
          }
        ],
        "supply": [
          {
            "denom": "utia",
            "amount": "8000000000000000"
          }
        ],
        "denom_metadata": [
          {
            "description": "The native staking token of the Celestia network.",
            "denom_units": [
              {
                "denom": "utia",
                "exponent": 0,
                "aliases": [
                  "microtia"
                ]
              },
              {
                "denom": "TIA",
                "exponent": 6,
                "aliases": []
              }
            ],
            "base": "utia",
            "display": "TIA",
            "name": "TIA",
            "symbol": "TIA",
            "uri": "",
            "uri_hash": ""
          }
        ]
      },
      "capability": {
        "index": "1",
        "owners": []
      },
      "crisis": {
        "constant_fee": {
          "denom": "utia",
          "amount": "1000"
        }
      },
      "distribution": {
        "params": {
          "community_tax": "0.020000000000000000",
          "base_proposer_reward": "0.010000000000000000",
          "bonus_proposer_reward": "0.040000000000000000",
          "withdraw_addr_enabled": true
        },
        "fee_pool": {
          "community_pool": []
        },
        "delegator_withdraw_infos": [],
        "previous_proposer": "",
        "outstanding_rewards": [],
        "validator_accumulated_commissions": [],
        "validator_historical_rewards": [],
        "validator_current_rewards": [],
        "delegator_starting_infos": [],
        "validator_slash_events": []
      },
      "evidence": {
        "evidence": []
      },
      "feegrant": {
        "allowances": []
      },
      "genutil": {
        "gen_txs": [
          {
            "body": {
              "messages": [
                {
                  "@type": "/cosmos.staking.v1beta1.MsgCreateValidator",
                  "description": {
                    "moniker": "qgb-e2e",
                    "identity": "",
                    "website": "",
                    "security_contact": "",
                    "details": ""
                  },
                  "commission": {
                    "rate": "0.100000000000000000",
                    "max_rate": "0.200000000000000000",
                    "max_change_rate": "0.010000000000000000"
                  },
                  "min_self_delegation": "1",
                  "delegator_address": "celestia1ly9zghpffkw7gca42qkqm3awmw2zatc62yey3m",
                  "validator_address": "celestiavaloper1ly9zghpffkw7gca42qkqm3awmw2zatc60mma8a",
                  "pubkey": {
                    "@type": "/cosmos.crypto.ed25519.PubKey",
                    "key": "X9yawE0ZCabQ9+VFql8WknVUC01pm5MmxLEMGqvCuhA="
                  },
                  "value": {
                    "denom": "utia",
                    "amount": "1000000000000000"
                  },
                  "orchestrator": "celestia1ly9zghpffkw7gca42qkqm3awmw2zatc62yey3m",
                  "eth_address": "0x966e6f22781EF6a6A82BBB4DB3df8E225DfD9488"
                }
              ],
              "memo": "a6b8e77b70336cb015723d2971b9d222c5ad1bd2@10.68.33.210:26656",
              "timeout_height": "0",
              "extension_options": [],
              "non_critical_extension_options": []
            },
            "auth_info": {
              "signer_infos": [
                {
                  "public_key": {
                    "@type": "/cosmos.crypto.secp256k1.PubKey",
                    "key": "AnDneC7RH8OEACcnQPSDi126LDt7hWXAUUtCZZIlbV3N"
                  },
                  "mode_info": {
                    "single": {
                      "mode": "SIGN_MODE_DIRECT"
                    }
                  },
                  "sequence": "0"
                }
              ],
              "fee": {
                "amount": [],
                "gas_limit": "210000",
                "payer": "",
                "granter": ""
              },
              "tip": null
            },
            "signatures": [
              "cLhCUPeCSlkZ0n2W/QpK10+F/NKGoCpVYMl5KKoeMCkrsYzyTPh1gT8/GioWYNVSqEpPhVucQT9+YwCEQfjhzQ=="
            ]
          }
        ]
      },
      "gov": {
        "starting_proposal_id": "1",
        "deposits": [],
        "votes": [],
        "proposals": [],
        "deposit_params": {
          "min_deposit": [
            {
              "denom": "stake",
              "amount": "10000000"
            }
          ],
          "max_deposit_period": "172800s"
        },
        "voting_params": {
          "voting_period": "172800s"
        },
        "tally_params": {
          "quorum": "0.334000000000000000",
          "threshold": "0.500000000000000000",
          "veto_threshold": "0.334000000000000000"
        }
      },
      "mint": {
        "minter": {
          "inflation": "0.130000000000000000",
          "annual_provisions": "0.000000000000000000"
        },
        "params": {
          "mint_denom": "utia",
          "inflation_rate_change": "0.130000000000000000",
          "inflation_max": "0.200000000000000000",
          "inflation_min": "0.070000000000000000",
          "goal_bonded": "0.670000000000000000",
          "blocks_per_year": "6311520"
        }
      },
      "params": null,
      "payment": {},
      "qgb": {
        "params": {
          "data_commitment_window": "101"
        }
      },
      "slashing": {
        "params": {
          "signed_blocks_window": "100",
          "min_signed_per_window": "0.500000000000000000",
          "downtime_jail_duration": "600s",
          "slash_fraction_double_sign": "0.050000000000000000",
          "slash_fraction_downtime": "0.010000000000000000"
        },
        "signing_infos": [],
        "missed_blocks": []
      },
      "staking": {
        "params": {
          "unbonding_time": "1814400s",
          "max_validators": 100,
          "max_entries": 7,
          "historical_entries": 10000,
          "bond_denom": "utia",
          "min_commission_rate": "0.000000000000000000"
        },
        "last_total_power": "0",
        "last_validator_powers": [],
        "validators": [],
        "delegations": [],
        "unbonding_delegations": [],
        "redelegations": [],
        "exported": false
      },
      "upgrade": {},
      "vesting": {}
    }
  }
`

var privVal = `
{
	"address": "048833F6126E0D51A54E8B7351916B548D232F99",
	"pub_key": {
	  "type": "tendermint/PubKeyEd25519",
	  "value": "X9yawE0ZCabQ9+VFql8WknVUC01pm5MmxLEMGqvCuhA="
	},
	"priv_key": {
	  "type": "tendermint/PrivKeyEd25519",
	  "value": "LoptEkS2g6SvHKiLkSmBexy7xbsNJ436z7WIHeISghNf3JrATRkJptD35UWqXxaSdVQLTWmbkybEsQwaq8K6EA=="
	}
  }
`
