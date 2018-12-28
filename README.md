# Syracuse

This microservice store user information into postgres.

### Table schema

```
   Column   |           Type           | Nullable |      Default
------------+--------------------------+-----------+-------------------
 id         | uuid                     | not null | gen_random_uuid()
 email      | text                     | not null |
 full_name  | text                     | not null |
 created_at | timestamp with time zone |          | now()
 updated_at | timestamp with time zone |          | now()
 deleted_at | timestamp with time zone |          |
```

### GRPC service

```
syntax = "proto3";

package citizens;

service Citizenship {
  rpc Get(GetRequest) returns (GetResponse) {}
  rpc Select(SelectRequest) returns (SelectResponse) {}
  rpc Create(CreateRequest) returns (CreateResponse) {}
  rpc Update(UpdateRequest) returns (UpdateResponse) {}
  rpc Delete(DeleteRequest) returns (DeleteResponse) {}
}

message Citizen {
  string id = 1;
  string full_name = 2;
  string email = 3;

  int64 created_at = 4;
  int64 updated_at = 5;
}

message GetRequest {
  string user_id = 1;
}

message GetResponse {
  Citizen data = 1;
}

message SelectRequest {}

message SelectResponse {
  repeated Citizen data = 1;
}

message CreateRequest {
  Citizen data = 1;
}

message CreateResponse {
  Citizen data = 1;
}

message UpdateRequest {
  string user_id = 1;
  Citizen data = 2;
}

message UpdateResponse {
  Citizen data = 1;
}

message DeleteRequest {
  string user_id = 1;
}

message DeleteResponse {
  Citizen data = 1;
}
```