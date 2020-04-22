# qorexGoHelpers Package
qorexGoHelpers provides a series of library functions and types to be used across all entities.

### Response
Gives a standard API response using http status codes and messages.  It can also contain an error on a failed api interaction or a data array.

##### SetStatus(statusCode int, statusText string, err error)
Used to set the returned response or any error

##### AppendData(entityList []interface{})
Used to append data to the internal data slice.  For efficiency data shouldn't be appended as individual items.

##### GetJson()
Returns a []byte for use with the JSON.Unmarshal functionality

### RowDecorator
Using the Wrap() method a collection of Sql/Row can be iterated over and have a callback applied to each individual row

### SecretManager
SecretManager allows access to AWS secretManager service by passing in config of secret and region to return a DBConnection struct

