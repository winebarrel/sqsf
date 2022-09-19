# sqsf

sqsf is a tool to receive AWS SQL messages like `tail -f`.

## Installation

```
brew install winebarrel/sqsf/sqsf
```

## Usage

```
Usage: sqsf [OPTION] QUEUE
  -decode
    	print decoded message body
  -delete
    	delete received message (default true)
  -version
    	print version and exit
  -vis-timeout int
    	visibility timeout (default 600)
```

### Example

```
$ sqsf my-queue-name
{
    "Attributes": null,
    "Body": "{\"version\":\"1.0\",\"timestamp\":\"2022-09-19T09:01:29.773Z\",\"requestContext\":{\"requestId\":\"7e658e64-4e9f-499f-a949-fad9eb41fff0\",\"functionArn\":\"arn:aws:lambda:ap-northeast-1:123456789012:function:hello:$LATEST\",\"condition\":\"Success\",\"approximateInvokeCount\":1},\"requestPayload\":{\"key1\":100,\"key2\":200,\"key3\":300},\"responseContext\":{\"statusCode\":200,\"executedVersion\":\"$LATEST\"},\"responsePayload\":100}",
    "MD5OfBody": "e3216d7baf92ab8d3842b2c5f742cbc5",
    "MD5OfMessageAttributes": null,
    "MessageAttributes": null,
    "MessageId": "3fdc12d6-3cb8-4c0d-aaa5-b6a6d40a0d54",
    "ReceiptHandle": "..."
}
^C # Running until CTRL-C is pressed

$ sqsf -decode my-queue-name
{
    "requestContext": {
        "approximateInvokeCount": 1,
        "condition": "Success",
        "functionArn": "arn:aws:lambda:ap-northeast-1:123456789012:function:hello:$LATEST",
        "requestId": "894310eb-fc64-4f12-aa2d-9ad6a4d2c8ae"
    },
    "requestPayload": {
        "key1": 100,
        "key2": 200,
        "key3": 300
    },
    "responseContext": {
        "executedVersion": "$LATEST",
        "statusCode": 200
    },
    "responsePayload": 100,
    "timestamp": "2022-09-19T09:01:55.043Z",
    "version": "1.0"
}
```
