package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

// SwaggerJSON embedded version of the swagger document used at generation time
var SwaggerJSON json.RawMessage

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/vnd.goswagger.examples.task-tracker.v1+json"
  ],
  "produces": [
    "application/vnd.goswagger.examples.task-tracker.v1+json"
  ],
  "schemes": [
    "http",
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "description": "This application implements a very simple issue tracker.\nIt's implemented as an API which is described by this swagger spec document.\n\nThe go-swagger project uses this specification to test the code generation.\nThis document contains all possible values for a swagger definition.\nThis means that it exercises the framework relatively well.\n",
    "title": "Issue Tracker API",
    "termsOfService": "/termsOfService.html",
    "contact": {
      "name": "Issue Tracker API Team",
      "url": "https://task-tracker.goswagger.io",
      "email": "nobody@nowhere.com"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "1.0.0"
  },
  "host": "localhost:8322",
  "basePath": "/api",
  "paths": {
    "/tasks": {
      "get": {
        "description": "Allows for specifying a number of filter parameters to\nnarrow down the results.\nAlso allows for specifying a **sinceId** and **pageSize** parameter\nto page through large result sets.\n",
        "tags": [
          "tasks"
        ],
        "summary": "Lists the tasks",
        "operationId": "listTasks",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "description": "The last id that was seen.",
            "name": "sinceId",
            "in": "query"
          },
          {
            "uniqueItems": true,
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "the tags to filter by",
            "name": "tags",
            "in": "query"
          },
          {
            "uniqueItems": true,
            "type": "array",
            "items": {
              "enum": [
                "open",
                "closed",
                "ignored",
                "rejected"
              ],
              "type": "string"
            },
            "collectionFormat": "pipes",
            "description": "the status to filter by",
            "name": "status",
            "in": "query"
          },
          {
            "$ref": "#/parameters/pageSize"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "type": "array",
              "title": "TaskList",
              "items": {
                "$ref": "#/definitions/TaskCard"
              }
            },
            "headers": {
              "X-Last-Task-Id": {
                "type": "integer",
                "format": "int64",
                "description": "the last task id known to the application"
              }
            }
          },
          "422": {
            "description": "Validation error",
            "schema": {
              "$ref": "#/definitions/ValidationError"
            }
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      },
      "post": {
        "description": "Allows for creating a task.\nThis operation requires authentication so that we know which user\ncreated the task.\n",
        "tags": [
          "tasks"
        ],
        "summary": "Creates a 'Task' object.",
        "operationId": "createTask",
        "security": [
          {
            "api_key": []
          },
          {
            "token_header": []
          }
        ],
        "parameters": [
          {
            "description": "The task to create",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Task"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Task created"
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      }
    },
    "/tasks/{id}": {
      "get": {
        "description": "The details view has more information than the card view.\nYou can see who reported the issue and who last updated it when.\n\nThere are also comments for each issue.\n",
        "tags": [
          "tasks"
        ],
        "summary": "Gets the details for a task.",
        "operationId": "getTaskDetails",
        "responses": {
          "200": {
            "description": "Task details",
            "schema": {
              "$ref": "#/definitions/Task"
            }
          },
          "422": {
            "description": "Validation error",
            "schema": {
              "$ref": "#/definitions/ValidationError"
            }
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      },
      "put": {
        "description": "Allows for updating a task.\nThis operation requires authentication so that we know which user\nlast updated the task.\n",
        "tags": [
          "tasks"
        ],
        "summary": "Updates the details for a task.",
        "operationId": "updateTask",
        "security": [
          {
            "api_key": []
          },
          {
            "token_header": []
          }
        ],
        "parameters": [
          {
            "description": "The task to update",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Task"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Task details",
            "schema": {
              "$ref": "#/definitions/Task"
            }
          },
          "422": {
            "description": "Validation error",
            "schema": {
              "$ref": "#/definitions/ValidationError"
            }
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      },
      "delete": {
        "description": "This is a soft delete and changes the task status to ignored.\n",
        "tags": [
          "tasks"
        ],
        "summary": "Deletes a task.",
        "operationId": "deleteTask",
        "security": [
          {
            "api_key": []
          },
          {
            "token_header": []
          }
        ],
        "responses": {
          "204": {
            "description": "Task deleted"
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      },
      "parameters": [
        {
          "$ref": "#/parameters/idPathParam"
        }
      ]
    },
    "/tasks/{id}/comments": {
      "get": {
        "description": "The comments require a size parameter.\n",
        "tags": [
          "tasks"
        ],
        "summary": "Gets the comments for a task",
        "operationId": "getTaskComments",
        "parameters": [
          {
            "$ref": "#/parameters/pageSize"
          },
          {
            "type": "string",
            "format": "date-time",
            "description": "The created time of the oldest seen comment",
            "name": "since",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "The list of comments",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Comment"
              }
            }
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      },
      "post": {
        "description": "The comment can contain ___github markdown___ syntax.\nFenced codeblocks etc are supported through pygments.\n",
        "tags": [
          "tasks"
        ],
        "summary": "Adds a comment to a task",
        "operationId": "addCommentToTask",
        "security": [
          {
            "api_key": []
          },
          {
            "token_header": []
          }
        ],
        "parameters": [
          {
            "$ref": "#/parameters/idPathParam"
          },
          {
            "description": "The comment to add",
            "name": "body",
            "in": "body",
            "schema": {
              "description": "These values can have github flavored markdown.\n",
              "type": "object",
              "title": "A comment to create",
              "required": [
                "content",
                "userId"
              ],
              "properties": {
                "content": {
                  "type": "string"
                },
                "userId": {
                  "type": "integer",
                  "format": "int64"
                }
              }
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Comment added"
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      },
      "parameters": [
        {
          "$ref": "#/parameters/idPathParam"
        }
      ]
    },
    "/tasks/{id}/files": {
      "post": {
        "description": "The file can't be larger than **5MB**",
        "consumes": [
          "multipart/form-data"
        ],
        "tags": [
          "tasks"
        ],
        "summary": "Adds a file to a task.",
        "operationId": "uploadTaskFile",
        "security": [
          {
            "api_key": []
          },
          {
            "token_header": []
          }
        ],
        "parameters": [
          {
            "type": "file",
            "description": "The file to upload",
            "name": "file",
            "in": "formData"
          },
          {
            "type": "string",
            "description": "Extra information describing the file",
            "name": "description",
            "in": "formData"
          }
        ],
        "responses": {
          "201": {
            "description": "File added"
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      },
      "parameters": [
        {
          "$ref": "#/parameters/idPathParam"
        }
      ]
    }
  },
  "definitions": {
    "Comment": {
      "description": "Users can comment on issues to discuss plans for resolution etc.\n",
      "type": "object",
      "title": "A comment for an issue.",
      "required": [
        "user",
        "content"
      ],
      "properties": {
        "content": {
          "description": "This is a free text field with support for github flavored markdown.\n",
          "type": "string",
          "title": "The content of the comment."
        },
        "createdAt": {
          "description": "This field is autogenerated when the content is posted.",
          "type": "string",
          "format": "date-time",
          "title": "The time at which this comment was created.",
          "readOnly": true
        },
        "user": {
          "$ref": "#/definitions/UserCard"
        }
      }
    },
    "Error": {
      "description": "Contains all the properties any error response from the API will contain.\nSome properties are optional so might be empty most of the time\n",
      "type": "object",
      "title": "Error Structure",
      "required": [
        "code",
        "message"
      ],
      "properties": {
        "code": {
          "description": "the error code, this is not necessarily the http status code",
          "type": "integer",
          "format": "int32"
        },
        "helpUrl": {
          "description": "an optional url for getting more help about this error",
          "type": "string",
          "format": "uri"
        },
        "message": {
          "description": "a human readable version of the error",
          "type": "string"
        }
      }
    },
    "Milestone": {
      "description": "Milestones can have a escription and due date.\nThis can be useful for filters and such.\n",
      "type": "object",
      "title": "A milestone is a particular goal that is important to the project for this issue tracker.",
      "required": [
        "name"
      ],
      "properties": {
        "description": {
          "description": "A description is a free text field that allows for a more detailed explanation of what the milestone is trying to achieve.\n",
          "type": "string",
          "title": "The description of the milestone."
        },
        "dueDate": {
          "description": "This property is optional, but when present it lets people know when they can expect this milestone to be completed.\n",
          "type": "string",
          "format": "date",
          "title": "An optional due date for this milestone."
        },
        "name": {
          "description": "Each milestone should get a unique name.\n",
          "type": "string",
          "title": "The name of the milestone.",
          "maxLength": 50,
          "minLength": 3,
          "pattern": "[A-Za-z][\\w- ]+"
        },
        "stats": {
          "description": "This object contains counts for the remaining open issues and the amount of issues that have been closed.\n",
          "type": "object",
          "title": "Some counters for this milestone.",
          "properties": {
            "closed": {
              "type": "integer",
              "format": "int32",
              "title": "The closed issues."
            },
            "open": {
              "type": "integer",
              "format": "int32",
              "title": "The remaining open issues."
            },
            "total": {
              "type": "integer",
              "format": "int32",
              "title": "The total number of issues for this milestone."
            }
          }
        }
      }
    },
    "Task": {
      "description": "A Task is the main entity in this application. Everything revolves around tasks and managing them.\n",
      "type": "object",
      "title": "a structure describing a complete task.",
      "allOf": [
        {
          "$ref": "#/definitions/TaskCard"
        },
        {
          "type": "object",
          "properties": {
            "attachments": {
              "description": "An issue can have at most 20 files attached to it.\n",
              "type": "object",
              "title": "The attached files.",
              "additionalProperties": {
                "type": "object",
                "maxProperties": 20,
                "properties": {
                  "contentType": {
                    "description": "The content type of the file is inferred from the upload request.\n",
                    "type": "string",
                    "title": "The content type of the file.",
                    "readOnly": true
                  },
                  "description": {
                    "description": "This is a free form text field with support for github flavored markdown.\n",
                    "type": "string",
                    "title": "Extra information to attach to the file.",
                    "minLength": 3
                  },
                  "name": {
                    "description": "This name is inferred from the upload request.\n",
                    "type": "string",
                    "title": "The name of the file.",
                    "readOnly": true
                  },
                  "size": {
                    "description": "This property was generated during the upload request of the file.",
                    "type": "number",
                    "format": "float64",
                    "title": "The file size in bytes.",
                    "readOnly": true
                  },
                  "url": {
                    "description": "This URL is generated on the server, based on where it was able to store the file when it was uploaded.\n",
                    "type": "string",
                    "format": "uri",
                    "title": "The url to download or view the file.",
                    "readOnly": true
                  }
                }
              }
            },
            "comments": {
              "description": "The detail view of an issue includes the 5 most recent comments.\nThis field is read only, comments are added through a separate process.\n",
              "type": "array",
              "title": "The 5 most recent items for this issue.",
              "items": {
                "$ref": "#/definitions/Comment"
              },
              "readOnly": true
            },
            "lastUpdated": {
              "description": "This field is read only so it's only sent as part of the response.\n",
              "type": "string",
              "format": "date-time",
              "title": "The time at which this issue was last updated.",
              "readOnly": true
            },
            "lastUpdatedBy": {
              "$ref": "#/definitions/UserCard"
            },
            "reportedBy": {
              "$ref": "#/definitions/UserCard"
            }
          }
        }
      ]
    },
    "TaskCard": {
      "description": "A task card is a minimalistic representation of a task. Useful for display in list views, like a card list.\n",
      "type": "object",
      "title": "a card for a task",
      "required": [
        "title",
        "status"
      ],
      "properties": {
        "assignedTo": {
          "$ref": "#/definitions/UserCard"
        },
        "description": {
          "description": "The task description is a longer, more detailed description of the issue.\nPerhaps it even mentions steps to reproduce.\n",
          "type": "string",
          "title": "The description of the task."
        },
        "effort": {
          "description": "the level of effort required to get this task completed",
          "type": "integer",
          "format": "int32",
          "maximum": 27,
          "multipleOf": 3
        },
        "id": {
          "description": "A unique identifier for the task. These are created in ascending order.",
          "type": "integer",
          "format": "int64",
          "title": "The id of the task.",
          "readOnly": true
        },
        "karma": {
          "description": "Karma is a lot like voting.  Users can donate a certain amount or karma to an issue.\nThis is used to determine the weight users place on an issue. Not that +1 comments aren't great.\n",
          "type": "number",
          "format": "float32",
          "title": "the karma donated to this item.",
          "minimum": 0,
          "exclusiveMinimum": true,
          "multipleOf": 0.5
        },
        "milestone": {
          "$ref": "#/definitions/Milestone"
        },
        "reportedAt": {
          "description": "This field is read-only, so it's only sent as part of the response.\n",
          "type": "string",
          "format": "date-time",
          "title": "The time at which this issue was reported.",
          "readOnly": true
        },
        "severity": {
          "type": "integer",
          "format": "int32",
          "maximum": 5,
          "minimum": 1
        },
        "status": {
          "description": "There are 4 possible values for a status.\nIgnored means as much as accepted but not now, perhaps later.\n",
          "type": "string",
          "title": "the status of the issue",
          "enum": [
            "open",
            "closed",
            "ignored",
            "rejected"
          ]
        },
        "tags": {
          "description": "a task can be tagged with text blurbs.",
          "type": "array",
          "title": "task tags.",
          "maxItems": 5,
          "uniqueItems": true,
          "items": {
            "type": "string",
            "minLength": 3,
            "pattern": "\\w[\\w- ]+"
          }
        },
        "title": {
          "description": "The title for a task, this needs to be at least 5 chars long.\nTitles don't allow any formatting, besides emoji.\n",
          "type": "string",
          "title": "The title of the task.",
          "maxLength": 150,
          "minLength": 5
        }
      }
    },
    "UserCard": {
      "description": "This representation of a user is mainly meant for inclusion in other models, or for list views.\n",
      "type": "object",
      "title": "A minimal representation of a user.",
      "required": [
        "id",
        "screenName"
      ],
      "properties": {
        "admin": {
          "description": "Only employees of the owning company can be admins.\nAdmins are like project owners but have access to all the projects in the application.\nThere aren't many admins, and it's only used for extremly critical issues with the application.\n",
          "type": "boolean",
          "title": "When true this user is an admin.",
          "readOnly": true
        },
        "availableKarma": {
          "description": "In this application users get a cerain amount of karma alotted.\nThis karma can be donated to other users to show appreciation, or it can be used\nby a user to vote on issues.\nOnce an issue is closed or rejected, the user gets his karma back.\n",
          "type": "number",
          "format": "float32",
          "title": "The amount of karma this user has available.",
          "maximum": 1000,
          "exclusiveMaximum": true,
          "readOnly": true
        },
        "id": {
          "description": "This id is automatically generated on the server when a user is created.\n",
          "type": "integer",
          "format": "int64",
          "title": "A unique identifier for a user.",
          "readOnly": true
        },
        "screenName": {
          "description": "This is used for vanity type urls as well as login credentials.\n",
          "type": "string",
          "title": "The screen name for the user.",
          "maxLength": 255,
          "minLength": 3,
          "pattern": "\\w[\\w_-]+"
        }
      }
    },
    "ValidationError": {
      "allOf": [
        {
          "$ref": "#/definitions/Error"
        },
        {
          "type": "object",
          "properties": {
            "field": {
              "description": "an optional field name to which this validation error applies",
              "type": "string"
            }
          }
        }
      ]
    }
  },
  "parameters": {
    "idPathParam": {
      "type": "integer",
      "format": "int64",
      "description": "The id of the item",
      "name": "id",
      "in": "path",
      "required": true
    },
    "pageSize": {
      "type": "integer",
      "format": "int32",
      "default": 20,
      "description": "Amount of items to return in a single page",
      "name": "pageSize",
      "in": "query"
    }
  },
  "responses": {
    "ErrorResponse": {
      "description": "Error response",
      "schema": {
        "$ref": "#/definitions/Error"
      },
      "headers": {
        "X-Error-Code": {
          "type": "string"
        }
      }
    }
  },
  "securityDefinitions": {
    "api_key": {
      "type": "apiKey",
      "name": "token",
      "in": "query"
    },
    "token_header": {
      "type": "apiKey",
      "name": "X-Token",
      "in": "header"
    }
  },
  "tags": [
    {
      "description": "manages tasks",
      "name": "tasks",
      "externalDocs": {
        "description": "An extensive explanation on what is possible can be found in the\nsupport site for this application.\n",
        "url": "https://go-swagger.github.io/examples/tasklist/help/tasks.html"
      }
    },
    {
      "description": "manages milestones",
      "name": "milestones",
      "externalDocs": {
        "description": "An extensive explanation on what is possible can be found in the\nsupport site for this application.\n",
        "url": "https://go-swagger.github.io/examples/tasklist/help/milestones.html"
      }
    }
  ],
  "externalDocs": {
    "description": "A much more elaborate guide to this application is available at the support\nsite.\n",
    "url": "https://go-swagger.github.io/examples/tasklist/help/tasks.html"
  }
}`))
}
