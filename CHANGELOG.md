# Changelog

Notable changes will be documented in this file

## Unreleased

* Unexport Serializable interface.
* Update to version 0.10.M1 compatibility
* On OfferingDescription: InputData is now Inputs, OutputData is now Outputs,
  RdfURI has been renamed to Category
* Capture and return error responses from API (remove ErrUnexpectedResponse type)
* Add integration tests to run against the live marketplace.

## v0.0.4

* Validate token method returns offering ID.

## v0.0.3

* Implement token validation to parse and validate incoming consumer request tokens
* Add api for activating offerings without sending complete description

## v0.0.2

* Add ability to delete offerings
* Tidy of code, unexport some exported properties

## v0.0.1

* First working version of the library
* Able to authenticate with marketplace
* Able to register an offering