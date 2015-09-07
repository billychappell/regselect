# regselect
Package regselect sets Windows registry values according to a JSON config file.

For your config files to unmarshal properly, please construct your files with the same format as indicated in the following example/schema:

Example:
```
	[
		{
			"path": "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
			"scope": "CURRENT_USER",
			"properties": [
				{
					"name": "ProxyServer",
					"type": "String",
					"value": "192.168.0.227:8080"
				},
				{
					"name": "ProxyEnable",
					"type": "DWord",
					"value": 1
				},
				{
					"name": "WarnOnIntranet",
					"type": "DWord",
					"value": 0
				}
			]
		}
	]
```

Schema:
```
	{
	  "$schema": "http://json-schema.org/draft-04/schema#",
	  "id": "/",
	  "type": "array",
	  "items": {
	    "id": "2",
	    "type": "object",
	    "properties": {
	      "path": {
	        "id": "path",
	        "type": "string"
	      },
	      "scope": {
	        "id": "scope",
	        "type": "string"
	      },
	      "properties": {
	        "id": "properties",
	        "type": "array",
	        "items": {
	          "id": "0",
	          "type": "object",
	          "properties": {
	            "name": {
	              "id": "name",
	              "type": "string"
	            },
	            "type": {
	              "id": "type",
	              "type": "string"
	            },
	            "value": {
	              "id": "value",
	              "type": "integer"
	            }
	          }
	        }
	      }
	    },
	    "required": [
	      "path",
	      "scope",
	      "properties"
	    ]
	  },
	  "required": [
	    "2"
	  ]
	}
```
