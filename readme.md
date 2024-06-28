# Summary
It's a very basic implementation of how a full-text search engine works.
Ideally people would use either Elasticsearch/Postgres for implementing search, but if you have a very static and small dataset upon which you want
to enable FTS then this project can be a good starting point :)

# Dataset used
The dataset being used is MediaWiki. Can be downloaded from [here](https://dumps.wikimedia.org/enwiki/latest/)

## Next steps:
* Converting it into a server rather than CLI.
* Endpoints needed:
    * /v1/doc/insert
        *** Will allow inserting new docs
    * /v1/doc/delete
        *** This will delete the doc, soft delete that will get acknowledged eventually.
    * /v1/doc/query
        *** Runs a query against the inmemory data-structure and returns the matching docs sorted on score desc order.
* Worker that will sink in the soft delete (hard delete)
* Go routines to hasten up the init time i.e. building in-memory datastrucutres needed for querying.
