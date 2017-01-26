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

/**
 * ObjectId encoding/decoding
 *
 */
public class ObjectIdUtil {
   // Forward slash must be encoded in URLs
   private static final String FORWARD_SLASH = "/";
   // Single encoding
   private static final String FORWARD_SLASH_ENCODED1 = "%2F";
   // Double encoding
   private static final String FORWARD_SLASH_ENCODED2 = "%252F";

   /**
    * Encode the given object id.
    *
    * @param objectId The objectId to be encoded
    * @return The encoded objectId
    */
   public static String encodeObjectId(String objectId) {
      if (objectId != null) {
         return objectId.replace(FORWARD_SLASH, FORWARD_SLASH_ENCODED2);
      }
      return objectId;
   }

   /**
    * Decode the given objectId when passed as a path variable in a Spring controller
    * (Spring already performs 1 level of decoding)
    *
    * @param objectId Encoded id
    * @return The decoded object id
    */
   public static String decodePathVariable(String objectId) {
      if (objectId != null) {
         return objectId.replace(FORWARD_SLASH_ENCODED1, FORWARD_SLASH);
      }
      return objectId;
   }

   /**
    * Decode the given objectId when passed as a URL parameter, i.e. reverse the
    * double encoding done by encodeObjectId.
    *
    * @param objectId Encoded id
    * @return The decoded object id
    */
   public static String decodeParameter(String objectId) {
      if (objectId != null) {
         return objectId.replace(FORWARD_SLASH_ENCODED2, FORWARD_SLASH);
      }
      return objectId;
   }

}
