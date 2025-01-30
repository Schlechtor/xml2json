# xml2json

This is a Go program that reads data from an XML file, applies transformations specified in a `config.yaml` file, and outputs the transformed data in JSON format. The transformations can include simple field mappings, combinations of fields, and calculations like age from a date of birth.

## Features
- Reads configuration from a `config.yaml` file.
- Processes data in an XML format.
- Supports transformations like combining fields and calculating age.
- Outputs the transformed data in a JSON format.

## Prerequisites

If you have nix you can just run the flake.
```bash
nix develop
```

## Running
Running the Program: Command-Line Usage

You can run the program from the command line by specifying the paths to the config.yaml and XML input file as arguments:

```bash
go run main.go <config.yaml> <input.xml>
```

Expected Output

The output will be a JSON object that contains the transformed data. For example:

```json
{
  "Patient": [
    {
      "age": 39,
      "id": "12345",
      "name": "John Doe",
      "operations": {
        "NewStuff": {
          "More": "Testing",
          "NewStuff2": {
            "More2": "Testing",
            "Next2": "other"
          },
          "Next": "other"
        },
        "Other": "other",
        "Test": "Testing"
      },
      "phone": "333-456-7890"
    },
    {
      "age": 32,
      "id": "67890",
      "name": "Jane Smith",
      "phone": "123-111-7890"
    },
    {
      "age": 32,
      "id": "11111",
      "name": "John Doe",
      "phone": "123-111-7890"
    }
  ]
}
```

## Testing

```bash
go test ./... -v
```

## Configuration

The transformation logic is defined in the config.yaml file. Here's an example of how to structure the config.yaml file: fieldName: "XMLElementName"

```yaml
topLevelField: "Patient"
patients:
  fields:
    id: "ID"
    name: "FirstName LastName"
    phone: "PhoneNumber"
    age:
      calculate: "DateOfBirth"
```

In this example:

    The id field maps to the ID attribute in the XML.
    The name field combines the FirstName and LastName fields from the XML.
    The phone field maps directly to the PhoneNumber field in the XML.
    The dob field is calculated based on the DateOfBirth field in the XML.

Fields:

    String mapping: A field in the configuration file can map directly to an XML element. For example:

Complex transformations: For calculations or combinations, you can specify an operation, like calculating age based on the date of birth:

    fieldName:
      calculate: "DateOfBirth"

## Adding Transformations

Start by creating a new file under the transformations folder, then add your logic into a function. Just make sure they implement the init function.

```go
func init() {
	Register("combine", combineFields)
}
```

## Other Considerations

Large datasets will take some time to process. Could add concurrency to handle those larger datasets.
XML property attributes.