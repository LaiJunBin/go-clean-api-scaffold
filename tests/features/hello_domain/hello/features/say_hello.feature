Feature: Say Hello API
  As an API client
  I want to ask the service for a greeting
  So that I can receive a hello message

  Scenario: Successful greeting
    When I send a "GET" request to "/api/v1/hello?name=Alice":
    Then the response code should be 200
    And the response should match json:
      """
      {
        "success": true,
        "data": {
          "message": "Hello, Alice!"
        }
      }
      """

  Scenario: Missing name
    When I send a "GET" request to "/api/v1/hello":
    Then the response code should be 400
    And the response should match json:
      """
      {
        "success": false,
        "message": "Bad Request",
        "details": [
          "name is required"
        ]
      }
      """
