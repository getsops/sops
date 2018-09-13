# language: en
@glacier @client
Feature: Amazon Glacier

  Scenario: Making a request
    When I call the "ListVaults" API
    Then the response should contain a "VaultList"

  Scenario: Handling errors
    When I attempt to call the "ListVaults" API with:
    | accountId | abcmnoxyz |
    Then the request should fail
