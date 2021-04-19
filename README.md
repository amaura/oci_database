# oci_database

```
./go_build_oci_database --help
Usage:
  go_build_oci_database [OPTIONS]

Application Options:
  -d, --db-name=                Database name
  -l, --db-workload=            Database workload : OLTP or DSS (default: OLTP)
  -o, --db-home-id=             Database home OCID
  -w, --wait-for-state=         Wait for state : AVAILABLE, TERMINATED, etc.
  -p, --admin-password=         Database password
  -u, --db-unique-name=         Database Unique Name
  -s, --character-set=          Character Set (default: AL32UTF8)
  -n, --national-character-set= National Character Set (default: AL16UTF16)
  -b, --pdb-name=               PDB name
  -x, --tde-wallet-password=    TDE Wallet Password
  -i, --wait-interval-seconds=  Wait Interval Seconds (default: 30)
  -m, --max-wait-seconds=       Max Wait Seconds (default: 3600)
  -t, --dry-run                 Display request only

Help Options:
  -h, --help                    Show this help message
```
