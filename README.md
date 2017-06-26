# Suggester

Suggester is designed as a web service for providing input suggestions. It has 4 URLs,

* Status

`GET /:prefix/status`

* Add Index

`PUT /:prefix/:unit_id/:word/:id`

* Del Index

`DELETE /:prefix/:unit_id/:word/:id`

* Search Index

`GET /:prefix/:unit_id/:word`

Also you can use `suggester` sub package to operate it.