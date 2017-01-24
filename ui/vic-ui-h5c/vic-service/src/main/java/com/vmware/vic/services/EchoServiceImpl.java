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

package com.vmware.vic.services;

/**
 * Implementation of the EchoService interface
 */
public class EchoServiceImpl implements EchoService {

   /* (non-Javadoc)
    * @see com.vmware.vic.EchoService#echo(java.lang.String)
    */
   public String echo(String message) {
      return message;
   }
}
