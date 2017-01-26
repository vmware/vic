/*

Copyright 2017 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Mac OS script starting an Ant build of the current flex project
Note: if Ant runs out of memory try defining ANT_OPTS=-Xmx512M

*/

package com.vmware.vic.mvc;

import java.net.URI;
import java.util.HashMap;
import java.util.Map;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

/**
 * Utility defining the result of an action called in an ActionsController.
 * This generic class can be re-used in any Java plugin.
 *
 * @see getJsonMap() to get the actual value to be returned in the invoke method.
 */
public class ActionResult {
   final static String ACTIONUID = "actionUid";
   final static String ERROR_MESSAGE = "errorMessage";
   final static String RESULT = "result";
   final static String OPERATION_TYPE = "operationType";
   final static String OP_ADD = "add";
   final static String OP_CHANGE = "change";
   final static String OP_DELETE = "delete";
   final static String OP_RELATIONSHIP_CHANGE = "relationship_change";

   private final static Log _logger = LogFactory.getLog(ActionResult.class);
   private Map<String, Object> _resultMap = new HashMap<String, Object>();
   private String _resourceBundle;

   /**
    * Constructor.
    * You must use one of the set___Result() methods later or setErrorMessage
    * before using this ActionResult.
    *
    * @param actionUid  The action id.
    * @param resourceBundle Name of the resource bundle to use for error messages.
    */
   public ActionResult(String actionUid, String resourceBundle) {
      _resultMap.put(ACTIONUID, actionUid);
      _resourceBundle = resourceBundle;
   }

   /**
    * @return the map to be used as the return value of the action's invoke method.
    */
   public Map<String, Object> getJsonMap() {
      if ((_resultMap.get(RESULT) == null) && (_resultMap.get(ERROR_MESSAGE) == null)) {
         _logger.error("Missing result or error message in ActionResult for " +
               _resultMap.get(ACTIONUID));
      }
      return _resultMap;
   }

   /**
    * Assign a result value to this ActionResult.
    *
    * @param result
    * @param errMsgKey  Message key of the error message to display when result is null
    *    or false, or leave null when no message is needed.
    */
   public void setResult(Object result, String errMsgKey) {
      _resultMap.put(RESULT, result);
      if ((errMsgKey != null) && ((result == null) ||
            (result instanceof Boolean && (Boolean)result == false))) {
         setErrorMessage(errMsgKey);
      }
   }

   /**
    * Set the result of an action creating a new object. This will update the UI model
    * to display the new object if the action was successful.
    *
    * @param result  the URI representing the object, or null if no object was created.
    * @param uriType the object type
    * @param errMsgKey  Message key of the error message to display when result is null,
    *    or leave null when no message is needed.
    */
   public void setObjectAddedResult(URI result, String uriType, String errMsgKey) {
      _resultMap.put(RESULT, result);
      _resultMap.put(OPERATION_TYPE, OP_ADD);
      _resultMap.put("uriType", uriType);
      if ((result == null) && (errMsgKey != null)) {
         setErrorMessage(errMsgKey);
      }
   }

   /**
    * Set the result of an action deleting an object. This will update the UI model to
    * remove the object if the action was successful.
    *
    * @param result true if the action was successful, false otherwise.
    * @param errMsgKey  Message key of the error message to display when result is false,
    *    or leave null when no message is needed.
    */
   public void setObjectDeletedResult(boolean result, String errMsgKey) {
      _resultMap.put(RESULT, result);
      _resultMap.put(OPERATION_TYPE, OP_DELETE);
      if (!result && (errMsgKey != null)) {
         setErrorMessage(errMsgKey);
      }
   }

   /**
    * Set the result of an action modifying an object. This will update the UI model to
    * display the object's changes if the action was successful.
    *
    * @param result
    * @param errMsgKey  Message key of the error message to display when result is false,
    *    or null to avoid displaying any error message.
    */
   public void setObjectChangedResult(boolean result, String errMsgKey) {
      _resultMap.put(RESULT, result);
      _resultMap.put(OPERATION_TYPE, OP_CHANGE);
      if (!result && (errMsgKey != null)) {
         setErrorMessage(errMsgKey);
      }
   }

   /**
    * Set an error message to be displayed in the UI when the action returns.
    *
    * @param msg  The message to display, already localized on the server side (or else
    *    you can use setErrorMessage())
    */
   public void setErrorLocalizedMessage(String msg) {
      _resultMap.put(ERROR_MESSAGE, msg);
   }

   /**
    * Set an error message to be displayed in the UI when the action returns. The message will be
    * localized on the client side for convenience.
    *
    * @param key  The message key.
    */
   public void setErrorMessage(String key) {
      setErrorMessage(key, null);
   }

   /**
    * Set an error message to be displayed in the UI when the action returns. The message will be
    * localized on the client side for convenience.
    * @param key  The message key.
    * @param params  An optional array of parameters when the string resources contains
    *    place holders {0}, {1} etc.  Use null otherwise.
    */
   public void setErrorMessage(String key, String[] params) {
      _resultMap.put(ERROR_MESSAGE, new ActionMessage(_resourceBundle, key, params));
   }

   class ActionMessage {
      public String bundleName;
      public String key;
      public String[] params;

      ActionMessage() {
         // empty constructor needed for Json serialization
      }

      ActionMessage(String bundleName, String key, String[] params) {
         this.bundleName = bundleName;
         this.key = key;
         this.params = params;
      }
   }
}

