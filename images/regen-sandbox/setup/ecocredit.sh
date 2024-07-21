source $(dirname $0)/utils.sh

set -e

TX_FLAGS="--from $ADDR1 --yes --fees 5000uregen"

echo "INFO: Creating credit class..."
regen tx ecocredit create-class $ADDR1 C "Test Credit Class" --class-fee 20000000uregen $TX_FLAGS | log_response

CLASS_ID=$(regen q ecocredit classes | jq -r '.classes[-1].id')
echo "INFO:   Credit Class ID: $CLASS_ID"

echo "INFO: Creating project..."
regen tx ecocredit create-project US "Horsetail Ranch" --class $CLASS_ID $TX_FLAGS | log_response

PROJECT_ID=$(regen q ecocredit projects | jq -r '.projects[-1].id')
echo "INFO:   Project ID: $CLASS_ID"

echo "INFO: Creating credit batch..."
TEMPDIR=$(mktemp -d)
trap "rm -rf $TEMPDIR" 0 2 3 15

cat > $TEMPDIR/batch.json <<EOL
{
  "class_id": "$CLASS_ID",
  "project_id": "PROJECT_ID",
  "issuer": "$ADDR1",
  "issuance": [
    {
      "recipient": "$ADDR1",
      "tradable_amount": "1000",
      "retired_amount": "500",
      "retirement_jurisdiction": "US-WA"
    },
    {
      "recipient": "$ADDR2",
      "tradable_amount": "1000",
      "retired_amount": "500",
      "retirement_jurisdiction": "US-OR"
    }
  ],
  "metadata": "regen:13toVgf5UjYBz6J29x28pLQyjKz5FpcW3f4bT5uRKGxGREWGKjEdXYG.rdf",
  "start_date": "2020-01-01T00:00:00Z",
  "end_date": "2021-01-01T00:00:00Z",
  "open": false
}
EOL

regen tx ecocredit create-batch $TEMPDIR/batch.json $TX_FLAGS | log_response
BATCH_ID=$(regen q ecocredit batches | jq -r '.batches[-1].id')
echo "INFO:   Batch ID: $BATCh_ID"


echo "INFO: Creating NCT basket (with $CLASS_ID as allowed credit class)"
regen tx ecocredit create-basket NCT --credit-type-abbrev C --allowed-classes $CLASS_ID --basket-fee 20000000uregen --description "Testing NCT Basket" $TX_FLAGS | log_response
