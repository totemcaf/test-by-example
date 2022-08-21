# Test-by-Example

Test-By-Example is an open source Behaviour Driven Development solution to write
acceptance tests for REST API.

It can be used as a codeless solution writing the test flows in YAML (or JSON), or it can also
be used as a library to write the test flows in GO language.

A test flow is a list of steps, each step is a request to REST api with its companion expected
response.

Both the request and response can have placeholder (see #Expressions) to insert predefined values, 
and also the response can have `extractors` to extract values from the response. The extracted values can
be used in the next steps. This is useful when a response provides generated values (for example
an id) that should be used in following requests.

The tool also has support to generate values, typically used in the request (see #Generators) to
simulate real data.

# How to use

## As a command line tool

To install the tool, run:

```bash
go install github.com/test-by-example
```

To run the tool, run:

```bash
test-by-example run TEST-FILE-PATH [TEST-FILE-PATH ...]
```

For a complete list of commands and options:

```bash
test-by-example --help
```

# Writing test flows

## Test flow
A test flow is simple a list of steps:

```yaml
apiVersion: test/v1-alpha
kind: TestFlow
metadata:
  name: sample-flow
spec:
  environment:
    apiKey: API_KEY

  values:
    baseURI: "https://gitlab.example.com/api/v4"

  steps: [] # a list of steps
```

Each step is executed in the order of the list.

If any step fails, the test flow is considered failed, and it is stopped.

The TestFlow can defined environment variables, and values to include in test context.

A Test Context is used to hold values used in placeholders to produce variable requests and
to verify variable responses.

## Step

A step is a request to the REST API with its companion expected response.

```yaml
    - get: $baseURI/projects
      name: Get the list of projects
      headers:
        Api-Key: ${apiKey}

      response:
        statusCode: 200
```

This simple step just executes a GET request to the given URI and verifies the expected status code. No
check is done in the response body or headers.

A complete Step can use any of the HTTP methods (GET, PUT, POST, DELETE, PATCH), it can configure none, one,
of several headers, an optional body, and a response.

The tool will execute the request and then verify the response.

You can verify the response status code, the response body (if any), and the response headers (if any) [TODO].

The response body is checked against the sample body provided. They are compared using the JSON comparison and in
case of differences, they are reported.

The values in the Specs can have placeholders (see [Expressions](#Expressions)), and the values in the response can have extractors.

A more complex Step can be:

```yaml
    - post: $baseURI/projects
      name: Create a project 
      headers:
        Api-Key: ${apiKey}
      body:
        name: "My Project"
        description: "My Project Description"
      response:
        statusCode: 200
        body:
          id: $(id)
          name: "My Project"
          description: "My Project Description"
          enabled: true
```

This step response body contains an example of the expected response. It can contain placeholders (see [Expressions](#Expressions)), 
and extractors (see [Extractors](#Extractors)) for values produced in the backend.

# Expressions

Expressions are used to insert values into the Specs inside field of bodies in the request and/or the response.

The values are taken from the Step context. Values in the context comes from the environment variables and values
defined in the TestFlow. They can also be generated in generators used in the request body (se [Generators](#Generators)).
Another source of values in the context are the Extractors defined in the response. 
 
| Expression                   | Sample                  | Use                                                                       |
|------------------------------|-------------------------|---------------------------------------------------------------------------|
| ${varName}                   | ${apiKe}                | Use context value for `varName`                                           |
| ${varName:generator}         | ${email:random.email}   | Generates and sets value in context                                       |
| ${varName:generator:options} | ${email:random.name:20} | Generates and sets value in context using generator with provided options |

# Extractors

Extractors allow to capture values received in the response and generated in the server. Because you have no way to
predict these value, the extractors allows to match any value received in the response.

Extractors also allow to verify the value type or format. The extracted value is stored in the context for later use. 

| Expression                   | Sample                  | Use                                                                       |
|------------------------------|-------------------------|---------------------------------------------------------------------------|
| $(varName)                   | $(creditID)             | When comparing extracts value and sets value of variable in context       |

# Generators

Generators provided random values for testing. They are implemented using the great 
[Gofakeit](https://github.com/brianvoe/gofakeit) library. 
See its [documentation](https://github.com/brianvoe/gofakeit) for details on generated values.

Generate values are stored in the context for later use.

| Name                | Sample                          | Options                         | Description                                                     |
|---------------------|---------------------------------|---------------------------------|-----------------------------------------------------------------|
| Random string       | ${user:random.string}           | max length : min length         | Generates a random string name                                  |
| Random name         | ${user:random.name}             |                                 | Generates a random person name                                  |
| Random company name | ${user:random.companyName}      |                                 | Generates a random company name                                 |
| Random email        | ${email:random.email}           | domain                          | Generates a random email address                                |
| Random phone        | ${workPhone:random.phone}       |                                 | Generates a random phone number                                 |
| Random Address      | ${workAddress:random.address}   |                                 | Generates a random complete address                             |
| Random Regex        | ${id:random.regex:/[a-z0-9-]+/} | A regular expression in slashes | Generates a random string that satisfies the regular expression |


# Modularization of Steps

A Step spec can be included directly in the TestFlow, or it can be defined in a separated file and referenced in a step
of the TestFlow using the Step name.

A Step file:

```yaml
apiVersion: test/v1-alpha
kind: TestStep
metadata:
  name: partner-creates-client
spec:
  post: ${credits}/partners/${partnerID}/clients
  headers:
    Api-Key: ${partnerApiKey}
  body:
    externalId: "${externalId:random.regex:/^[a-zA-Z]{3}-[a-zA-Z0-9]{3}-[0-9]{6}$/}"
    taxId: ${taxId:random.regex:/[A-Z][A-Z][A-Z][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]/}
    email: ${email:random.email}
    address: ${address:random.address}

  response:
    statusCode: 200
    body:
      id: $(clientID)
      externalId: $externalId
      taxId: $taxId
      emailAddress: $email
      address: $address
      enabled: true
```

This step can be included in the TestFlow by referencing it by name:

```yaml
apiVersion: test/v1-alpha
kind: TestFlow
metadata:
  name: credit-request-flow
spec:
  environment:
    partnerApiKey: PARTNER_API_KEY
    partnerID: PARTNER_ID

  values:
    credits: "https://api.some-server.com/credits/v1"
  steps:
    - name: partner-creates-client
    - name: partner-starts-bnpl-flow
```

A Step inside a TestFlow to be considered a 'reference' should include none URL and none body,
nor response.

# Requirements

// Capture of values !!

* [X] Support defined values
* [X] JSON bodies and responses
* [X] Capture values from responses, simple and complex
* [ ] Check response values data types (checkers)
    * [ ] Accept datetime in ISO format $(:datetime)
    * [ ] Accept MongoDB ObjectID $(:objectID)
    * [ ] Accept UUID $(:uuid)
    * [X] Accept Integer Numbers $(:int)
    * [ ] Accept Decimal Numbers $(:decimal)
    * [ ] Accept Decimal Numbers $(:float)
    * [ ] Allow base64 encrypted data (shown decoded) $(:base64)
    * [ ] Allow base64 encrypted data (shown encoded) $(:base64encoded)
    * [X] Allow regexp $(:/regexp/)
* [ ] All HTTP methods
* [X] Allow to define headers
* [ ] Allow to check response headers
* [X] Random sample values
* [ ] Support extractors in expressions 
* [ ] Support any value in diff (something like ignore this value)
* [X] Allow Step definitions and Flows referencing defined steps so a step can be used in different flows
* [ ] Allow to define variables in steps
* [X] Allow to define variables in flows
* [X] Allow to include Steps from other files
* [ ] Allow to include Flows from other files
* [X] Allow to configure log verbosity (to debug runs)
* [X] Add a license

# See

* [Kubernetes API Conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
* [Specification by example](https://en.wikipedia.org/wiki/Specification_by_example)
