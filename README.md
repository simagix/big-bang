# big-bang
Generate random data with relationships among collections.

## Build
```
dep ensure -update
go build -o bigbang
```

## Example
```
./bigbang --config examples/bigbang-config.json
```

## Validation
```
db.agents.aggregate([{$lookup: {
  from: 'agencies',
  localField: 'agencyID',
  foreignField: '_id',
  as: 'agency'
}}, {$lookup: {
  from: 'expertises',
  localField: 'skills',
  foreignField: '_id',
  as: 'skills'
}}])
```
